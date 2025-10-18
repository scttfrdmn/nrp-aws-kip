# Getting Started with NRP-AWS-KIP

This guide will help you set up and deploy NRP-AWS-KIP to enable cloud bursting from your NRP Kubernetes cluster to AWS.

## Prerequisites

### Required

- Access to an NRP Kubernetes cluster with admin permissions
- AWS account with appropriate permissions
- Go 1.21+ (for development)
- kubectl configured for your cluster
- Terraform 1.0+ (for infrastructure setup)

### Optional but Recommended

- VPN or Direct Connect between NRP and AWS
- Docker (for building container images)
- golangci-lint (for development)

## Step 1: AWS Infrastructure Setup

### Deploy with Terraform

1. Navigate to the Terraform directory:

```bash
cd deployments/terraform
```

2. Copy and customize the variables:

```bash
cp terraform.tfvars.example terraform.tfvars
```

3. Edit `terraform.tfvars`:

```hcl
aws_region         = "us-west-2"
project_name       = "nrp-aws-kip"
create_vpc         = true  # or false if using existing VPC
vpc_cidr           = "10.100.0.0/16"
availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
nrp_cidr_blocks    = ["10.0.0.0/8"]  # Update with your NRP CIDR

common_tags = {
  Project     = "NRP-AWS-KIP"
  Team        = "SDSU"
  Environment = "production"
}
```

4. Initialize and apply Terraform:

```bash
terraform init
terraform plan
terraform apply
```

5. Note the outputs - you'll need these for configuration:

```bash
terraform output
```

## Step 2: Configure NRP-AWS-KIP

### Create Configuration File

1. Copy the example config:

```bash
cp config.yaml.example config.yaml
```

2. Update `config.yaml` with values from Terraform outputs:

```yaml
aws:
  region: us-west-2
  vpcId: vpc-xxxxx           # From Terraform output
  subnetIds:
    - subnet-xxxxx           # From Terraform output
    - subnet-yyyyy
  securityGroupIds:
    - sg-xxxxx               # From Terraform output
  keyName: nrp-burst-key     # Optional: for SSH access
  iamInstanceProfile: nrp-burst-instance-profile  # From Terraform
  tags:
    Project: NRP
    ManagedBy: nrp-aws-kip
    Team: SDSU
  useSpotInstances: true
  spotMaxPrice: "0.50"

node:
  name: nrp-aws-burst
  operatingSystem: Linux
  cpu: "1000"                # Adjust based on your needs
  memory: "4000Gi"
  pods: "500"
  labels:
    type: burst
    cloud: aws
    nrp.org/burst-node: "true"
  taints:
    - key: nrp.org/burst
      value: "true"
      effect: NoSchedule

limits:
  maxPods: 500
  maxInstances: 100
  costLimit: "1000.00"       # Optional: USD per month

logging:
  level: info
  format: json

metrics:
  enabled: true
  port: 10255
  path: /metrics
```

## Step 3: Build and Deploy

### Option A: Deploy to Kubernetes (Recommended)

1. Update the Kubernetes manifests with your configuration:

```bash
cd deployments/kubernetes
```

2. Edit `deployment.yaml`:
   - Update the ConfigMap with your `config.yaml` contents
   - Add your AWS credentials to the Secret

3. Build the Docker image:

```bash
make docker-build
```

4. Push to your container registry (if needed):

```bash
docker tag nrp-aws-kip:latest your-registry/nrp-aws-kip:latest
docker push your-registry/nrp-aws-kip:latest
```

5. Deploy to Kubernetes:

```bash
kubectl apply -f deployment.yaml
```

6. Verify the deployment:

```bash
kubectl get pods -n nrp-aws-kip
kubectl logs -n nrp-aws-kip -l app=nrp-aws-kip
```

7. Check the virtual node:

```bash
kubectl get nodes
kubectl describe node nrp-aws-burst
```

### Option B: Run Locally (Development)

1. Set AWS credentials:

```bash
export AWS_REGION=us-west-2
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
```

2. Build the binary:

```bash
make build
```

3. Run locally:

```bash
./bin/nrp-aws-kip --config config.yaml --kubeconfig ~/.kube/config
```

## Step 4: Test the Burst Functionality

### Create a Test Pod

1. Create a test pod manifest (`test-burst-pod.yaml`):

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-burst-pod
  namespace: default
