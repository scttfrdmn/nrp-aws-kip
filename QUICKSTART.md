# Quick Start Guide

Get NRP-AWS-KIP running in 15 minutes!

## Prerequisites

- kubectl and access to NRP Kubernetes cluster
- AWS account with admin permissions
- AWS CLI configured

## 1. Deploy AWS Infrastructure (5 min)

```bash
cd deployments/terraform
cp terraform.tfvars.example terraform.tfvars

# Edit terraform.tfvars with your values
vim terraform.tfvars

# Deploy
terraform init
terraform apply -auto-approve

# Save outputs
terraform output > ../../aws-config.txt
```

## 2. Configure NRP-AWS-KIP (2 min)

```bash
cd ../..
cp config.yaml.example config.yaml

# Update config.yaml with Terraform outputs
# Edit the aws section with vpc, subnets, security groups from aws-config.txt
vim config.yaml
```

## 3. Deploy to Kubernetes (5 min)

```bash
# Update deployment with your config
kubectl create namespace nrp-aws-kip

# Create secret with AWS credentials
kubectl create secret generic nrp-aws-kip-aws-credentials \
  -n nrp-aws-kip \
  --from-literal=AWS_ACCESS_KEY_ID=your_key \
  --from-literal=AWS_SECRET_ACCESS_KEY=your_secret

# Update ConfigMap in deployment.yaml with your config.yaml contents
vim deployments/kubernetes/deployment.yaml

# Deploy
kubectl apply -f deployments/kubernetes/deployment.yaml

# Verify
kubectl get pods -n nrp-aws-kip
kubectl get nodes  # Should see nrp-aws-burst node
```

## 4. Test It (3 min)

```bash
# Create test pod
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: test-burst
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
EOF

# Watch it create
kubectl get pods -w

# Check EC2 instance
aws ec2 describe-instances \
  --filters "Name=tag:nrp-aws-kip/pod,Values=test-burst"

# Clean up
kubectl delete pod test-burst
```

## Next Steps

- Read [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) to understand how it works
- See [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md) for detailed setup
- Check [docs/FAQ.md](docs/FAQ.md) for common questions

## Troubleshooting

**Pods stuck in Pending?**
```bash
kubectl describe pod test-burst
kubectl logs -n nrp-aws-kip -l app=nrp-aws-kip
```

**Virtual node not showing?**
```bash
kubectl get nodes
kubectl logs -n nrp-aws-kip -l app=nrp-aws-kip
```

**AWS errors?**
```bash
# Check credentials
aws sts get-caller-identity

# Check permissions
aws ec2 describe-instances --max-items 1
```

## Support

- GitHub Issues: https://github.com/scttfrdmn/nrp-aws-kip/issues
- Documentation: https://github.com/scttfrdmn/nrp-aws-kip/tree/main/docs
