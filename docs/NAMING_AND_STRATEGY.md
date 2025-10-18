# Project Naming & Strategic Decisions

## What Does "Kip" Stand For?

**Kip = Kubernetes Instance Provider** (or "Kubernetes Cloud Instance Provider")

From Elotl's GitHub README: *"Kip, the Kubernetes Cloud Instance Provider"*

### The "Cell" Metaphor

Kip uses biological terminology:
- **Kip** (the provider) = The organism
- **Cells** = Individual EC2 instances
- **Itzo** = The cell agent (runs inside instances)

This metaphor is clever but perhaps too abstract for research/academic context.

---

## Do We Need a New Name?

### Arguments for Keeping "nrp-aws-kip"

**Pros:**
- ✅ Acknowledges inspiration from Kip
- ✅ "Kubernetes Instance Provider" is descriptive
- ✅ Searchable (Kip + NRP + AWS)
- ✅ Respects Elotl's contribution

**Cons:**
- ❌ May confuse people ("Is this Kip or not?")
- ❌ "Kip" doesn't convey university/research context
- ❌ Lowercase "kip" looks odd in branding
- ❌ Ties identity to a dead project

### Arguments for a New Name

**Pros:**
- ✅ Clear identity (not "Kip fork")
- ✅ Can convey research/university context
- ✅ Fresh start, modern branding
- ✅ Better for publications
- ✅ No confusion with original Kip

**Cons:**
- ❌ Lose "Instance Provider" descriptiveness
- ❌ New acronym to establish
- ❌ May need to explain more

---

## Name Suggestions

### Option 1: Keep Current Name

**nrp-aws-kip** or **NRP-AWS-KIP**

*Tagline: "Kubernetes Instance Provider for NRP and AWS"*

**When to use this:**
- If you want to acknowledge Kip heritage
- If "Kubernetes Instance Provider" resonates
- If you're not concerned about confusion

---

### Option 2: Research/Academic Names

#### **BURSTER**
**B**ursting **U**tility for **R**esearch **S**upercomputing **T**o **E**xternal **R**esources

*Tagline: "Cloud bursting for research computing"*

**Pros:**
- Clear verb ("burst to cloud")
- Memorable, action-oriented
- Good for academic context
- Not an overloaded acronym

---

#### **SURGE**
**S**calable **U**niversity **R**esearch **G**PU **E**xtension

*Tagline: "Surge to AWS for GPU capacity"*

**Pros:**
- Conveys scaling/bursting
- Short, punchy
- Emphasizes GPU (key use case)
- University-focused

---

#### **RAPIDS**
**R**esearch **A**ccess **P**latform for **I**nstance **D**eployment and **S**cheduling

*Tagline: "Rapid access to cloud GPUs for research"*

**Pros:**
- RAPIDS is established in data science world (NVIDIA)
- Conveys speed
- Professional sounding

**Cons:**
- Potential confusion with NVIDIA RAPIDS

---

#### **ASCENT**
**A**cademic **S**caling and **C**loud **E**xtension for **N**RP **T**enants

*Tagline: "Ascend to the cloud for unlimited compute"*

**Pros:**
- Upward motion (scaling)
- Academic focus
- NRP-specific (in acronym)

---

### Option 3: Technical/Descriptive Names

#### **K8s-CloudBurst** or **CloudBurst**

*Tagline: "Kubernetes cloud bursting made simple"*

**Pros:**
- Extremely clear what it does
- Self-documenting
- Good for documentation

**Cons:**
- Generic (could apply to any cloud bursting)
- No acronym pizzazz

---

#### **VK-AWS** or **Virtual-Node-AWS**

*Tagline: "Virtual Kubelet provider for AWS"*

**Pros:**
- Technically accurate
- Clear scope

**Cons:**
- Boring
- No personality
- Doesn't convey research context

---

#### **PodLaunch** or **LaunchPod**

*Tagline: "Launch pods on AWS from anywhere"*

**Pros:**
- Clear what it does
- Catchy
- Easy to remember

**Cons:**
- Doesn't convey bursting concept
- Could apply to any pod launcher

---

### Option 4: NRP-Specific Names

#### **NRP-Burst**

*Tagline: "NRP's gateway to AWS compute"*

**Pros:**
- NRP front-and-center
- Clear purpose (bursting)
- Simple

**Cons:**
- Less interesting for non-NRP users
- Limits scope if you want broader adoption

---

