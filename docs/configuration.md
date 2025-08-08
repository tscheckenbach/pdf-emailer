# Configuration Guide

## Overview

The application uses environment variables for configuration, making it easy to deploy in different environments without code changes.

## Environment Variables

### Core Configuration

| Variable | Default | Description | Required |
|----------|---------|-------------|----------|
| `SERVER_PORT` | `8080` | HTTP server port | No |
| `ENVIRONMENT` | `development` | Application environment | No |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) | No |

### SMTP Configuration

| Variable | Default | Description | Required |
|----------|---------|-------------|----------|
| `SMTP_HOST` | `localhost` | SMTP server hostname | Yes |
| `SMTP_PORT` | `1025` | SMTP server port | Yes |
| `SMTP_USERNAME` | - | SMTP username | No |
| `SMTP_PASSWORD` | - | SMTP password | No |
| `SMTP_FROM_EMAIL` | `noreply@example.com` | Sender email address | Yes |
| `SMTP_FROM_NAME` | `PDF Emailer` | Sender name | Yes |
| `SMTP_USE_TLS` | `false` | Enable TLS for SMTP | No |

## Environment Files

### File Structure

```
pdf-emailer/
├── .env.example          # Example configuration (committed)
├── .env                  # Local development (not committed)
├── .env.docker          # Docker development (not committed)
└── .env.docker.example  # Docker example (committed)
```

### File Priority

The application loads environment files in this order (first found wins):

1. `.env`
2. `.env.local`
3. `.env.docker`

### Security Note

**Never commit sensitive files to version control:**
- ❌ `.env` (contains real credentials)
- ❌ `.env.docker` (contains real credentials)
- ✅ `.env.example` (contains example values)
- ✅ `.env.docker.example` (contains example values)

## Configuration Examples

### Development Environment

```bash
# .env
SERVER_PORT=8080
ENVIRONMENT=development
LOG_LEVEL=debug

# SMTP Configuration (Mailpit for local testing)
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM_EMAIL=noreply@pdf-emailer.local
SMTP_FROM_NAME=PDF Emailer
SMTP_USE_TLS=false
```

### Production Environment

```bash
# .env
SERVER_PORT=8080
ENVIRONMENT=production
LOG_LEVEL=info

# SMTP Configuration (Gmail example)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=your-email@gmail.com
SMTP_FROM_NAME=Your Company Name
SMTP_USE_TLS=true
```

### Docker Environment

```bash
# .env.docker
SERVER_PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info

# SMTP Configuration (Mailpit in Docker)
SMTP_HOST=mailpit
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM_EMAIL=noreply@pdf-emailer.local
SMTP_FROM_NAME=PDF Emailer
SMTP_USE_TLS=false
```

## SMTP Providers

### Gmail

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_USE_TLS=true
```

**Note**: Use App Passwords, not your regular Gmail password.

### Outlook/Hotmail

```bash
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USERNAME=your-email@outlook.com
SMTP_PASSWORD=your-password
SMTP_USE_TLS=true
```

### SendGrid

```bash
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
SMTP_USE_TLS=true
```

### Mailgun

```bash
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USERNAME=your-mailgun-username
SMTP_PASSWORD=your-mailgun-password
SMTP_USE_TLS=true
```

### Custom SMTP Server

```bash
SMTP_HOST=your-smtp-server.com
SMTP_PORT=587
SMTP_USERNAME=your-username
SMTP_PASSWORD=your-password
SMTP_USE_TLS=true
```

## Logging Configuration

### Log Levels

| Level | Description | Use Case |
|-------|-------------|----------|
| `debug` | Detailed debug information | Development |
| `info` | General information | Production |
| `warn` | Warning messages | Production |
| `error` | Error messages only | Troubleshooting |

### Log Format

The application uses structured logging with `log/slog`. Logs include:

- Timestamp
- Log level
- Message
- Structured fields (key-value pairs)

Example:
```json
{"time":"2025-01-15T10:30:00Z","level":"INFO","msg":"server started","port":"8080"}
```

## Environment-Specific Settings

### Development

```bash
ENVIRONMENT=development
LOG_LEVEL=debug
```

Features:
- Detailed logging
- Development error messages
- Mailpit for email testing
- Hot reload with Air

### Production

```bash
ENVIRONMENT=production
LOG_LEVEL=info
```

Features:
- Minimal logging
- Production error handling
- Real SMTP server
- Optimized performance

### Testing

```bash
ENVIRONMENT=testing
LOG_LEVEL=warn
```

Features:
- Test-specific configuration
- Minimal logging
- Mock services

## Configuration Validation

The application validates configuration on startup:

### Required Fields

- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_FROM_EMAIL`
- `SMTP_FROM_NAME`

### Optional Validation

- SMTP connection test (if `ENVIRONMENT=production`)
- Template directory existence
- Upload directory permissions

## Environment Setup Scripts

### Quick Setup

```bash
# Set up both local and Docker environments
./scripts/setup-env.sh -a

# Set up only Docker environment
./scripts/setup-env.sh -d

# Set up only local environment
./scripts/setup-env.sh -l

# Show help
./scripts/setup-env.sh -h
```

### Manual Setup

```bash
# Copy example files
cp env.example .env
cp env.example .env.docker

# Edit configuration
nano .env
nano .env.docker
```

## Docker Configuration

### Docker Compose

The application includes a `docker-compose.yml` file with:

- Application container
- Mailpit for email testing
- Volume mounts for persistence
- Environment file support

### Environment Variables in Docker

```bash
# Use environment file
docker compose --env-file .env.docker up

# Override specific variables
SMTP_HOST=custom-host docker compose up

# Use environment file with overrides
docker compose --env-file .env.docker -e SMTP_HOST=custom-host up
```

## Security Considerations

### Environment Variables

1. **Never commit sensitive data**
   ```bash
   # Add to .gitignore
   .env
   .env.docker
   .env.local
   ```

2. **Use strong passwords**
   ```bash
   # Generate secure password
   openssl rand -base64 32
   ```

3. **Rotate credentials regularly**
   - Update SMTP passwords periodically
   - Use app-specific passwords
   - Monitor for suspicious activity

### SMTP Security

1. **Use TLS when possible**
   ```bash
   SMTP_USE_TLS=true
   ```

2. **Use dedicated SMTP accounts**
   - Don't use personal email accounts
   - Create service accounts for applications

3. **Monitor SMTP usage**
   - Check email sending logs
   - Monitor for unusual activity

## Troubleshooting

### Common Configuration Issues

1. **SMTP Connection Failed**
   ```bash
   # Check SMTP settings
   cat .env | grep SMTP
   
   # Test SMTP connection
   telnet smtp.gmail.com 587
   ```

2. **Port Already in Use**
   ```bash
   # Check port usage
   lsof -i :8080
   
   # Change port
   SERVER_PORT=8081
   ```

3. **Template Not Found**
   ```bash
   # Check template directory
   ls -la assets/templates/email/
   
   # Check environment
   echo $ENVIRONMENT
   ```

### Configuration Validation

```bash
# Validate configuration
go run main.go --validate-config

# Check environment variables
env | grep -E "(SERVER|SMTP|LOG)"
```

### Getting Help

- Check the [Installation Guide](installation.md)
- Review the [Development Guide](development.md)
- See the [Email Templates Guide](email-templates.md)
- Check application logs for error messages
