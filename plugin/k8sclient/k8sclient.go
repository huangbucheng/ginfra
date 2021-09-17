package k8sclient

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
)

type KubeClient struct {
	kubeconfigPath string
	cs             *kubernetes.Clientset
}

type JobRequest struct {
	Namespace               string
	JobName                 string
	Image                   string
	CpuRequest              string // 700m
	MemoryRequest           string // 512Mi
	CpuLimit                string
	MemoryLimit             string
	Mounts                  []VolumeMount
	Envs                    map[string]string
	TTLSecondsAfterFinished int32
}

type DeploymentRequest struct {
	Namespace     string
	DeployName    string
	Image         string
	CpuRequest    string // 700m
	MemoryRequest string // 512Mi
	CpuLimit      string
	MemoryLimit   string
	Mounts        []VolumeMount
	Envs          map[string]string
	Port          int32
}

type ServiceRequest struct {
	Namespace   string
	ServiceName string
	Port        int32
	Selector    map[string]string
}

type VolumeMount struct {
	Name      string
	MountPath string // /usr/local/service/runner/code
	HostPath  string // /tmp/arena/test
}

func (c *KubeClient) WithKubeConfig(kubeconfigPath string) error {
	// uses the current context in kubeconfig
	// path-to-kubeconfig -- for example, /root/.kube/config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	// creates the clientset
	c.cs, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}

// InClusterExp 在TKE中不能访问API Server
func (c *KubeClient) InClusterConfig() error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	c.cs, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func (c *KubeClient) WatchPod(ctx context.Context, namespace, podname string) error {
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := c.cs.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = c.cs.CoreV1().Pods(namespace).Get(
			context.TODO(), podname, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s not found in test namespace\n", podname)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			return err
		} else {
			fmt.Printf("Found %s pod in namespace\n", podname)
		}

		time.Sleep(10 * time.Second)
	}
}

func (c *KubeClient) GetJob(ctx context.Context, namespace, jobname string) (*batchv1.Job, error) {
	// access the API to list job pods
	//jobs, err := c.cs.BatchV1().Jobs(namespace).List(ctx, metav1.ListOptions{})
	//fmt.Printf("There are following jobs in the cluster:%v, err:%v\n", jobs, err)

	return c.cs.BatchV1().Jobs(namespace).Get(ctx, jobname, metav1.GetOptions{})
}

func (c *KubeClient) DeleteJob(ctx context.Context, namespace, jobname string) error {
	propagationPolicy := metav1.DeletePropagationBackground
	return c.cs.BatchV1().Jobs(namespace).Delete(ctx, jobname, metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy})
}

func (c *KubeClient) CreateFootballJob(ctx context.Context, req *JobRequest) (
	*batchv1.Job, error) {
	var volumes []v1.Volume
	var mounts []v1.VolumeMount
	for _, v := range req.Mounts {
		vm, v := getHostPathVolumeMount(v.Name, v.MountPath, v.HostPath)
		volumes = append(volumes, v)
		mounts = append(mounts, vm)
	}

	jobsClient := c.cs.BatchV1().Jobs(req.Namespace)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.JobName,
			Namespace: req.Namespace,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            pointer.Int32Ptr(1),
			TTLSecondsAfterFinished: pointer.Int32Ptr(req.TTLSecondsAfterFinished),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					RestartPolicy: "Never",
					Containers: []v1.Container{
						{
							Name:  "main",
							Image: req.Image,
							Env:   EnvToVars(req.Envs),
							//Command: []string{"sleep"},
							//Args:    []string{"10000"},
							SecurityContext: &v1.SecurityContext{
								Privileged:               pointer.BoolPtr(false),
								AllowPrivilegeEscalation: pointer.BoolPtr(false),
								ReadOnlyRootFilesystem:   pointer.BoolPtr(true),
								//RunAsNonRoot:             pointer.BoolPtr(true),
							},
							Resources: getResourceRequirements(
								getResourceList(req.CpuRequest, req.MemoryRequest),
								getResourceList(req.CpuRequest, req.MemoryRequest),
							),
							VolumeMounts: mounts,
						},
					},
					Volumes:          volumes,
					ImagePullSecrets: []v1.LocalObjectReference{{Name: "qcloudregistrykey"}},
				},
			},
		},
	}

	return jobsClient.Create(ctx, job, metav1.CreateOptions{})
}

