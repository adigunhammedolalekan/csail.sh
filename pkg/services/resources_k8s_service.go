package services

import (
	"fmt"
	"github.com/saas/hostgolang/pkg/res"
	"github.com/saas/hostgolang/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"strings"
)

type ResourcesService interface {
	DeployResource(app *types.App, envs []types.ResourceEnv, res res.Res) (*types.ResourceDeploymentResult, error)
	DeleteResource(app *types.App, name string) error
}

type defaultResourcesService struct {
	client *kubernetes.Clientset
}

func NewResourcesService(c *kubernetes.Clientset) ResourcesService {
	return &defaultResourcesService{client: c}
}

func (d *defaultResourcesService) DeployResource(app *types.App, resourcesEnvs []types.ResourceEnv, res res.Res) (*types.ResourceDeploymentResult, error) {
	serviceName := fmt.Sprintf("svc-%s-%s", res.Name(), app.AppName)
	svc, err := d.createResourceService(serviceName, res.Port())
	if err != nil {
		return nil, err
	}
	st, err := d.createResourceStatefulSet(app.AppName, svc, res)
	if err != nil {
		return nil, err
	}
	deployment, err := d.client.AppsV1().Deployments(stormNs).Get(app.AppName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	lbAddress := ""
	if lbAddr := svc.Spec.LoadBalancerIP; lbAddr == "" {
		for _, addr := range svc.Spec.ExternalIPs {
			if addr != "" {
				lbAddress = addr
				break
			}
		}
	}
	hostEnvKey := fmt.Sprintf("%s_HOST", strings.ToUpper(res.Name()))
	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		newEnvs := make([]v1.EnvVar, 0)
		for _, r := range resourcesEnvs {
			newEnvs = append(newEnvs, v1.EnvVar{Name: r.EnvKey, Value: r.EnvValue})
		}
		newEnvs = append(newEnvs, v1.EnvVar{Name: hostEnvKey, Value: lbAddress})
		deployment.Spec.Template.Spec.Containers[0].Env = newEnvs
		if _, err := d.client.AppsV1().Deployments(stormNs).Update(deployment); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := &types.ResourceDeploymentResult{
		ID:   string(st.UID),
		Ip:   lbAddress,
		Port: fmt.Sprintf("%d", res.Port()),
	}
	return result, nil
}

func (d *defaultResourcesService) createResourceService(serviceName string, targetPort int) (*v1.Service, error) {
	svcClient := d.client.CoreV1()
	if s, err := svcClient.Services(stormNs).Get(serviceName, metav1.GetOptions{}); err == nil && s.Name != "" {
		if err := svcClient.Services(stormNs).Delete(serviceName, &metav1.DeleteOptions{}); err != nil {
			return nil, err
		}
	}
	label := map[string]string{
		"res": serviceName,
	}

	svc := &v1.Service{}
	svc.Name = serviceName
	svc.Labels = label
	svc.Namespace = stormNs
	// servicePort := findAvailablePort()
	port := v1.ServicePort{
		Name:       truncString(fmt.Sprintf("%15s", serviceName+"-serviceport")),
		Protocol:   "TCP",
		Port:       int32(targetPort),
		TargetPort: intstr.FromInt(targetPort),
	}
	svc.Spec = v1.ServiceSpec{
		Ports:    []v1.ServicePort{port},
		Selector: label,
		Type:     v1.ServiceTypeLoadBalancer,
	}
	s, err := svcClient.Services(stormNs).Create(svc)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (d *defaultResourcesService) DeleteResource(app *types.App, name string) error {
	svcName := fmt.Sprintf("svc-%s-%s", name, app.AppName)
	stName := fmt.Sprintf("st-%s-%s", name, app.AppName)
	if err := d.client.CoreV1().Services(stormNs).Delete(svcName, &metav1.DeleteOptions{}); err != nil {
		return err
	}
	if err := d.client.AppsV1().StatefulSets(stormNs).Delete(stName, &metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func (d *defaultResourcesService) createResourceStatefulSet(appName string, svc *v1.Service, res res.Res) (*appsv1.StatefulSet, error) {
	name := fmt.Sprintf("st-%s-%s", res.Name(), appName)
	if st, err := d.client.AppsV1().StatefulSets(stormNs).Get(name, metav1.GetOptions{}); err == nil && st.Name != "" {
		if err := d.client.AppsV1().StatefulSets(stormNs).Delete(name, &metav1.DeleteOptions{}); err != nil {
			return nil, err
		}
	}
	cpu, err := resource.ParseQuantity(fmt.Sprintf("%fm", res.Quota().Cpu*1000))
	if err != nil {
		return nil, err
	}
	memory, err := resource.ParseQuantity(fmt.Sprintf("%fMi", res.Quota().Memory*1000))
	if err != nil {
		return nil, err
	}
	labels := svc.Labels
	st := &appsv1.StatefulSet{}
	st.Name = name
	st.Labels = labels
	st.Spec.ServiceName = svc.Name
	st.Spec.Selector = &metav1.LabelSelector{MatchLabels: labels}
	container := v1.Container{}
	container.Name = fmt.Sprintf("%s-cont", name)
	container.Image = res.Image()
	envs := make([]v1.EnvVar, 0)
	for k, v := range res.Envs() {
		envs = append(envs, v1.EnvVar{Name: k, Value: v})
	}
	if len(res.Args()) > 0 {
		args := res.Args()
		container.Args = args
	}
	container.Env = envs
	container.Resources = v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceMemory: memory,
			v1.ResourceCPU:    cpu,
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    cpu,
			v1.ResourceMemory: memory,
		},
	}
	container.Ports = []v1.ContainerPort{{
		Name:          truncString(fmt.Sprintf("%15s", name+"port")),
		ContainerPort: int32(res.Port()),
		Protocol:      "TCP",
	}}
	st.Spec.Replicas = Int32(1)
	template := v1.PodTemplateSpec{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				container,
			},
		},
	}
	template.Labels = labels
	st.Spec.Template = template

	// volume
	storageQuantity, err := resource.ParseQuantity(fmt.Sprintf("%fMi", res.Quota().StorageSize))
	if err != nil {
		return nil, err
	}
	pvc := v1.PersistentVolumeClaim{}
	pvc.Name = fmt.Sprintf("%s-%s-pvc", appName, res.Name())
	pvc.Spec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce}
	pvc.Spec.Resources = v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceStorage: storageQuantity,
		},
	}
	st.Spec.VolumeClaimTemplates = []v1.PersistentVolumeClaim{pvc}
	return d.client.AppsV1().StatefulSets(stormNs).Create(st)
}
