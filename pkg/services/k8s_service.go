package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/saas/hostgolang/pkg/types"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/util/retry"

	"github.com/saas/hostgolang/pkg/config"
	"io/ioutil"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

const stormNs = "namespace-storm"
const registrySecretName = "storm-secret"

//go:generate mockgen -destination=mocks/k8s_service_mock.go -package=mocks github.com/saas/hostgolang/pkg/services K8sService
type K8sService interface {
	DeployService(opt *types.CreateDeploymentOpts) (*types.DeploymentResult, error)
	GetLogs(appName string) (string, error)
	UpdateEnvs(appName string, envs []types.Environment) error
	ScaleApp(appName string, replicas int) error
	ListRunningPods(appName string) ([]types.Instance, error)
}

type defaultK8sService struct {
	client *kubernetes.Clientset
	cfg    *config.Config
}

func NewK8sService(client *kubernetes.Clientset, config *config.Config) K8sService {
	d := &defaultK8sService{client: client, cfg: config}
	if err := d.createNameSpace(); err != nil {
		log.Println("Error occurred while creating namespace: ", err)
	}
	return d
}

func (d *defaultK8sService) createNameSpace() error {
	c := d.client.CoreV1().Namespaces()
	if _, err := c.Get(stormNs, metav1.GetOptions{}); err != nil {
		log.Println("Error: ", err, ": Creating namespace...")
		ns := &v1.Namespace{}
		ns.Name = stormNs
		if _, err := c.Create(ns); err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultK8sService) DeployService(opt *types.CreateDeploymentOpts) (*types.DeploymentResult, error) {
	var serviceType = v1.ServiceTypeNodePort
	if !opt.IsLocal {
		serviceType = v1.ServiceTypeLoadBalancer
	}
	if err := d.createRegistrySecret(); err != nil {
		return nil, err
	}
	name := opt.Name
	svc, err := d.createService(name, serviceType)
	if err != nil {
		return nil, err
	}
	ports := svc.Spec.Ports
	deployment, err := d.client.AppsV1().Deployments(stormNs).Get(name, metav1.GetOptions{})
	if err == nil && deployment.Name != "" {
		if err := d.client.AppsV1().Deployments(stormNs).Delete(name, &metav1.DeleteOptions{}); err != nil {
			return nil, err
		}
	}
	if err := d.createDeployment(opt.Tag, name, opt.Envs, svc.Labels, ports, opt.Replicas, opt.Memory, opt.Cpu); err != nil {
		return nil, err
	}

	newDeployment, err := d.client.AppsV1().Deployments(stormNs).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	nodeName := newDeployment.Spec.Template.Spec.NodeName
	go func() {
		time.Sleep(3 * time.Minute)
		log.Println("Running GR...")
		nd, err := d.client.AppsV1().Deployments(stormNs).Get(name, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
			return
		}
		_, err = d.getNodeIp(nd.Spec.Template.Spec.NodeName)
		log.Println("GC error: ", err)
	}()

	addr, err := d.getNodeIp(nodeName)
	log.Println("GetNode error: ", err)
	if err != nil || addr == "" {
		addr = "http://localhost"
	}
	log.Println("NodeIp: ", addr)
	accessAddr := ""
	for _, p := range ports {
		if opt.IsLocal {
			if nodePort := p.NodePort; nodePort != 0 {
				accessAddr = fmt.Sprintf("%s:%d", addr, nodePort)
			}
		}
	}
	return &types.DeploymentResult{Address: accessAddr}, nil
}

func (d *defaultK8sService) UpdateEnvs(appName string, envs []types.Environment) error {
	if envs == nil {
		envs = make([]types.Environment, 0)
	}
	c := d.client.AppsV1().Deployments(stormNs)
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		deployment, err := c.Get(appName, metav1.GetOptions{})
		if err != nil {
			return errors.New("deployment not found")
		}
		portValue := ""
		podEnvs := deployment.Spec.Template.Spec.Containers[0].Env
		for _, pe := range podEnvs {
			if strings.ToLower(pe.Name) == "port" {
				portValue = pe.Value
				break
			}
		}
		envs = append(envs, types.Environment{EnvKey: "PORT", EnvValue: portValue})
		newEnvs := make([]v1.EnvVar, 0)
		for _, newEnv := range envs {
			newEnvs = append(newEnvs, v1.EnvVar{Name: newEnv.EnvKey, Value: newEnv.EnvValue})
		}
		deployment.Spec.Template.Spec.Containers[0].Env = newEnvs
		if _, err := c.Update(deployment); err != nil {
			return err
		}
		return nil
	})
}