#### **Nautilus-AWS-Bridge** (Nautilus is NRP's K8s platform)

*Tagline: "Bridge Nautilus to AWS"*

**Pros:**
- References NRP's actual platform
- "Bridge" is good metaphor

**Cons:**
- Long name
- Limits to Nautilus (NRP has other clusters)

---

#### **ORCA**
**O**rchestration for **R**esearch **C**loud **A**ccess

*Tagline: "Dive into cloud resources"*

**Pros:**
- Ocean theme matches Nautilus
- Orca = powerful, intelligent
- Research-focused

---

### Option 5: Hybrid Names (Keep Kip Heritage)

#### **Kip-R** (Kip for Research)

*Tagline: "Kip, modernized for research computing"*

**Pros:**
- Acknowledges Kip
- "R" distinguishes it
- Short

**Cons:**
- Still tied to Kip name

---

#### **Kip-NRP** or **Kip-Next**

*Tagline: "Next generation Kubernetes Instance Provider"*

**Pros:**
- Clear evolution
- Keeps KIP acronym utility

**Cons:**
- Perpetuates confusion

---

## Recommendation: **BURSTER** 🎯

### Why BURSTER Works Best

1. **Clear Purpose**: "Burst" is exactly what it does
2. **Action-Oriented**: Verb form (burst-er, thing that bursts)
3. **Memorable**: Short, punchy, one word
4. **Academic Context**: "Research" in acronym
5. **Not Overloaded**: Not a common term in K8s space
6. **Professional**: Sounds serious, not cutesy
7. **Flexible**: Works for NRP and beyond
8. **Good for Papers**: "We present BURSTER, a system for..."

### Branding

**Full Name**: BURSTER - Bursting Utility for Research Supercomputing To External Resources

**Pronunciation**: "Burster" (like "blaster")

**Logo Concept**: ☁️ ⚡ 💥 (cloud with burst rays)

**Tagline**: "Cloud bursting for research computing"

**GitHub**: `github.com/scttfrdmn/burster` or `github.com/aws/burster`

**Website**: `burster.io` (if available)

---

## Strategic Positioning

### As BURSTER (not nrp-aws-kip):

**Messaging**:
```
BURSTER is an open-source Kubernetes provider that enables
research institutions to seamlessly burst workloads to AWS,
with native support for GPU-intensive AI/ML computing.

Built by AWS in partnership with the National Research Platform
and San Diego State University.
```

**Positioning**:
- 🎓 Research-first (not general enterprise)
- 🖥️ GPU-focused (AI/ML workloads)
- 🔬 Open source (community-driven)
- 🌉 Bridge solution (on-prem ↔ cloud)

**Target Audience**:
1. Universities with HPC/research computing
2. Research labs needing occasional GPU access
3. National labs with burst requirements
4. Academic consortia like NRP

**Differentiation**:
- ✨ Not general cloud bursting (research-specific)
- 🎯 Not just another Virtual Kubelet (GPU-optimized)
- 🤝 Not vendor lock-in (open source, Apache 2.0)
- 📚 Not black box (research-grade codebase)

---

## Instance Selection Strategy (Your Point)

### Explicit > Auto-Selection (Agree 100%)

**Your rationale is perfect**:
- GPU users know what they want (P6 vs G6e vs P5)
- Research workloads have specific requirements
- Cost/performance tradeoffs are personal
- "Smart" auto-selection often gets it wrong

### Implementation Approach

**Priority 1: Explicit Selection** ✅
```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    burster.io/instance-type: "p6.48xlarge"  # Explicit
    burster.io/launch-type: "spot"           # Explicit
```

**Priority 2: Template-Based** ✅
```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    burster.io/workload-template: "llm-training"  # Maps to p5.48xlarge
```

**Priority 3: Auto-Selection** (Future)
```yaml
# Only if no annotation specified
# Use simple heuristics:
# - GPU requested → smallest GPU instance that fits
# - No GPU → smallest CPU instance that fits
```

### Configuration
```yaml
# config.yaml
instanceSelection:
  mode: "explicit"  # explicit, template, auto

  # Template definitions
  templates:
    llm-training:
      instanceType: p5.48xlarge
      launchType: spot
      maxPrice: "20.00"

    vision-training:
      instanceType: g5.4xlarge
      launchType: spot

    inference:
      instanceType: g6.2xlarge
      launchType: on-demand

  # Auto-selection (disabled by default)
  auto:
    enabled: false
    strategy: "cheapest-fit"  # or "performance-first"
```

