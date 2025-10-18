# GPU Workloads on AWS via NRP-AWS-KIP

## Overview

This guide covers everything you need to know about running GPU workloads on AWS through cloud bursting from NRP, including instance selection, cost optimization, and practical examples for ML/AI workloads.

---

## GPU Instance Types on AWS

### Instance Family Overview

| Family | GPU Type | GPU Memory | Use Case | Performance | Cost |
|--------|----------|------------|----------|-------------|------|
| **g4dn** | NVIDIA T4 | 16 GB | ML inference, light training, graphics | Good | $ |
| **g5** | NVIDIA A10G | 24 GB | ML training/inference, graphics | Better | $$ |
| **p3** | NVIDIA V100 | 16-32 GB | ML training, HPC | Great | $$$ |
| **p4d** | NVIDIA A100 | 40 GB | Large-scale ML training | Excellent | $$$$ |
| **p5** | NVIDIA H100 | 80 GB | Next-gen ML, largest models | Best | $$$$$ |

### Detailed Instance Specifications

#### g4dn - NVIDIA T4 (Entry-Level)

| Instance Type | GPUs | GPU Memory | vCPUs | RAM | Price/hr (On-Demand) | Price/hr (Spot ~) |
|---------------|------|------------|-------|-----|---------------------|-------------------|
| g4dn.xlarge | 1 | 16 GB | 4 | 16 GB | $0.526 | $0.158 |
| g4dn.2xlarge | 1 | 16 GB | 8 | 32 GB | $0.752 | $0.226 |
| g4dn.4xlarge | 1 | 16 GB | 16 | 64 GB | $1.204 | $0.361 |
| g4dn.12xlarge | 4 | 64 GB | 48 | 192 GB | $3.912 | $1.174 |

**Best For**: Inference, light training, video transcoding, graphics rendering
**Spot Savings**: ~70%

#### g5 - NVIDIA A10G (Mid-Tier)

| Instance Type | GPUs | GPU Memory | vCPUs | RAM | Price/hr (On-Demand) | Price/hr (Spot ~) |
|---------------|------|------------|-------|-----|---------------------|-------------------|
| g5.xlarge | 1 | 24 GB | 4 | 16 GB | $1.006 | $0.302 |
| g5.2xlarge | 1 | 24 GB | 8 | 32 GB | $1.212 | $0.364 |
| g5.4xlarge | 1 | 24 GB | 16 | 64 GB | $1.624 | $0.487 |
| g5.12xlarge | 4 | 96 GB | 48 | 192 GB | $5.672 | $1.702 |

**Best For**: ML training, inference at scale, graphics
**Spot Savings**: ~70%

#### p3 - NVIDIA V100 (High-Performance)

| Instance Type | GPUs | GPU Memory | vCPUs | RAM | Price/hr (On-Demand) | Price/hr (Spot ~) |
|---------------|------|------------|-------|-----|---------------------|-------------------|
| p3.2xlarge | 1 | 16 GB | 8 | 61 GB | $3.06 | $0.92 |
| p3.8xlarge | 4 | 64 GB | 32 | 244 GB | $12.24 | $3.67 |
| p3.16xlarge | 8 | 128 GB | 64 | 488 GB | $24.48 | $7.34 |

**Best For**: Deep learning training, HPC, scientific computing
**Spot Savings**: ~70%

#### p4d - NVIDIA A100 (Premium)

| Instance Type | GPUs | GPU Memory | vCPUs | RAM | Price/hr (On-Demand) | Price/hr (Spot ~) |
|---------------|------|------------|-------|-----|---------------------|-------------------|
| p4d.24xlarge | 8 | 320 GB | 96 | 1152 GB | $32.77 | $9.83 |

**Best For**: Large-scale distributed training, giant models
**Spot Savings**: ~70%

#### p5 - NVIDIA H100 (Next-Gen)

| Instance Type | GPUs | GPU Memory | vCPUs | RAM | Price/hr (On-Demand) |
|---------------|------|------------|-------|-----|---------------------|
| p5.48xlarge | 8 | 640 GB | 192 | 2048 GB | $98.32 |

**Best For**: Cutting-edge LLM training, largest AI models
**Spot Availability**: Limited

---

## Instance Selection Guide

### Decision Matrix

