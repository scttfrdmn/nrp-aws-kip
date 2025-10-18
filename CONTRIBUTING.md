# Contributing to NRP-AWS-KIP

Thank you for your interest in contributing to NRP-AWS-KIP! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and professional in all interactions. We're building this for the research community.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/nrp-aws-kip.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit and push
7. Open a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Docker (for building images)
- kubectl and access to a test Kubernetes cluster
- AWS account for testing
- golangci-lint for linting

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Running Locally

```bash
# Set AWS credentials
export AWS_REGION=us-west-2
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret

# Run
./bin/nrp-aws-kip --config config.yaml --kubeconfig ~/.kube/config
```

## Code Style

### Go Code

- Follow standard Go conventions
- Use `gofmt` for formatting (run `make fmt`)
- Run `go vet` (run `make vet`)
- Use `golangci-lint` for linting (run `make lint`)
- Write clear, descriptive variable and function names
- Add comments for exported functions and types

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(provider): add support for GPU instances

Adds logic to detect GPU resource requests and select appropriate
EC2 instance types (p3, p4, etc.).

Closes #123
```

```
fix(aws): prevent instance leak on pod deletion failure

Ensures EC2 instances are terminated even if pod deletion encounters
an error. Adds retry logic and error handling.

Fixes #456
```

## Changelog

We use [Keep a Changelog](https://keepachangelog.com/) format. When making changes:

1. Add your changes to the `[Unreleased]` section of `CHANGELOG.md`
2. Categorize under:
   - `Added` - new features
   - `Changed` - changes to existing functionality
   - `Deprecated` - soon-to-be removed features
   - `Removed` - removed features
   - `Fixed` - bug fixes
   - `Security` - security fixes

Example:
```markdown
## [Unreleased]

### Added
- GPU instance type selection based on pod GPU requests (#123)

### Fixed
- Instance leak when pod deletion fails (#456)
```

## Versioning

We use [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality
- **PATCH** version for backwards-compatible bug fixes

Version is tracked in the `VERSION` file.

## Pull Request Process

1. **Update Documentation**: Ensure README and docs are updated
2. **Update Changelog**: Add your changes to CHANGELOG.md
3. **Update Tests**: Add or update tests for your changes
4. **Pass CI**: Ensure all CI checks pass
5. **Code Review**: Address review comments
6. **Squash Commits**: Clean up commit history if requested

### PR Title

Use conventional commit format:
```
feat: add GPU support for p3 instances
fix: resolve memory leak in pod status sync
docs: update getting started guide
```

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How did you test this?

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] Changelog updated
- [ ] No new warnings generated
```

## Testing Guidelines

### Unit Tests

- Write unit tests for all new code
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies (AWS, Kubernetes)

Example:
```go
func TestSelectInstanceType(t *testing.T) {
    tests := []struct {
        name     string
        cpu      string
        memory   string
        expected string
    }{
        {"small", "1", "2Gi", "t3.medium"},
        {"large", "4", "8Gi", "t3.xlarge"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

- Test AWS integration with real AWS services (use test accounts)
- Test Kubernetes integration with test clusters
- Clean up resources after tests

## Documentation

### Code Documentation

- Add godoc comments for all exported functions, types, and packages
- Include examples in documentation where helpful
- Document edge cases and assumptions

Example:
```go
// SelectInstanceType chooses an appropriate EC2 instance type based on
// the resource requests in the pod specification. It considers CPU, memory,
// and GPU requirements.
//
// The selection algorithm prioritizes cost efficiency while ensuring
// resource requests can be satisfied. For pods without explicit requests,
// it defaults to t3.medium.
//
// Example:
//   instanceType := SelectInstanceType(pod)
//   fmt.Println(instanceType) // "t3.large"
func SelectInstanceType(pod *corev1.Pod) string {
    // Implementation
}
```

### User Documentation

- Update README.md for user-facing changes
- Update docs/ for detailed documentation
- Include examples and tutorials
- Keep troubleshooting section updated

## Areas for Contribution

### High Priority

- [ ] Container runtime integration (Docker/containerd)
- [ ] CloudWatch Logs integration for pod logs
- [ ] SSM Session Manager for exec functionality
- [ ] GPU instance support
- [ ] Persistent state storage (DynamoDB/etcd)
- [ ] Cost tracking and reporting
- [ ] Pod status change notifications
- [ ] Comprehensive integration tests

### Medium Priority

- [ ] Multi-region support
- [ ] Advanced instance type selection
- [ ] Network policy implementation
- [ ] EBS volume management
- [ ] Spot instance interruption handling
- [ ] Metrics and monitoring improvements
- [ ] Helm chart

### Low Priority

- [ ] Multi-cloud support (Azure, GCP)
- [ ] Web UI for management
- [ ] Cost optimization recommendations
- [ ] Auto-scaling based on metrics

## Questions or Issues?

- Open an issue on GitHub
- Reach out on NRP Slack
- Email: [your contact email]

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
