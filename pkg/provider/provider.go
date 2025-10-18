package provider

import (
	"context"
	"io"

	"github.com/scttfrdmn/nrp-aws-kip/internal/aws"
	"github.com/scttfrdmn/nrp-aws-kip/pkg/config"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSProvider implements the Virtual Kubelet provider interface for AWS
type AWSProvider struct {
	config      *config.Config
	awsClient   *aws.Client
	nodeName    string
	region      string
	operatingOS string
}

// NewAWSProvider creates a new AWS provider instance
func NewAWSProvider(ctx context.Context, cfg *config.Config) (*AWSProvider, error) {
	awsClient, err := aws.NewClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &AWSProvider{
		config:      cfg,
		awsClient:   awsClient,
		nodeName:    cfg.Node.Name,
		region:      cfg.AWS.Region,
		operatingOS: cfg.Node.OperatingSystem,
	}, nil
}

// CreatePod takes a Kubernetes Pod and deploys it within the provider
func (p *AWSProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	return p.awsClient.CreateInstance(ctx, pod)
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider
func (p *AWSProvider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	// For now, we'll recreate the instance
	// In production, you might want to handle updates more gracefully
	if err := p.DeletePod(ctx, pod); err != nil {
		return err
	}
	return p.CreatePod(ctx, pod)
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider
func (p *AWSProvider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	return p.awsClient.DeleteInstance(ctx, pod)
}

// GetPod retrieves a pod by name from the provider
func (p *AWSProvider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	return p.awsClient.GetPodStatus(ctx, namespace, name)
}

// GetPodStatus retrieves the status of a pod by name from the provider
func (p *AWSProvider) GetPodStatus(ctx context.Context, namespace, name string) (*corev1.PodStatus, error) {
	pod, err := p.awsClient.GetPodStatus(ctx, namespace, name)
	if err != nil {
		return nil, err
	}
	return &pod.Status, nil
}

// GetPods retrieves a list of all pods running on the provider
func (p *AWSProvider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
	return p.awsClient.ListPods(ctx)
}

// GetContainerLogs retrieves the logs of a container by name from the provider
func (p *AWSProvider) GetContainerLogs(ctx context.Context, namespace, podName, containerName string, opts api.ContainerLogOpts) (io.ReadCloser, error) {
	return p.awsClient.GetContainerLogs(ctx, namespace, podName, containerName, opts)
}

// RunInContainer executes a command in a container in the pod, copying data
// between in/out/err and the container's stdin/stdout/stderr
func (p *AWSProvider) RunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string, attach api.AttachIO) error {
	return p.awsClient.RunInContainer(ctx, namespace, podName, containerName, cmd, attach)
}

// ConfigureNode enables a provider to configure the node object that
// will be used for Kubernetes
func (p *AWSProvider) ConfigureNode(ctx context.Context, node *corev1.Node) {
	node.Status.Capacity = p.capacity()
	node.Status.Allocatable = p.capacity()
	node.Status.Conditions = p.nodeConditions()
	node.Status.Addresses = p.nodeAddresses()
	node.Status.DaemonEndpoints = p.nodeDaemonEndpoints()
	node.Status.NodeInfo = p.nodeInfo()
}

// capacity returns the capacity of the virtual node
func (p *AWSProvider) capacity() corev1.ResourceList {
	return corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(p.config.Node.CPU),
		corev1.ResourceMemory: resource.MustParse(p.config.Node.Memory),
		corev1.ResourcePods:   resource.MustParse(p.config.Node.Pods),
	}
}

// nodeConditions returns a list of conditions (Ready, OutOfDisk, etc), for the node
func (p *AWSProvider) nodeConditions() []corev1.NodeCondition {
	return []corev1.NodeCondition{
		{
			Type:               corev1.NodeReady,
			Status:             corev1.ConditionTrue,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletReady",
			Message:            "NRP AWS burst node is ready",
		},
		{
			Type:               corev1.NodeMemoryPressure,
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasSufficientMemory",
			Message:            "kubelet has sufficient memory available",
		},
		{
			Type:               corev1.NodeDiskPressure,
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasNoDiskPressure",
			Message:            "kubelet has no disk pressure",
		},
		{
			Type:               corev1.NodePIDPressure,
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasSufficientPID",
			Message:            "kubelet has sufficient PID available",
		},
		{
			Type:               corev1.NodeNetworkUnavailable,
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "RouteCreated",
			Message:            "RouteController created a route",
		},
	}
}

// nodeAddresses returns a list of addresses for the node
func (p *AWSProvider) nodeAddresses() []corev1.NodeAddress {
	return []corev1.NodeAddress{
		{
			Type:    corev1.NodeHostName,
			Address: p.nodeName,
		},
	}
}

// nodeDaemonEndpoints returns NodeDaemonEndpoints
func (p *AWSProvider) nodeDaemonEndpoints() corev1.NodeDaemonEndpoints {
	return corev1.NodeDaemonEndpoints{
		KubeletEndpoint: corev1.DaemonEndpoint{
			Port: 10250,
		},
	}
}

// nodeInfo returns a NodeSystemInfo object
func (p *AWSProvider) nodeInfo() corev1.NodeSystemInfo {
	return corev1.NodeSystemInfo{
		MachineID:               "nrp-aws-kip",
		SystemUUID:              "nrp-aws-kip",
		BootID:                  "nrp-aws-kip",
		KernelVersion:           "5.10.0",
		OSImage:                 "Amazon Linux 2",
		ContainerRuntimeVersion: "aws://1.0.0",
		KubeletVersion:          "v1.28.0",
		KubeProxyVersion:        "v1.28.0",
		OperatingSystem:         p.operatingOS,
		Architecture:            "amd64",
	}
}

// NotifyPods instructs the notifier to call the passed in function when
// the pod status changes
func (p *AWSProvider) NotifyPods(ctx context.Context, notifierCb func(*corev1.Pod)) {
	// TODO: Implement pod status change notification
	// This would typically involve polling EC2 instance states
	// or using EventBridge/CloudWatch Events
}