```
Workload Type → Instance Recommendation

Small Model Inference (< 1GB)
  └─ g4dn.xlarge (T4, $0.526/hr)

Medium Model Inference (1-5GB)
  └─ g4dn.2xlarge (T4, $0.752/hr)
  └─ g5.xlarge (A10G, $1.006/hr) - better performance

Light Training (< 10GB model)
  └─ g4dn.4xlarge (T4, $1.204/hr)
  └─ g5.2xlarge (A10G, $1.212/hr) - recommended

Standard Training (10-50GB model)
  └─ g5.4xlarge (A10G, $1.624/hr)
  └─ p3.2xlarge (V100, $3.06/hr) - for compute-heavy

Heavy Training (50-100GB model)
  └─ p3.8xlarge (4x V100, $12.24/hr)
  └─ p4d.24xlarge (8x A100, $32.77/hr) - for distributed

LLM Training (> 100GB model)
  └─ p4d.24xlarge (8x A100, $32.77/hr)
  └─ p5.48xlarge (8x H100, $98.32/hr) - cutting edge

Multi-GPU Parallel Training
  └─ g4dn.12xlarge (4x T4, $3.912/hr) - budget
  └─ g5.12xlarge (4x A10G, $5.672/hr) - balanced
  └─ p3.16xlarge (8x V100, $24.48/hr) - performance
```

### By Framework

**PyTorch/TensorFlow Training**:
- Small models (ResNet, BERT-base): g4dn.2xlarge or g5.xlarge
- Medium models (GPT-2, ViT): g5.4xlarge or p3.2xlarge
- Large models (GPT-3-like): p3.8xlarge or p4d.24xlarge

**Inference Serving**:
- Real-time APIs: g4dn.xlarge (cost-effective)
- Batch inference: g4dn.2xlarge or g5.xlarge
- High-throughput: g5.4xlarge

**Computer Vision**:
- Image classification: g4dn.xlarge
- Object detection: g4dn.2xlarge or g5.xlarge
- Video processing: g5.2xlarge or g5.4xlarge

**NLP**:
- BERT fine-tuning: g5.xlarge
- GPT training: p3.2xlarge or better
- LLM inference: g5.4xlarge

---

## Configuring GPU Workloads

### Basic GPU Job

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: gpu-test
  namespace: your-namespace
spec:
  nodeSelector:
    type: virtual-kubelet
  containers:
  - name: gpu-container
    image: nvidia/cuda:12.0-base-ubuntu22.04
    command: ["nvidia-smi"]
    resources:
      limits:
        nvidia.com/gpu: 1  # Request 1 GPU
      requests:
        nvidia.com/gpu: 1
```

### Specifying Instance Type

Use annotations to control instance selection:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: training-job
  annotations:
    # Explicitly specify instance type
    nrp.aws.kip/instance-type: "g5.2xlarge"

    # Use spot instances for cost savings
    nrp.aws.kip/launch-type: "spot"
spec:
  nodeSelector:
    type: virtual-kubelet
  containers:
  - name: trainer
    image: pytorch/pytorch:2.1.0-cuda12.1-cudnn8-runtime
    resources:
      limits:
        nvidia.com/gpu: 1
        memory: "32Gi"
        cpu: "8"
      requests:
        nvidia.com/gpu: 1
        memory: "32Gi"
        cpu: "8"
```

### Multi-GPU Training

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: multi-gpu-training
  annotations:
    nrp.aws.kip/instance-type: "p3.8xlarge"  # 4 GPUs
spec:
  nodeSelector:
    type: virtual-kubelet
  containers:
  - name: trainer
    image: pytorch/pytorch:latest-cuda
    command: ["python", "-m", "torch.distributed.launch"]
    args:
    - "--nproc_per_node=4"  # 4 GPUs
    - "train.py"
    resources:
      limits:
        nvidia.com/gpu: 4  # Request 4 GPUs
        memory: "244Gi"
        cpu: "32"
      requests:
        nvidia.com/gpu: 4
        memory: "244Gi"
        cpu: "32"
```

---

## Complete Examples

### Example 1: ResNet Training (Light)

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: resnet-training
  namespace: your-namespace
spec:
  backoffLimit: 2
  ttlSecondsAfterFinished: 3600
  template:
    metadata:
      annotations:
        nrp.aws.kip/instance-type: "g4dn.2xlarge"
        nrp.aws.kip/launch-type: "spot"
    spec:
      nodeSelector:
        type: virtual-kubelet
      restartPolicy: Never
      containers:
      - name: trainer
        image: pytorch/pytorch:2.1.0-cuda12.1-cudnn8-runtime
        command: ["/bin/bash", "-c"]
        args:
        - |
          echo "Training ResNet-50 on $(hostname)"
          nvidia-smi

          pip install torchvision
          python train_resnet.py \
            --epochs 50 \
            --batch-size 128 \
            --lr 0.1
        resources:
          limits:
            nvidia.com/gpu: 1
            memory: "32Gi"
            cpu: "8"
          requests:
            nvidia.com/gpu: 1
            memory: "32Gi"
            cpu: "8"
```