func (d *defaultK8sService) ScaleApp(appName string, replicas int) error {
	c := d.client.AppsV1().Deployments(stormNs)
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		deployment, err := c.Get(appName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		deployment.Spec.Replicas = Int32(int32(replicas))
		if _, err := c.Update(deployment); err != nil {
			return err
		}
		return nil
	})
}

func (d *defaultK8sService) ListRunningPods(appName string) ([]types.Instance, error) {
	pods, err := d.getPodsBySelector(fmt.Sprintf("web=%s-service", strings.ToLower(appName)))
	if err != nil {
		return nil, err
	}
	data := make([]types.Instance, 0)
	for _, p := range pods {
		started := ""
		if p.Status.StartTime != nil {
			started = p.Status.StartTime.String()
		}
		in := types.Instance{
			Id:      string(p.UID),
			Status:  string(p.Status.Phase),
			Name:    p.Name,
			Started: started,
		}
		data = append(data, in)
	}
	return data, nil
}

func (d *defaultK8sService) createNodePortService(name string) (*v1.Service, error) {
	return d.createService(name, v1.ServiceTypeNodePort)
}

func (d *defaultK8sService) createLoadBalancerService(name string) (*v1.Service, error) {
	return d.createService(name, v1.ServiceTypeLoadBalancer)
}

func (d *defaultK8sService) createService(serviceName string, serviceType v1.ServiceType) (*v1.Service, error) {
	name := serviceName
	client := d.client.CoreV1()
	svc, err := client.Services(stormNs).Get(name, metav1.GetOptions{})
	if err == nil && svc.Name != "" {
		// service already exists, delete it because we'll need to recreate it
		if err := client.Services(stormNs).Delete(name, &metav1.DeleteOptions{}); err != nil {
			return nil, err
		}
	}

	labels := d.createLabel(serviceName)
	svc = &v1.Service{}
	svc.Name = name
	svc.Labels = labels
	svc.Namespace = stormNs
	servicePort := findAvailablePort()
	port := v1.ServicePort{
		Name:       truncString(fmt.Sprintf("%15s", name+"-serviceport")),
		Protocol:   "TCP",
		Port:       int32(servicePort),
		TargetPort: intstr.FromInt(servicePort),
	}
	svc.Spec = v1.ServiceSpec{
		Ports:    []v1.ServicePort{port},
		Selector: labels,
		Type:     serviceType,
	}
	return client.Services(stormNs).Create(svc)
}

func (d *defaultK8sService) createDeployment(tag, name string, envs, labels map[string]string, ports []v1.ServicePort, replicas int32, mem, vCpu float64) error {
	cpu, err := resource.ParseQuantity(fmt.Sprintf("%fm", vCpu*1000))
	if err != nil {
		return err
	}
	memory, err := resource.ParseQuantity(fmt.Sprintf("%fMi", mem*1000))
	if err != nil {
		return err
	}

	deployment := &appsV1.Deployment{}
	deployment.Name = name
	deployment.Labels = labels

	container := v1.Container{}
	envVars := make([]v1.EnvVar, 0, len(envs))

	for k, v := range envs {
		envVars = append(envVars, v1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	var port int32 = 0
	for _, p := range ports {
		if targetPort := p.TargetPort.IntVal; targetPort != 0 {
			port = targetPort
		}
	}
	envVars = append(envVars, v1.EnvVar{Name: "PORT", Value: fmt.Sprintf("%d", port)})
	container.Name = name
	container.Env = envVars
	container.Image = tag
	container.Ports = []v1.ContainerPort{{
		Name:          truncString(fmt.Sprintf("%15s", name+"-port")),
		ContainerPort: port,
		Protocol:      "TCP",
	}}
	container.ImagePullPolicy = v1.PullAlways
	container.Resources = v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    cpu,
			v1.ResourceMemory: memory,
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    cpu,
			v1.ResourceMemory: memory,
		},
	}
	podTemplate := v1.PodTemplateSpec{}
	podTemplate.Labels = labels
	podTemplate.Name = name
	podTemplate.Spec = v1.PodSpec{
		Containers: []v1.Container{
			container,
		},
		ImagePullSecrets: []v1.LocalObjectReference{{Name: registrySecretName}},
	}
	deployment.Spec = appsV1.DeploymentSpec{
		Replicas: Int32(replicas),
		Selector: &metav1.LabelSelector{MatchLabels: labels},
		Template: podTemplate,
	}
	if _, err := d.client.AppsV1().Deployments(stormNs).Create(deployment); err != nil {
		return err
	}
	return nil
}

