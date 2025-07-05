# GitHub Actions Secrets Configuration

This document outlines the required secrets that need to be configured in your GitHub repository for the CI/CD pipelines to work properly.

## Required Secrets

### Go to GitHub Repository Settings > Secrets and variables > Actions > Repository secrets

Add the following secrets:

### 🔑 Required API Keys
These are critical for the application to function:

- **`OPENAI_API_KEY`** (Required)
  - Your OpenAI API key
  - Get it from: https://platform.openai.com/api-keys
  - The application will panic without this

- **`GEMINI_API_KEY`** (Required)
  - Your Google Gemini API key
  - Get it from: https://ai.google.dev/
  - The application will panic without this

### 🔐 Security Configuration

- **`JWT_SECRET`** (Recommended)
  - A strong secret key for JWT token signing
  - Generate with: `openssl rand -base64 32`
  - Default: "secret" (not secure for production)

### 🗄️ Database Configuration (for production deployment)

- **`DB_HOST`** - Database host (e.g., `your-db-host.com`)
- **`DB_PORT`** - Database port (default: `5432`)
- **`DB_USER`** - Database username
- **`DB_PASS`** - Database password
- **`DB_NAME`** - Database name
- **`DB_SSL_MODE`** - SSL mode (e.g., `require` for production, `disable` for local)

### 📦 Redis Configuration (for production deployment)

- **`REDIS_HOST`** - Redis host
- **`REDIS_PORT`** - Redis port (default: `6379`)

### 🗂️ Minio Object Storage Configuration

- **`MINIO_ENDPOINT`** - Minio endpoint (e.g., `s3.amazonaws.com` for AWS S3)
- **`MINIO_ACCESS_KEY`** - Access key for Minio/S3
- **`MINIO_SECRET_KEY`** - Secret key for Minio/S3
- **`MINIO_REGION`** - Region (default: `us-east-1`)
- **`MINIO_USE_SSL`** - Use SSL (default: `true` for production)

### 🐳 Docker Hub Configuration (for deployment)

- **`DOCKER_USERNAME`** - Your Docker Hub username
- **`DOCKER_PASSWORD`** - Your Docker Hub password or access token

### 🚀 Deployment Configuration (optional)

If using the deployment workflow:

- **`DEPLOY_HOST`** - Server hostname/IP for deployment
- **`DEPLOY_USER`** - SSH username for deployment
- **`DEPLOY_KEY`** - SSH private key for deployment

## How to Add Secrets

1. Go to your GitHub repository
2. Click on **Settings** tab
3. In the left sidebar, click **Secrets and variables** > **Actions**
4. Click **New repository secret**
5. Add the secret name and value
6. Click **Add secret**

## Environment-Specific Secrets

The workflows support different environments:

- **CI/CD Pipeline**: Uses test database and minimal configuration
- **Staging**: Uses staging environment secrets
- **Production**: Uses production environment secrets

You can create environment-specific secrets by:

1. Go to **Settings** > **Environments**
2. Create environments (`staging`, `production`)
3. Add environment-specific secrets

## Security Best Practices

1. **Never commit secrets to the repository**
2. **Use strong, unique passwords and keys**
3. **Rotate secrets regularly**
4. **Use environment-specific secrets for different deployments**
5. **Enable two-factor authentication on all service accounts**
6. **Use the principle of least privilege for API keys**

## Testing Locally

For local development, copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
# Edit .env with your local values
```

## Troubleshooting

### Common Issues:

1. **Build fails with "OpenAI API key not found"**
   - Ensure `OPENAI_API_KEY` is set in repository secrets

2. **Build fails with "Gemini API key not found"**
   - Ensure `GEMINI_API_KEY` is set in repository secrets

3. **Docker push fails**
   - Verify `DOCKER_USERNAME` and `DOCKER_PASSWORD` are correctly set

4. **Database connection fails in tests**
   - Check if database service is running properly in GitHub Actions

## Validation

You can validate your secrets setup by:

1. Pushing a commit to trigger the CI pipeline
2. Checking the Actions tab for any secret-related errors
3. Reviewing the build logs for missing environment variables