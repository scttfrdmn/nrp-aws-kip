# Frequently Asked Questions (FAQ)

## General Questions

### What is NRP-AWS-KIP?

NRP-AWS-KIP is a cloud bursting solution that allows the National Research Platform's Kubernetes clusters to dynamically extend their capacity by launching workloads on AWS EC2 instances. It implements the Virtual Kubelet pattern to act as a bridge between Kubernetes and AWS.

### Why would I use cloud bursting?

Cloud bursting is useful when:
- Your on-premises cluster is at capacity
- You have occasional high-compute workloads (batch jobs, simulations)
- You need access to specialized hardware (GPUs, high-memory instances)
- You want to avoid over-provisioning on-premises infrastructure

### How much does it cost?

Costs depend on:
- EC2 instance types used
- Duration workloads run
- Data transfer between NRP and AWS
- Whether you use spot vs on-demand instances

Spot instances can save 50-90% compared to on-demand pricing.

## Technical Questions

### How does it work?

1. NRP-AWS-KIP registers as a virtual node in your Kubernetes cluster
2. When you schedule a pod with appropriate tolerations/selectors
3. Kubernetes schedules it to the virtual node
4. NRP-AWS-KIP creates an EC2 instance matching the pod's resource requirements
5. The instance runs your workload
6. When the pod is deleted, the instance is terminated

### What instance types are supported?

Currently, the instance selection is basic (t3 family). The implementation can be extended to support:
- Compute-optimized (c5, c6i)
- Memory-optimized (r5, r6i)
- GPU instances (p3, p4, g4dn)
- Any EC2 instance type

### Can I use spot instances?

Yes! Set `useSpotInstances: true` in the configuration. Spot instances can significantly reduce costs but may be interrupted by AWS.

### How do I ensure my pods use burst capacity?

Add a taint toleration and node selector to your pod spec:

```yaml
spec:
  nodeSelector:
    nrp.org/burst-node: "true"
  tolerations:
    - key: nrp.org/burst
      operator: Equal
      value: "true"
      effect: NoSchedule
```

### Does it support persistent storage?

Not yet. Current implementation is for stateless workloads. EBS volume support is planned for future releases.

### Can pods communicate with on-premises services?

Yes, if you have network connectivity (VPN or Direct Connect) between NRP and AWS. The security groups must allow the necessary traffic.

### How are container images handled?

Pods should use public container images or images from registries accessible from AWS. Private registries require proper authentication setup on the EC2 instances.

### Can I SSH into the instances?

If you configure an SSH key pair in Terraform and allow SSH in the security groups, yes. However, this is primarily for debugging.

### How do logs work?

Currently, container logs are not automatically collected. CloudWatch Logs integration is planned for future releases. For now, you can:
- Write logs to stdout/stderr and they'll be on the instance
- Use a logging sidecar to ship logs
- Write logs to a centralized logging system

### Can I exec into pods?

Not yet. This requires SSM Session Manager integration, which is planned for future releases.

## Networking Questions

### What networking options are available?

1. **VPN**: Site-to-site VPN between NRP and AWS VPC
2. **Direct Connect**: Dedicated connection (recommended for production)
3. **Public Internet**: Pods use public IPs (not recommended)

### Do I need a VPN?

For production use, yes. VPN or Direct Connect provides secure, reliable connectivity between NRP and AWS.

### What about DNS resolution?

- Kubernetes CoreDNS works for cluster services
- AWS instances can use AWS DNS or custom DNS servers
- Cross-cluster DNS requires additional configuration

### Are there bandwidth costs?

AWS charges for data transfer out to the internet. Data transfer within the same region is free. If using Direct Connect, pricing varies. Check AWS pricing for details.

## Security Questions

### Is this secure?

Security depends on your configuration:
- Use VPN/Direct Connect for network isolation
- Properly configure security groups
- Use IAM roles instead of access keys
- Follow AWS and Kubernetes security best practices

### How are AWS credentials managed?

Options:
1. **IAM Instance Profile** (recommended for the controller pod)
2. **IRSA** (IAM Roles for Service Accounts) - best practice
3. **Access Keys** (less secure, use only for testing)

### Can I limit what instances can do?

Yes, configure the IAM role attached to instances with minimal required permissions.

### Are instances isolated?

Instances run in your VPC with your security groups. They're isolated from other AWS customers but not from each other by default. Use security groups and network ACLs for additional isolation.

## Cost Questions

### How do I track costs?

- Use AWS Cost Explorer with tag filters
- Instances are tagged with pod namespace/name
- Set up AWS budgets with alerts
- Configure `costLimit` in the config

### How can I reduce costs?

1. Use spot instances (`useSpotInstances: true`)
2. Right-size instance types for workloads
3. Ensure pods terminate promptly when done
4. Use Reserved Instances or Savings Plans for predictable workloads
5. Monitor and clean up orphaned resources

### What if I exceed my budget?

The `costLimit` setting is informational only. You should:
- Set up AWS Budgets with email/SNS alerts
- Use AWS Cost Anomaly Detection
- Regularly review Cost Explorer
- Consider Service Control Policies (SCPs) for hard limits

## Troubleshooting

### My pods are stuck in Pending

Check:
1. Virtual node is registered: `kubectl get nodes`
2. Pod has correct tolerations and node selector
3. Virtual node has capacity: `kubectl describe node nrp-aws-burst`
4. NRP-AWS-KIP logs: `kubectl logs -n nrp-aws-kip -l app=nrp-aws-kip`

### Instances are created but pods don't run

Current implementation launches instances but doesn't fully run containers. Container runtime integration is needed (planned enhancement).

### AWS API errors

Check:
- AWS credentials are valid
- IAM permissions are sufficient
- AWS service limits aren't exceeded
- Region/VPC/subnet configuration is correct

### Network connectivity issues

Verify:
- VPN/Direct Connect is established
- Routing tables are configured
- Security groups allow required traffic
- Network ACLs aren't blocking traffic

### Instances aren't terminated

This could be a bug. Check:
- NRP-AWS-KIP logs for errors
- AWS Console for orphaned instances
- Consider adding a Lambda function for cleanup

## Development Questions

### Can I contribute?

Yes! See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

### What are the high-priority features?

- Container runtime integration
- CloudWatch Logs integration
- Exec support via SSM
- GPU support
- Persistent state storage

### How do I report bugs?

Open an issue on GitHub with:
- Description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Logs and error messages
- Your configuration (redacted)

### Can I use this for production?

This is an early-stage project. It works for basic use cases but has limitations:
- No container runtime integration yet
- No persistent storage
- Basic instance type selection
- No log/exec support

Use for testing and development. Production use requires additional work.

## Miscellaneous

### Why Go?

- Kubernetes ecosystem is Go-based
- Excellent concurrency support
- Great performance
- Strong AWS SDK support

### What about other clouds?

The architecture could be extended to support Azure (Azure Container Instances) or GCP (Cloud Run, GKE Autopilot). Contributions welcome!

### How is this different from cluster-api-provider-aws?

- **Cluster API**: Creates full Kubernetes clusters in AWS
- **NRP-AWS-KIP**: Extends existing cluster with virtual nodes

NRP-AWS-KIP is lighter-weight and designed specifically for bursting scenarios.

### Is there a Helm chart?

Not yet, but it's on the roadmap. For now, use the raw Kubernetes manifests in `deployments/kubernetes/`.

### Where can I get help?

- GitHub Issues
- NRP Slack channel
- Email maintainers
- AWS Support (for AWS-specific issues)
