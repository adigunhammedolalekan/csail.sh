package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/saas/hostgolang/pkg/types"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/retry"
	"net/http"
	"os"

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
const validPortRange = 65535
const letsEncryptProd = "https://acme-v02.api.letsencrypt.org/directory"

//go:generate mockgen -destination=../mocks/k8s_service_mock.go -package=mocks github.com/saas/hostgolang/pkg/services K8sService
type K8sService interface {
	DeployService(opt *types.CreateDeploymentOpts) (*types.DeploymentResult, error)
	GetLogs(appName string) (string, error)
	UpdateEnvs(appName string, envs []types.Environment) error
	ScaleApp(appName string, replicas int) error
	ListRunningPods(appName string) ([]types.Instance, error)
	AddDomain(appName, domain string) error
	RemoveDomain(appName, domain string) error
	PodExec(appName, resName string, cmds []string) (string, error)
}

type defaultK8sService struct {
	client *kubernetes.Clientset
	dy     dynamic.Interface
	cfg    *config.Config
	restConfig *rest.Config
}

func NewK8sService(client *kubernetes.Clientset, dy dynamic.Interface,
	config *config.Config, restConfig *rest.Config) K8sService {
	d := &defaultK8sService{client: client, cfg: config, dy: dy, restConfig: restConfig}
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
	var serviceType = v1.ServiceTypeClusterIP
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

	accessAddr := fmt.Sprintf("%s.%s", svc.Name, stormNs)
	for _, p := range ports {
		if targetPort := p.Port; targetPort != 0 {
			addr := fmt.Sprintf("http://%s:%d", accessAddr, targetPort)
			accessAddr = addr
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
			if strings.ToLower(strings.TrimSpace(pe.Name)) == "port" {
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
	servicePort := findPort()
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

func (d *defaultK8sService) AddDomain(appName, domain string) error {
	name := fmt.Sprintf("%s-ingress", domain)
	in, err := d.client.ExtensionsV1beta1().Ingresses(stormNs).Get(name, metav1.GetOptions{})
	if err == nil && in.Name != "" {
		if err := d.deleteIngress(appName, domain); err != nil {
			return err
		}
	}
	return d.createIngress(appName, domain)
}

func (d *defaultK8sService) RemoveDomain(appName, domain string) error {
	return d.deleteIngress(appName, domain)
}

func (d *defaultK8sService) deleteIngress(appName, domain string) error {
	name := fmt.Sprintf("%s-ingress", domain)
	c := d.client.ExtensionsV1beta1().Ingresses(stormNs)
	if err := c.Delete(name, &metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func (d *defaultK8sService) createIngress(appName, domain string) error {
	c := d.client.CoreV1().Services(stormNs)
	svc, err := c.Get(appName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if err := d.provisionCertificate(appName, domain); err != nil {
		log.Println("ProvisionCertError ", err)
		return err
	}
	issuerName := fmt.Sprintf("%s-issuer", domain)
	ingressName := fmt.Sprintf("%s-ingress", domain)
	secretName := fmt.Sprintf("%s-sec", domain)
	ingress := &v1beta1.Ingress{}
	ingress.Name = ingressName
	ingress.Annotations = map[string]string{
		"kubernetes.io/ingress.class":                 "nginx",
		"nginx.ingress.kubernetes.io/proxy-body-size": "10000m",
		"cert-manager.io/cluster-issuer":              issuerName,
	}

	ports := svc.Spec.Ports
	var port int32 = 0
	if len(ports) > 0 {
		port = ports[0].Port
	}
	// define backend
	backend := v1beta1.IngressBackend{
		ServiceName: svc.Name,
		ServicePort: intstr.FromInt(int(port)),
	}
	ingressRule := &v1beta1.IngressRule{Host: domain}
	ingressRule.HTTP = &v1beta1.HTTPIngressRuleValue{
		Paths: []v1beta1.HTTPIngressPath{{Backend: backend}},
	}
	ingress.Spec.Backend = &backend
	ingress.Spec.TLS = []v1beta1.IngressTLS{
		{Hosts: []string{domain}, SecretName: secretName},
	}
	ingress.Spec.Rules = []v1beta1.IngressRule{*ingressRule}
	if _, err := d.client.ExtensionsV1beta1().Ingresses(stormNs).Create(ingress); err != nil {
		return err
	}
	return nil
}

func (d *defaultK8sService) provisionCertificate(appName, domain string) error {
	ns := "cert-manager"
	issuerResources := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1alpha2", Resource: "clusterissuer"}
	secretName := fmt.Sprintf("%s-key", domain)
	issuerName := fmt.Sprintf("%s-issuer", domain)
	certIssuer := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cert-manager.io/v1alpha2",
			"kind":       "ClusterIssuer",
			"metadata": map[string]interface{}{
				"name":      issuerName,
				"namespace": ns,
			},
			"spec": map[string]interface{}{
				"acme": map[string]interface{}{
					"email":  "adigunhammed.lekan@gmail.com",
					"server": letsEncryptProd,
					"privateKeySecretRef": map[string]interface{}{
						"name": secretName,
					},
					"solvers": []map[string]interface{}{
						{"http01": map[string]interface{}{
							"ingress": map[string]interface{}{
								"class": "nginx",
							},
						}},
					},
				},
			},
		},
	}
	_, err := d.dy.Resource(issuerResources).Create(certIssuer, metav1.CreateOptions{})
	if err != nil {
		log.Println("IssuerError ", err)
		return err
	}
	return d.createCertificate(issuerName, appName, domain)
}

func (d *defaultK8sService) createCertificate(issuerName, appName, domain string) error {
	secretName := fmt.Sprintf("%s-sec", domain)
	name := fmt.Sprintf("%s-cert", domain)
	ns := "cert-manager"

	certResource := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1alpha2", Resource: "certificate"}
	cert := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cert-manager.io/v1alpha2",
			"kind":       "Certificate",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": ns,
			},
			"spec": map[string]interface{}{
				"secretName": secretName,
				"issuerRef": map[string]interface{}{
					"name": issuerName,
					"kind": "ClusterIssuer",
				},
				"commonName": domain,
				"dnsNames":   []string{domain},
			},
		},
	}
	_, err := d.dy.Resource(certResource).Create(cert, metav1.CreateOptions{})
	return err
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

func findPort() int {
	port := rand.Intn(validPortRange)
	if port < 1025 {
		return findPort()
	}
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return port
	}
	if err := conn.Close(); err != nil { /*no-op*/}
	return findPort()
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

func (d *defaultK8sService) PodExec(appName, resName string, cmds []string) (string, error) {
	log.Printf("Execing command for App=%s, Res=%s, Cmds=%s", appName, resName, cmds)
	selector := fmt.Sprintf("res=svc-%s-%s", resName, appName)
	pods, err := d.getPodsBySelector(selector)
	if err != nil {
		return "", err
	}
	if len(pods) <= 0 {
		return "", errors.New("no pod found")
	}
	pod := pods[0]
	containers := pod.Spec.Containers
	container := ""
	if len(containers) > 0 {
		container = containers[0].Name
	}
	log.Printf("Pod=%s <==> Container=%s", pod.Name, container)
	c := d.client.CoreV1().RESTClient()

	req := c.Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		Param("container", container).
		Param("stdout", "true").
		Param("stderr", "true")

	for _, cmd := range cmds {
		req.Param("command", cmd)
	}

	executor, err := remotecommand.NewSPDYExecutor(d.restConfig, http.MethodPost, req.URL())
	if err != nil {
		return "", err
	}
	out := &bytes.Buffer{}
	os.Stdout.Sync()
	os.Stderr.Sync()
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:              nil,
		Stdout:             out,
		Stderr:             os.Stderr,
		Tty:                false,
		TerminalSizeQueue:  nil,
	})
	if err != nil {
		return "", err
	}
	return out.String(), nil
}