func (c *KubeClient) CreateFootballDeployment(ctx context.Context, req *DeploymentRequest) (
	*appsv1.Deployment, error) {
	var volumes []v1.Volume
	var mounts []v1.VolumeMount
	for _, v := range req.Mounts {
		vm, v := getHostPathVolumeMount(v.Name, v.MountPath, v.HostPath)
		volumes = append(volumes, v)
		mounts = append(mounts, vm)
	}

	client := c.cs.AppsV1().Deployments(req.Namespace)
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.DeployName,
			Namespace: req.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"runner": req.DeployName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"runner": req.DeployName,
					},
				},
				Spec: v1.PodSpec{
					RestartPolicy: "Always",
					Containers: []v1.Container{
						{
							Name:  "main",
							Image: req.Image,
							Env:   EnvToVars(req.Envs),
							//Command: []string{"sleep"},
							//Args:    []string{"10000"},
							SecurityContext: &v1.SecurityContext{
								Privileged:               pointer.BoolPtr(false),
								AllowPrivilegeEscalation: pointer.BoolPtr(false),
								ReadOnlyRootFilesystem:   pointer.BoolPtr(true),
								//RunAsNonRoot:             pointer.BoolPtr(true),
							},
							Resources: getResourceRequirements(
								getResourceList(req.CpuRequest, req.MemoryRequest),
								getResourceList(req.CpuRequest, req.MemoryRequest),
							),
							VolumeMounts: mounts,
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: req.Port,
								},
							},
						},
					},
					Volumes:          volumes,
					ImagePullSecrets: []v1.LocalObjectReference{{Name: "qcloudregistrykey"}},
				},
			},
		},
	}

	return client.Create(ctx, deploy, metav1.CreateOptions{})
}

func (c *KubeClient) CreateFootballService(ctx context.Context, req *ServiceRequest) (
	*v1.Service, error) {
	client := c.cs.CoreV1().Services(req.Namespace)

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.ServiceName,
			Namespace: req.Namespace,
		},
		Spec: v1.ServiceSpec{
			Ports:     []v1.ServicePort{{Port: req.Port}},
			Selector:  req.Selector,
			ClusterIP: "",
		},
	}

	// Create Service
	return client.Create(ctx, service, metav1.CreateOptions{})
}

func (c *KubeClient) UpdateFootballService(ctx context.Context, namespace string, service *v1.Service) (
	*v1.Service, error) {
	client := c.cs.CoreV1().Services(namespace)

	// Update Service
	return client.Update(ctx, service, metav1.UpdateOptions{})
}

func (c *KubeClient) GetService(ctx context.Context, namespace, servicename string) (*v1.Service, error) {
	return c.cs.CoreV1().Services(namespace).Get(ctx, servicename, metav1.GetOptions{})
}

func (c *KubeClient) DeleteService(ctx context.Context, namespace, servicename string) error {
	return c.cs.CoreV1().Services(namespace).Delete(ctx, servicename, metav1.DeleteOptions{})
}

func (c *KubeClient) GetDeploy(ctx context.Context, namespace, deployname string) (*appsv1.Deployment, error) {
	return c.cs.AppsV1().Deployments(namespace).Get(ctx, deployname, metav1.GetOptions{})
}

func (c *KubeClient) DeleteDeploy(ctx context.Context, namespace, deployname string) error {
	propagationPolicy := metav1.DeletePropagationBackground
	return c.cs.AppsV1().Deployments(namespace).Delete(ctx, deployname, metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy})
}

// 返回指定的cpu、memory资源值
// 写法参考k8s见：
// https://github.com/kubernetes/kubernetes/blob/b3875556b0edf3b5eaea32c69678edcf4117d316/pkg/kubelet/cm/helpers_linux_test.go#L36-L53
func getResourceList(cpu, memory string) v1.ResourceList {
	res := v1.ResourceList{}
	if cpu != "" {
		res[v1.ResourceCPU] = resource.MustParse(cpu)
	}
	if memory != "" {
		res[v1.ResourceMemory] = resource.MustParse(memory)
	}
	return res
}

// 返回ResourceRequirements对象，详细见getResourceList函数注释
func getResourceRequirements(requests, limits v1.ResourceList) v1.ResourceRequirements {
	res := v1.ResourceRequirements{}
	res.Requests = requests
	res.Limits = limits
	return res
}

// 把对象转化成k8s所能接受的环境变量格式
func EnvToVars(envMap map[string]string) []v1.EnvVar {
	var envVars []v1.EnvVar
	for k, v := range envMap {
		envVar := v1.EnvVar{
			Name:  k,
			Value: v,
		}
		envVars = append(envVars, envVar)
	}
	return envVars
}

func getHostPathVolumeMount(name, mountpath, hostpath string) (v1.VolumeMount, v1.Volume) {
	vm := v1.VolumeMount{
		Name:      name,
		MountPath: mountpath,
	}
	v := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: hostpath,
			},
		},
	}
	return vm, v
}