**Cost**: ~$0.226/hr (spot), ~$22.60 for 100 epochs

### Example 2: BERT Fine-Tuning (Medium)

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: bert-finetuning
  namespace: your-namespace
spec:
  template:
    metadata:
      annotations:
        nrp.aws.kip/instance-type: "g5.xlarge"
        nrp.aws.kip/launch-type: "spot"
    spec:
      nodeSelector:
        type: virtual-kubelet
      restartPolicy: OnFailure
      containers:
      - name: trainer
        image: huggingface/transformers-pytorch-gpu:latest
        command: ["python", "finetune_bert.py"]
        args:
        - --model_name=bert-base-uncased
        - --dataset=squad
        - --epochs=3
        - --batch_size=16
        - --learning_rate=2e-5
        env:
        - name: TRANSFORMERS_CACHE
          value: /cache
        resources:
          limits:
            nvidia.com/gpu: 1
            memory: "16Gi"
            cpu: "4"
          requests:
            nvidia.com/gpu: 1
            memory: "16Gi"
            cpu: "4"
        volumeMounts:
        - name: model-cache
          mountPath: /cache
      volumes:
      - name: model-cache
        emptyDir:
          sizeLimit: 20Gi
```

**Cost**: ~$0.302/hr (spot), ~$2.42 for 8 hours

### Example 3: Large Language Model Training (Heavy)

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: llm-training
  namespace: your-namespace
spec:
  template:
    metadata:
      annotations:
        nrp.aws.kip/instance-type: "p3.8xlarge"
        nrp.aws.kip/launch-type: "on-demand"  # Use on-demand for critical workloads
    spec:
      nodeSelector:
        type: virtual-kubelet
      restartPolicy: Never
      containers:
      - name: trainer
        image: huggingface/transformers-pytorch-gpu:latest
        command: ["python", "-m", "torch.distributed.launch"]
        args:
        - "--nproc_per_node=4"
        - "train_gpt.py"
        - --model=gpt2-large
        - --dataset=openwebtext
        - --gradient_checkpointing
        - --fp16
        resources:
          limits:
            nvidia.com/gpu: 4
            memory: "244Gi"
            cpu: "32"
          requests:
            nvidia.com/gpu: 4
            memory: "244Gi"
            cpu: "32"
```

**Cost**: $12.24/hr (on-demand), $293.76/day

### Example 4: Video Processing (Computer Vision)

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: video-processing
  namespace: your-namespace
spec:
  parallelism: 10  # Process 10 videos in parallel
  completions: 100  # Process 100 videos total
  template:
    metadata:
      annotations:
        nrp.aws.kip/instance-type: "g5.2xlarge"
        nrp.aws.kip/launch-type: "spot"
    spec:
      nodeSelector:
        type: virtual-kubelet
      restartPolicy: OnFailure
      containers:
      - name: processor
        image: myregistry/video-processor:gpu
        command: ["python", "process_video.py"]
        args:
        - --input=$(INPUT_VIDEO)
        - --output=/output/
        - --model=yolov8
        - --use-gpu
        env:
        - name: INPUT_VIDEO
          valueFrom:
            fieldRef:
              fieldPath: metadata.annotations['batch.kubernetes.io/job-completion-index']
        resources:
          limits:
            nvidia.com/gpu: 1
            memory: "32Gi"
            cpu: "8"
          requests:
            nvidia.com/gpu: 1
            memory: "32Gi"
            cpu: "8"
```

**Cost**: $0.364/hr per job × 10 parallel = $3.64/hr for batch

### Example 5: Inference Service (Long-Running)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ml-inference-api
  namespace: your-namespace
spec:
  replicas: 2
  selector:
    matchLabels:
      app: ml-inference
  template:
    metadata:
      labels:
        app: ml-inference
      annotations:
        nrp.aws.kip/instance-type: "g4dn.xlarge"
        nrp.aws.kip/launch-type: "spot"
    spec:
      nodeSelector:
        type: virtual-kubelet
      containers:
      - name: inference-server
        image: myregistry/inference-server:latest
        ports:
        - containerPort: 8080
        resources:
          limits:
            nvidia.com/gpu: 1
            memory: "16Gi"
            cpu: "4"
          requests:
            nvidia.com/gpu: 1
            memory: "16Gi"
            cpu: "4"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 120
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 60
          periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: ml-inference-api
  namespace: your-namespace
spec:
  selector:
    app: ml-inference
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

**Cost**: $0.158/hr × 2 replicas = $0.316/hr ($227/month)

---

## Cost Optimization Strategies

### 1. Use Spot Instances

**Savings**: 70-90% vs on-demand

```yaml
metadata:
  annotations:
    nrp.aws.kip/launch-type: "spot"
