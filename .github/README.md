# GitHub Actions CI/CD Setup

This repository includes comprehensive GitHub Actions workflows for continuous integration and deployment.

## 🚀 Workflows

### 1. CI/CD Pipeline (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Jobs:**
- **Test**: Runs tests with PostgreSQL and Redis services
- **Lint**: Code quality checks with golangci-lint
- **Security**: Security scanning with Gosec
- **Build Docker**: Builds and pushes Docker images (main branch only)

### 2. Production Deployment (`.github/workflows/deploy.yml`)

**Triggers:**
- Git tags starting with `v` (e.g., `v1.0.0`)
- Manual workflow dispatch

**Features:**
- Environment-specific deployments (staging/production)
- Docker image building and pushing
- Example server deployment configuration

### 3. Security Scanning (`.github/workflows/security.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- Weekly scheduled scan (Mondays at 9 AM UTC)

**Scans:**
- Go vulnerability check with `govulncheck`
- Dependency scanning with Trivy
- Secret scanning with TruffleHog

## 🔑 Required Secrets

See [SECRETS.md](.github/SECRETS.md) for detailed information about configuring repository secrets.

### Essential Secrets (Required)
```
OPENAI_API_KEY      # OpenAI API key
GEMINI_API_KEY      # Google Gemini API key
JWT_SECRET          # JWT signing secret
```

### Production Secrets (For Deployment)
```
DOCKER_USERNAME     # Docker Hub username
DOCKER_PASSWORD     # Docker Hub password
DB_HOST            # Production database host
DB_PASS            # Production database password
# ... and more (see SECRETS.md)
```

## 🛠️ Local Development

### Setup Environment
```bash
make setup-env
# Edit .env with your values
```

### Build and Test
```bash
make build          # Build the application
make test           # Run tests
make lint           # Run linter
make docker-build   # Build Docker image
```

### Test CI Locally (Optional)
Install [act](https://github.com/nektos/act) to test GitHub Actions locally:
```bash
# Test the CI workflow
make test-ci
```

## 📋 Getting Started

1. **Configure Secrets**: Follow [SECRETS.md](.github/SECRETS.md) to set up repository secrets
2. **Create Environments**: Set up `staging` and `production` environments in GitHub
3. **Test the Pipeline**: Push a commit or create a PR to trigger the CI pipeline
4. **Deploy**: Create a git tag (e.g., `v1.0.0`) to trigger deployment

## 🔧 Customization

### Modify Workflows
- Edit workflow files in `.github/workflows/`
- Adjust environment variables and secrets as needed
- Update Docker registry and deployment targets

### Environment Configuration
- Create environment-specific secrets in GitHub Settings
- Modify deployment steps in `deploy.yml`
- Update server deployment configuration

## 📊 Monitoring

The workflows provide:
- Test coverage reports
- Security scan results in GitHub Security tab
- Build and deployment status
- Docker image tags and metadata

## 🚨 Troubleshooting

### Common Issues

1. **Missing Secrets**: Check that all required secrets are configured
2. **Build Failures**: Verify Go version and dependencies
3. **Docker Push Failures**: Ensure Docker Hub credentials are correct
4. **Test Failures**: Check database and Redis service configuration

### Debug Tips

- Check the Actions tab for detailed logs
- Review the workflow YAML for syntax errors
- Verify environment variables and secrets
- Test builds locally before pushing