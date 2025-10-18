package aws

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/scttfrdmn/nrp-aws-kip/pkg/config"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	corev1 "k8s.io/api/core/v1"
)

// Client handles interactions with AWS services
type Client struct {
	config    *config.Config
	ec2Client *ec2.Client
	instances map[string]*Instance // namespace/name -> Instance
}

// Instance represents an EC2 instance mapped to a Kubernetes pod
type Instance struct {
	InstanceID string
	Pod        *corev1.Pod
	State      types.InstanceState
}

// NewClient creates a new AWS client
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Client{
		config:    cfg,
		ec2Client: ec2.NewFromConfig(awsCfg),
		instances: make(map[string]*Instance),
	}, nil
}

// CreateInstance creates an EC2 instance for the given pod
func (c *Client) CreateInstance(ctx context.Context, pod *corev1.Pod) error {
	// Determine instance type based on pod resource requests
	instanceType := c.selectInstanceType(pod)

	// Build user data script
	userData := c.buildUserData(pod)

	// Prepare tags
	tags := c.buildTags(pod)

	// Prepare launch specification
	runInput := &ec2.RunInstancesInput{
		ImageId:      aws.String(c.selectAMI()),
		InstanceType: types.InstanceType(instanceType),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		SubnetId:     aws.String(c.config.AWS.SubnetIDs[0]), // TODO: implement subnet selection
		SecurityGroupIds: c.config.AWS.SecurityGroupIDs,
		KeyName:      aws.String(c.config.AWS.KeyName),
		UserData:     aws.String(userData),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags:         tags,
			},
		},
	}

	if c.config.AWS.IAMInstanceProfile != "" {
		runInput.IamInstanceProfile = &types.IamInstanceProfileSpecification{
			Name: aws.String(c.config.AWS.IAMInstanceProfile),
		}
	}

	// Handle spot instances if configured
	if c.config.AWS.UseSpotInstances {
		runInput.InstanceMarketOptions = &types.InstanceMarketOptionsRequest{
			MarketType: types.MarketTypeSpot,
			SpotOptions: &types.SpotMarketOptions{
				SpotInstanceType: types.SpotInstanceTypeOneTime,
			},
		}
		if c.config.AWS.SpotMaxPrice != "" {
			runInput.InstanceMarketOptions.SpotOptions.MaxPrice = aws.String(c.config.AWS.SpotMaxPrice)
		}
	}

	// Launch the instance
	result, err := c.ec2Client.RunInstances(ctx, runInput)
	if err != nil {
		return fmt.Errorf("failed to launch EC2 instance: %w", err)
	}

	if len(result.Instances) == 0 {
		return fmt.Errorf("no instances were launched")
	}

	instance := &Instance{
		InstanceID: *result.Instances[0].InstanceId,
		Pod:        pod.DeepCopy(),
		State:      *result.Instances[0].State,
	}

	key := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
	c.instances[key] = instance

	return nil
}

// DeleteInstance terminates the EC2 instance for the given pod
func (c *Client) DeleteInstance(ctx context.Context, pod *corev1.Pod) error {
	key := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
	instance, exists := c.instances[key]
	if !exists {
		return nil // Already deleted or never existed
	}

	_, err := c.ec2Client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{instance.InstanceID},
	})
	if err != nil {
		return fmt.Errorf("failed to terminate instance: %w", err)
	}

	delete(c.instances, key)
	return nil
}

// GetPodStatus retrieves the pod status from AWS
func (c *Client) GetPodStatus(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	key := fmt.Sprintf("%s/%s", namespace, name)
	instance, exists := c.instances[key]
	if !exists {
		return nil, fmt.Errorf("pod not found")
	}

	// Query current instance state
	result, err := c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instance.InstanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance not found")
	}

	ec2Instance := result.Reservations[0].Instances[0]
	instance.State = *ec2Instance.State

	// Update pod status based on instance state
	pod := instance.Pod.DeepCopy()
	pod.Status.Phase = c.mapInstanceStateToPodPhase(ec2Instance.State.Name)

	if ec2Instance.PrivateIpAddress != nil {
		pod.Status.PodIP = *ec2Instance.PrivateIpAddress
		pod.Status.HostIP = *ec2Instance.PrivateIpAddress
	}

	// Update container statuses
	for i := range pod.Spec.Containers {
		containerStatus := corev1.ContainerStatus{
			Name:  pod.Spec.Containers[i].Name,
			Ready: pod.Status.Phase == corev1.PodRunning,
			State: corev1.ContainerState{},
		}

		if pod.Status.Phase == corev1.PodRunning {
			containerStatus.State.Running = &corev1.ContainerStateRunning{}
		}

		pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, containerStatus)
	}

	return pod, nil
}