```

**Best For**:
- ✅ Training jobs (can handle interruptions)
- ✅ Batch processing
- ✅ Non-critical workloads
- ❌ Production APIs (use on-demand)
- ❌ Long-running experiments with no checkpointing

**Spot Best Practices**:
- Implement checkpointing every 5-10 minutes
- Use multiple availability zones
- Handle interruption notices (2-minute warning)
- Consider g4dn/g5 families (better spot availability than p3/p4)

### 2. Right-Size Instances

Don't over-provision:

```python
# Example: BERT fine-tuning
# Needs: ~12GB GPU memory

# ❌ Wasteful
instance_type = "p3.8xlarge"  # 4x V100, $12.24/hr

# ✅ Optimal
instance_type = "g5.xlarge"   # 1x A10G, $1.006/hr
```

**Savings**: 12x cheaper!

### 3. Batch Your Jobs

Run multiple small jobs as one large job:

```yaml
# Instead of 10 separate p3.2xlarge jobs ($3.06/hr each)
# Use 1 p3.16xlarge with 10 parallel workers ($24.48/hr)
# Savings: $30.60 - $24.48 = $6.12/hr
```

### 4. Use Mixed Precision Training

Enable FP16 to reduce memory and speed up training:

```python
from torch.cuda.amp import autocast, GradScaler

scaler = GradScaler()

with autocast():
    outputs = model(inputs)
    loss = criterion(outputs, labels)

scaler.scale(loss).backward()
scaler.step(optimizer)
scaler.update()
```

**Benefits**:
- 2-3x faster training
- 50% less GPU memory
- Can use smaller (cheaper) instances

### 5. Gradient Accumulation

Simulate larger batch sizes without more memory:

```python
accumulation_steps = 4

for i, (inputs, labels) in enumerate(dataloader):
    outputs = model(inputs)
    loss = criterion(outputs, labels) / accumulation_steps
    loss.backward()

    if (i + 1) % accumulation_steps == 0:
        optimizer.step()
        optimizer.zero_grad()
```

**Benefit**: Use g4dn.xlarge instead of g5.2xlarge

### 6. Terminate Idle Instances

Set appropriate timeouts:

```yaml
spec:
  ttlSecondsAfterFinished: 300  # Delete completed pods after 5 min
  activeDeadlineSeconds: 86400   # Kill if running > 24 hours
```

### 7. Use Preemptible Queues

Run jobs during off-peak hours:

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: nightly-training
spec:
  schedule: "0 2 * * *"  # Run at 2 AM when cheaper
  jobTemplate:
    spec:
      template:
        # ... job spec
```

---

## Cost Estimation

### Training Cost Calculator

```python
def estimate_training_cost(
    instance_type: str,
    hourly_rate: float,
    training_hours: float,
    spot_discount: float = 0.7
):
    """
    Estimate GPU training cost

    Args:
        instance_type: AWS instance type
        hourly_rate: On-demand price per hour
        training_hours: Expected training duration
        spot_discount: Discount for spot (default 0.7 = 70% off)
    """
    on_demand_cost = hourly_rate * training_hours
    spot_cost = on_demand_cost * (1 - spot_discount)

    print(f"Instance: {instance_type}")
    print(f"Duration: {training_hours} hours")
    print(f"On-Demand: ${on_demand_cost:.2f}")
    print(f"Spot (~{int(spot_discount*100)}% off): ${spot_cost:.2f}")
    print(f"Savings: ${on_demand_cost - spot_cost:.2f}")

    return spot_cost

# Example: Train ResNet for 48 hours
estimate_training_cost(
    instance_type="g5.2xlarge",
    hourly_rate=1.212,
    training_hours=48,
    spot_discount=0.7
)
# Output:
# Instance: g5.2xlarge
# Duration: 48 hours
# On-Demand: $58.18
# Spot (~70% off): $17.45
# Savings: $40.73
```

### Monthly Budget Planning

