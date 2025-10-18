# NRP-AWS-KIP Architecture

## Overview

NRP-AWS-KIP (National Research Platform AWS Kubernetes Instance Provider) enables cloud bursting from NRP's Kubernetes clusters to AWS. It implements the Virtual Kubelet pattern to seamlessly extend cluster capacity.

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────┐
│ NRP Kubernetes Cluster                                   │
│                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Regular Node │  │ Regular Node │  │ Virtual Node │  │
│  │              │  │              │  │ (AWS KIP)    │  │
│  └──────────────┘  └──────────────┘  └──────┬───────┘  │
│                                              │           │
└──────────────────────────────────────────────┼───────────┘
                                               │
                            VPN / Direct Connect / Internet
                                               │
┌──────────────────────────────────────────────┼───────────┐
│ AWS Account                                  │           │
│                                              │           │
│  ┌───────────────────────────────────────────▼────────┐ │
│  │ VPC                                                 │ │
│  │                                                     │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  │ │
│  │  │ EC2        │  │ EC2        │  │ EC2        │  │ │
│  │  │ Instance   │  │ Instance   │  │ Instance   │  │ │
│  │  │ (Pod A)    │  │ (Pod B)    │  │ (Pod C)    │  │ │
│  │  └────────────┘  └────────────┘  └────────────┘  │ │
│  │                                                     │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                           │
│  ┌─────────────────────────────────────────────────────┐ │
│  │ CloudWatch Logs                                     │ │
│  └─────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────┘
```

## Components

### 1. Virtual Kubelet Provider

The core component that implements the Virtual Kubelet provider interface:

- **Location**: `pkg/provider/provider.go`
- **Responsibilities**:
  - Registers as a virtual node in the K8s cluster
  - Implements pod lifecycle operations (Create, Update, Delete, Get)
  - Reports node capacity and status
  - Handles pod scheduling decisions

### 2. AWS Client

Manages interactions with AWS services:

- **Location**: `internal/aws/client.go`
- **Responsibilities**:
  - EC2 instance lifecycle management
  - Instance type selection based on pod resources
  - Networking and security group management
  - Spot instance handling
  - Instance tagging and tracking

### 3. Configuration Manager

Handles application configuration:

- **Location**: `pkg/config/config.go`
- **Responsibilities**:
  - Parse and validate configuration files
  - Manage AWS credentials and settings
  - Define resource limits and quotas
  - Configure node properties

## Pod to EC2 Mapping

### Instance Selection Algorithm

```go
Pod Resources → Instance Type Selection
├── CPU Requests
├── Memory Requests
├── GPU Requirements (future)
└── Custom Annotations

Example:
- 1 CPU, 2GB RAM    → t3.medium
- 2 CPUs, 4GB RAM   → t3.large
- 4 CPUs, 8GB RAM   → t3.xlarge
- 8+ CPUs, 16GB RAM → t3.2xlarge
```

### Container Execution

Currently, EC2 instances are launched with user data scripts. Future improvements:

1. **Container Runtime**: Install Docker/containerd on instances
2. **Image Management**: Pull container images from registries
3. **Execution**: Run pod containers on instances
4. **Monitoring**: Forward logs to CloudWatch

## Networking

### Connectivity Options

1. **VPN Connection**
   - Site-to-site VPN between NRP and AWS VPC
   - Encrypted tunnel
   - Lower bandwidth, higher latency

2. **AWS Direct Connect** (Recommended)
   - Dedicated network connection
   - Higher bandwidth, lower latency
   - More reliable for production workloads

3. **Public Internet**
   - Pods communicate via public IPs
   - Less secure, higher latency
   - Suitable for testing only

### Security

- Security groups restrict traffic to NRP CIDR blocks
- IAM instance profiles provide AWS API access
- Network ACLs provide subnet-level protection

## Scheduling

### Taints and Tolerations

The virtual node is tainted to prevent accidental scheduling:

```yaml
taints:
  - key: nrp.org/burst
    value: "true"
    effect: NoSchedule
```

Pods must include matching tolerations:

```yaml
tolerations:
  - key: nrp.org/burst
    operator: Equal
    value: "true"
    effect: NoSchedule
```

### Node Selectors

Pods can target the burst node explicitly:

```yaml
nodeSelector:
  nrp.org/burst-node: "true"
```

## Cost Management

### Tracking

- EC2 instances tagged with pod namespace/name
- CloudWatch metrics for instance counts
- Cost allocation tags for billing

### Optimization

1. **Spot Instances**: Use spot instances for cost savings (configurable)
2. **Right-sizing**: Match instance types to pod requirements
3. **Auto-cleanup**: Terminate instances when pods are deleted
4. **Limits**: Configure max instances and cost limits

## Monitoring and Observability

### Metrics

- Prometheus metrics exposed on `:10255/metrics`
- Instance counts, states, and resource usage
- Pod creation/deletion rates

### Logging

- Application logs (JSON format)
- CloudWatch Logs for instance logs
- Structured logging with context

### Health Checks

- Liveness probe: `/healthz`
- Readiness probe: `/healthz`

## Scalability

### Limits

Configurable limits to prevent runaway costs:

```yaml
limits:
  maxPods: 500       # Maximum pods on this virtual node
  maxInstances: 100  # Maximum EC2 instances
  costLimit: 1000.00 # Optional monthly cost limit in USD
```

### Concurrency

- Multiple pods can be created/deleted concurrently
- AWS API rate limits are respected
- Instance state is tracked in-memory (future: persistent store)

## Future Enhancements

1. **Persistent State**: Use DynamoDB or etcd for state persistence
2. **Advanced Scheduling**: Consider availability zones, instance placement
3. **GPU Support**: Map GPU requests to GPU instance types
4. **Auto-scaling**: Scale based on pending pods or cluster metrics
5. **Cost Prediction**: Estimate costs before scheduling
6. **Multi-region**: Support bursting to multiple AWS regions
7. **Hybrid Clouds**: Support Azure, GCP alongside AWS
8. **Container Runtime**: Full container runtime integration
9. **Network Policies**: Implement K8s network policies in AWS
10. **Storage**: EBS volume management for persistent storage

## Development Workflow

```
1. Developer modifies code
2. Run tests: make test
3. Build binary: make build
4. Build Docker image: make docker-build
5. Deploy to K8s: kubectl apply -f deployments/kubernetes/
6. Monitor logs: kubectl logs -n nrp-aws-kip -l app=nrp-aws-kip
```

## Deployment Models

### Model 1: Centralized

Single NRP-AWS-KIP deployment manages all burst workloads:

- Easier to manage
- Single point of failure
- Resource contention

### Model 2: Per-Team

Each research team has their own NRP-AWS-KIP instance:

- Isolated AWS accounts/VPCs
- Independent cost tracking
- More complex management

### Model 3: Per-Workload-Type

Separate instances for different workload types:

- GPU workloads → GPU-optimized burst
- Batch jobs → Spot instances
- Interactive → On-demand instances

## Security Considerations

1. **AWS Credentials**: Use IAM roles or IRSA (IAM Roles for Service Accounts)
2. **Network Isolation**: VPC isolation, security groups
3. **Pod Security**: Enforce pod security policies
4. **Secrets Management**: Use Kubernetes secrets or AWS Secrets Manager
5. **Audit Logging**: CloudTrail for AWS API calls
6. **Compliance**: Ensure compliance with institutional policies