**This gives users**:
- Full control when they want it (explicit)
- Convenience when they want it (templates)
- Fallback if needed (auto, optional)

---

## Examining Kip Issues (Your Other Point)

Let me extract lessons from Kip's GitHub issues:

### Common Issues in Kip

**Category: Pod Lifecycle**
- Pods stuck in pending state
- Pods not cleaning up after termination
- Status not updating correctly
- Orphaned instances after crash

**Lessons**:
- ✅ Robust state machine is critical
- ✅ Idempotency in all operations
- ✅ Graceful handling of partial failures
- ✅ Cleanup jobs for orphaned resources

---

**Category: Networking**
- DNS resolution failures
- Service discovery not working
- Pod-to-pod communication issues
- Security group problems

**Lessons**:
- ✅ Test networking thoroughly
- ✅ Document network requirements clearly
- ✅ Provide troubleshooting guides
- ✅ Make security group config explicit

---

**Category: GPU Support**
- GPU not detected in pods
- Wrong GPU instance selected
- NVIDIA driver issues
- Multi-GPU allocation problems

**Lessons**:
- ✅ Explicit GPU instance selection (your point!)
- ✅ Pre-built AMIs with drivers
- ✅ Clear GPU detection/testing
- ✅ Document GPU requirements

---

**Category: Logs & Debugging**
- No logs from pods
- kubectl logs not working
- Hard to debug pod failures
- No visibility into instance issues

**Lessons**:
- ✅ Log streaming is table-stakes
- ✅ CloudWatch integration early
- ✅ Rich error messages
- ✅ Debug mode for troubleshooting

---

**Category: Cost Control**
- Spot instances terminating unexpectedly
- Instances not terminating (cost overrun)
- No visibility into costs
- Surprise bills

**Lessons**:
- ✅ Budget controls built-in
- ✅ Clear spot instance handling
- ✅ Automated cleanup
- ✅ Cost visibility dashboard

---

**Category: Configuration**
- Complex configuration
- Unclear parameters
- Defaults don't work
- Hard to update config

**Lessons**:
- ✅ Simple config with good defaults
- ✅ Excellent documentation
- ✅ Config validation
- ✅ Hot-reload where possible

---

**Category: Instance Selection**
- Wrong instance type chosen
- Over-provisioning (waste $$$)
- Under-provisioning (OOM kills)
- No control over selection

**Lessons**:
- ✅ Your point! Explicit selection priority
- ✅ Clear templates for common cases
- ✅ Let users override everything
- ✅ Document instance selection clearly

---

## Action Items

### 1. Decide on Name

**My recommendation**: **BURSTER**

**Your decision**: (Want to discuss pros/cons of other options?)

### 2. Update Project

If choosing BURSTER:
```bash
# Rename repo
mv nrp-aws-kip burster

# Update all references
- README.md
- Documentation
- Code (package names)
- GitHub repo

# New branding
- Logo
- Tagline
- Positioning
```

### 3. Document Kip Lessons

Create `docs/KIP_LESSONS_LEARNED.md`:
- Study Kip's GitHub issues (I can help extract)
- Document common problems
- Design decisions to avoid them
- Testing checklist

### 4. Instance Selection Design

Create `docs/INSTANCE_SELECTION.md`:
- Explicit selection (priority 1)
- Template-based (priority 2)
- Auto-selection (future, optional)
- User documentation

---

## Quick Poll: Name Preferences

Rank these (1 = favorite):

- [ ] **BURSTER** (my recommendation)
- [ ] **SURGE**
- [ ] **ORCA**
- [ ] **ASCENT**
- [ ] Keep **nrp-aws-kip**
- [ ] Other idea: _____________

---

## Bottom Line

**Recommendation**:
1. ✅ **Name it BURSTER** - Clear, memorable, research-focused
2. ✅ **Explicit instance selection first** - Your instinct is correct
3. ✅ **Study Kip issues** - Learn from their pain points
4. ✅ **Modern architecture** - 2025 best practices

**This positions you perfectly**:
- Clear identity (not "Kip but newer")
- Research focus (universities are your market)
- User-friendly (explicit control over instances)
- Battle-tested design (learn from Kip's mistakes)

Want me to:
1. Create a detailed Kip issues analysis?
2. Design the explicit instance selection system?
3. Create the BURSTER branding guide?
4. Something else?