// ListPods returns all pods managed by this provider
func (c *Client) ListPods(ctx context.Context) ([]*corev1.Pod, error) {
	pods := make([]*corev1.Pod, 0, len(c.instances))

	for key := range c.instances {
		parts := splitKey(key)
		if len(parts) != 2 {
			continue
		}

		pod, err := c.GetPodStatus(ctx, parts[0], parts[1])
		if err != nil {
			continue
		}
		pods = append(pods, pod)
	}

	return pods, nil
}

// GetContainerLogs retrieves logs for a container (stub implementation)
func (c *Client) GetContainerLogs(ctx context.Context, namespace, podName, containerName string, opts api.ContainerLogOpts) (io.ReadCloser, error) {
	// TODO: Implement log retrieval via CloudWatch Logs or SSM
	return nil, fmt.Errorf("container logs not yet implemented")
}

// RunInContainer executes a command in a container (stub implementation)
func (c *Client) RunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string, attach api.AttachIO) error {
	// TODO: Implement via SSM Session Manager
	return fmt.Errorf("exec not yet implemented")
}

// Helper functions

func (c *Client) selectInstanceType(pod *corev1.Pod) string {
	// Simple instance type selection based on resource requests
	// TODO: Implement more sophisticated selection logic

	var cpuMillis int64
	var memoryBytes int64

	for _, container := range pod.Spec.Containers {
		if cpu := container.Resources.Requests.Cpu(); cpu != nil {
			cpuMillis += cpu.MilliValue()
		}
		if mem := container.Resources.Requests.Memory(); mem != nil {
			memoryBytes += mem.Value()
		}
	}

	// Simple mapping (this should be more sophisticated in production)
	if cpuMillis < 2000 && memoryBytes < 4*1024*1024*1024 {
		return "t3.medium"
	} else if cpuMillis < 4000 && memoryBytes < 8*1024*1024*1024 {
		return "t3.large"
	} else if cpuMillis < 8000 && memoryBytes < 16*1024*1024*1024 {
		return "t3.xlarge"
	}

	return "t3.2xlarge"
}

func (c *Client) selectAMI() string {
	// TODO: Implement AMI selection based on region and requirements
	// For now, return a placeholder
	return "ami-0c55b159cbfafe1f0" // Amazon Linux 2 in us-east-1 (example)
}

func (c *Client) buildUserData(pod *corev1.Pod) string {
	// TODO: Build appropriate user data to run the pod's containers
	// This could use Docker, containerd, or other container runtimes
	return `#!/bin/bash
echo "NRP AWS KIP Instance"
# TODO: Install container runtime and start containers
`
}

func (c *Client) buildTags(pod *corev1.Pod) []types.Tag {
	tags := []types.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String(fmt.Sprintf("%s-%s", pod.Namespace, pod.Name)),
		},
		{
			Key:   aws.String("nrp-aws-kip/namespace"),
			Value: aws.String(pod.Namespace),
		},
		{
			Key:   aws.String("nrp-aws-kip/pod"),
			Value: aws.String(pod.Name),
		},
	}

	// Add configured tags
	for k, v := range c.config.AWS.Tags {
		tags = append(tags, types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	return tags
}

func (c *Client) mapInstanceStateToPodPhase(state types.InstanceStateName) corev1.PodPhase {
	switch state {
	case types.InstanceStateNamePending:
		return corev1.PodPending
	case types.InstanceStateNameRunning:
		return corev1.PodRunning
	case types.InstanceStateNameStopping, types.InstanceStateNameStopped:
		return corev1.PodSucceeded
	case types.InstanceStateNameTerminated:
		return corev1.PodSucceeded
	default:
		return corev1.PodUnknown
	}
}

func splitKey(key string) []string {
	result := make([]string, 0, 2)
	start := 0

	for i, char := range key {
		if char == '/' {
			result = append(result, key[start:i])
			start = i + 1
		}
	}

	if start < len(key) {
		result = append(result, key[start:])
	}

	return result
}
