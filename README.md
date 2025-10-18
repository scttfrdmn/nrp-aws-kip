# nrp-aws-kip (Archived)

**⚠️ This project has been superseded by [ORCA](../orca/).**

## What Happened?

This directory contains initial research and planning for a Kubernetes cloud bursting solution. After extensive analysis (see `docs/` directory), the project was renamed and rebuilt from scratch as **ORCA** (Orchestration for Research Cloud Access).

## Why the Change?

1. **Better name** - ORCA reflects the ocean theme of NRP's Nautilus platform
2. **Fresh start** - Building from scratch with modern tools (K8s 1.31, AWS SDK v2, Go 1.21)
3. **Clear identity** - Not a fork of Kip, but a new project learning from its lessons
4. **Research focus** - Designed specifically for university/research computing workloads

## What's Here?

This directory preserves the initial thinking and research:

- `docs/NAMING_AND_STRATEGY.md` - Name selection process (chose ORCA)
- `docs/WHY_BUILD_CUSTOM.md` - Why build from scratch vs fork Kip
- `docs/ALTERNATIVES_ANALYSIS.md` - Analysis of existing solutions
- `docs/KIP_COMPARISON.md` - Detailed Kip comparison
- `docs/GPU_WORKLOADS.md` - GPU instance research (October 2025 pricing)
- `claude-chat.txt` - Original research conversation

## Active Development

**All active development happens in**: [`../orca/`](../orca/)

See the [ORCA README](../orca/README.md) for:
- Current architecture
- Development setup
- Contributing guidelines
- Roadmap

## Historical Value

This directory is kept for historical reference:
- Documents decision-making process
- Preserves research on alternatives
- Shows evolution from concept to implementation
- Useful for understanding "why ORCA"

---

**Status**: 🔒 Archived (reference only)
**Active Project**: ✅ [ORCA](../orca/)
**Date Archived**: October 2025

## Original Concept

This was originally conceived as "nrp-aws-kip" (National Research Platform AWS Kubernetes Instance Provider), a Virtual Kubelet-based solution for cloud bursting. The core concept remains the same in ORCA, but with:
- Better naming and branding
- Modern implementation (2025 tools)
- Research-first design philosophy
- Explicit instance selection
- Comprehensive testing framework

The research and planning done here directly informed ORCA's architecture.