func (d *defaultK8sService) createLabel(name string) map[string]string {
	return map[string]string{"web": fmt.Sprintf("%s-service", strings.ToLower(name))}
}

func (d *defaultK8sService) createRegistrySecret() error {
	secret := &v1.Secret{}
	secret.Name = registrySecretName
	secret.Type = v1.SecretTypeDockerConfigJson
	data, err := d.dockerConfigJson()
	if err != nil {
		return err
	}
	secret.Data = map[string][]byte{
		v1.DockerConfigJsonKey: data,
	}
	c := d.client.CoreV1().Secrets(stormNs)
	if s, err := c.Get(registrySecretName, metav1.GetOptions{}); err == nil && s.Name != "" {
		if err := c.Delete(registrySecretName, nil); err != nil {
			return err
		}
	}
	if _, err := c.Create(secret); err != nil {
		return err
	}
	log.Println("secret created")
	return nil
}

func (d *defaultK8sService) GetLogs(appName string) (string, error) {
	selector := fmt.Sprintf("web=%s-service", strings.ToLower(appName))
	pods, err := d.getPodsBySelector(selector)
	if err != nil {
		return "", err
	}
	logs := ""
	for _, p := range pods {
		l := d.client.CoreV1().Pods(stormNs).GetLogs(p.Name, &v1.PodLogOptions{})
		r, err := l.Stream()
		if err != nil {
			log.Println("failed to get stream handle for ", p.Name)
			continue
		}
		data, err := ioutil.ReadAll(r)
		if err != nil {
			log.Println("failed to read stream data for ", p.Name)
			continue
		}
		logs += string(data)
	}
	return logs, nil
}

func (d *defaultK8sService) getPodsBySelector(selector string) ([]v1.Pod, error) {
	pods, err := d.client.CoreV1().Pods(stormNs).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func (d *defaultK8sService) getNodeIp(nodeName string) (string, error) {
	nd, err := d.client.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	log.Println("Node: ", nodeName, nd)
	addrs := nd.Status.Addresses
	nodeIp := ""
	for _, addr := range addrs {
		if addr.Type == v1.NodeExternalIP && addr.Address != "" {
			nodeIp = addr.Address
			break
		}
	}
	return nodeIp, nil
}

// dockerConfigJson returns a json rep of user's
// docker registry auth credentials.
func (d *defaultK8sService) dockerConfigJson() ([]byte, error) {
	type authData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Auth     string `json:"auth"`
	}
	username, password := d.cfg.Registry.Username, d.cfg.Registry.Password
	ad := authData{
		Username: username,
		Password: password,
	}
	type auths struct {
		Auths map[string]authData `json:"auths"`
	}
	usernamePassword := fmt.Sprintf("%s:%s", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(usernamePassword))
	ad.Auth = encoded
	a := &auths{Auths: map[string]authData{
		d.cfg.Registry.Url: ad,
	}}
	return json.Marshal(a)
}

func findAvailablePort() int {
	port := rand.Intn(59999)
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return port
	}
	if err := conn.Close(); err != nil { /*no-op*/
	}
	return findAvailablePort()
}

func Int32(i int32) *int32 {
	return &i
}

func truncString(s string) string {
	if len(s) < 15 {
		return s
	}
	return strings.TrimSpace(s[0:14])
}