spec:
  nodeSelector:
    nrp.org/burst-node: "true"
  tolerations:
    - key: nrp.org/burst
      operator: Equal
      value: "true"
      effect: NoSchedule
  containers:
    - name: test
      image: nginx:latest
      resources:
        requests:
          cpu: 1
          memory: 2Gi
        limits:
          cpu: 2
          memory: 4Gi
```

2. Deploy the test pod:

```bash
kubectl apply -f test-burst-pod.yaml
```

3. Monitor the pod:

```bash
kubectl get pods -w
kubectl describe pod test-burst-pod
```

4. Check the EC2 instance in AWS:

```bash
aws ec2 describe-instances \
  --filters "Name=tag:nrp-aws-kip/pod,Values=test-burst-pod" \
  --region us-west-2
```

5. Clean up:

```bash
kubectl delete pod test-burst-pod
```

## Step 5: Configure Network Connectivity

### Option A: VPN Connection

1. Create a customer gateway in AWS:

```bash
aws ec2 create-customer-gateway \
  --type ipsec.1 \
  --public-ip <NRP-VPN-IP> \
  --bgp-asn 65000
```

2. Create VPN connection:

```bash
aws ec2 create-vpn-connection \
  --type ipsec.1 \
  --customer-gateway-id <cgw-id> \
  --vpn-gateway-id <vgw-id>
```

3. Configure NRP side (work with NRP network team)

### Option B: Direct Connect (Recommended for Production)

Work with your NRP and AWS teams to establish a Direct Connect connection.

## Step 6: Monitoring

### View Metrics

```bash
# Port-forward to access metrics
kubectl port-forward -n nrp-aws-kip svc/nrp-aws-kip-metrics 10255:10255

# View metrics
curl http://localhost:10255/metrics
```

### View Logs

```bash
# Follow logs
kubectl logs -n nrp-aws-kip -l app=nrp-aws-kip -f

# View recent logs
kubectl logs -n nrp-aws-kip -l app=nrp-aws-kip --tail=100
```

### CloudWatch Logs (AWS)

```bash
# View instance logs
aws logs tail /nrp-aws-kip/instances --follow
```

## Step 7: Production Considerations

### Security

1. **Use IAM Roles**: Instead of access keys, use IRSA (IAM Roles for Service Accounts)
2. **Rotate Credentials**: If using access keys, rotate them regularly
3. **Network Policies**: Implement Kubernetes network policies
4. **Audit Logging**: Enable CloudTrail for AWS API calls

### Cost Management

1. **Set Budgets**: Create AWS budgets with alerts
2. **Review Regularly**: Monitor EC2 costs in AWS Cost Explorer
3. **Optimize Instance Types**: Adjust instance selection algorithm
4. **Use Spot Instances**: Enable spot instances for non-critical workloads

### High Availability

1. **Multiple Replicas**: Run multiple replicas of nrp-aws-kip
2. **Leader Election**: Implement leader election (future enhancement)
3. **Multi-Region**: Consider multi-region deployment

### Monitoring and Alerts

1. **Prometheus**: Scrape metrics for monitoring
2. **Alertmanager**: Set up alerts for failures
3. **Grafana**: Create dashboards for visualization

## Troubleshooting

### Pod Stuck in Pending

```bash
# Check node status
kubectl describe node nrp-aws-burst

# Check pod events
kubectl describe pod <pod-name>

# Check nrp-aws-kip logs
kubectl logs -n nrp-aws-kip -l app=nrp-aws-kip
```

### No EC2 Instance Created

1. Check AWS credentials
2. Verify VPC/subnet configuration
3. Check security group rules
4. Review CloudTrail for API errors

### Cannot Connect to Pod

1. Verify network connectivity (VPN/Direct Connect)
2. Check security group rules
3. Verify routing between NRP and AWS

## Next Steps

- Read [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed architecture
- Review [CONTRIBUTING.md](../CONTRIBUTING.md) for development guidelines
- Check out example workloads in `/examples` (coming soon)
- Join the NRP Slack for support

## Getting Help

- GitHub Issues: https://github.com/scttfrdmn/nrp-aws-kip/issues
- NRP Documentation: https://docs.nationalresearchplatform.org
- AWS Support: https://aws.amazon.com/support