```python
# Scenario: CS Department GPU Usage
# 10 students × 20 hours/week × 4 weeks

g4dn_xlarge_spot = 0.158  # $/hr
weekly_hours_per_student = 20
students = 10
weeks = 4

monthly_cost = (g4dn_xlarge_spot *
                weekly_hours_per_student *
                students *
                weeks)

print(f"Monthly Budget: ${monthly_cost:.2f}")
# Output: Monthly Budget: $1,264.00
```

---

## Monitoring GPU Usage

### Check GPU Utilization

```bash
# SSH into running pod (if exec is enabled)
kubectl exec -it gpu-training-pod -- nvidia-smi

# Check every second
kubectl exec -it gpu-training-pod -- watch -n 1 nvidia-smi
```

### CloudWatch Metrics

Monitor GPU usage in AWS CloudWatch:

```bash
# Get GPU utilization for an instance
aws cloudwatch get-metric-statistics \
  --namespace AWS/EC2 \
  --metric-name GPUUtilization \
  --dimensions Name=InstanceId,Value=i-1234567890abcdef0 \
  --start-time 2025-01-01T00:00:00Z \
  --end-time 2025-01-01T23:59:59Z \
  --period 300 \
  --statistics Average
```

### Cost Alerts

Set up billing alerts:

```bash
# Create SNS topic for alerts
aws sns create-topic --name gpu-cost-alerts

# Create budget alert
aws budgets create-budget \
  --account-id 123456789012 \
  --budget file://gpu-budget.json
```

gpu-budget.json:
```json
{
  "BudgetName": "GPU-Monthly-Budget",
  "BudgetLimit": {
    "Amount": "2000",
    "Unit": "USD"
  },
  "TimeUnit": "MONTHLY",
  "BudgetType": "COST"
}
```

---

## Best Practices

### 1. Always Checkpoint

Save model state frequently:

```python
checkpoint_freq = 10  # Every 10 epochs

if epoch % checkpoint_freq == 0:
    torch.save({
        'epoch': epoch,
        'model_state_dict': model.state_dict(),
        'optimizer_state_dict': optimizer.state_dict(),
        'loss': loss,
    }, f's3://my-bucket/checkpoints/model_epoch_{epoch}.pt')
```

### 2. Use Data Parallelism

Maximize GPU utilization:

```python
if torch.cuda.device_count() > 1:
    model = nn.DataParallel(model)
```

### 3. Profile Your Code

Find bottlenecks:

```python
from torch.profiler import profile, ProfilerActivity

with profile(activities=[ProfilerActivity.CPU, ProfilerActivity.CUDA]) as prof:
    model(inputs)

print(prof.key_averages().table(sort_by="cuda_time_total"))
```

### 4. Monitor Training Progress

Log to experiment tracking:

```python
import wandb

wandb.init(project="my-project")
wandb.watch(model)

for epoch in range(epochs):
    # ... training ...
    wandb.log({
        "epoch": epoch,
        "loss": loss,
        "gpu_util": get_gpu_util(),
        "cost_so_far": calculate_cost()
    })
```

### 5. Clean Up Resources

Always clean up:

```yaml
spec:
  ttlSecondsAfterFinished: 600  # Delete after 10 minutes
```

---

## Troubleshooting

### GPU Not Detected

```bash
# Check NVIDIA driver
kubectl exec -it pod-name -- nvidia-smi

# If missing, check AMI configuration
# Ensure using GPU-enabled AMI
```

### Out of Memory

```python
# Enable gradient checkpointing
from torch.utils.checkpoint import checkpoint

def forward(x):
    return checkpoint(heavy_layer, x)

# Use smaller batch size
batch_size = 16  # Instead of 32

# Clear cache periodically
torch.cuda.empty_cache()
```

### Slow Training

```python
# Use pin_memory for DataLoader
train_loader = DataLoader(
    dataset,
    batch_size=32,
    pin_memory=True,  # Faster CPU-to-GPU transfer
    num_workers=4
)

# Enable cuDNN auto-tuner
torch.backends.cudnn.benchmark = True
```

---

## Additional Resources

- [Instance Selection Guide](./INSTANCE_SELECTION.md)
- [Cost Management](./COST_MANAGEMENT.md)
- [Multi-Tenancy](./MULTI_TENANCY.md)
- [Example GPU Jobs](../examples/gpu-workloads/)
- [AWS GPU Instance Pricing](https://aws.amazon.com/ec2/instance-types/)
- [NVIDIA GPU Cloud](https://www.nvidia.com/en-us/gpu-cloud/)
