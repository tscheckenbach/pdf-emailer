# Development Guide

## Getting Started

### Prerequisites

- Go 1.24.5 or later
- Docker and Docker Compose
- Git
- Node.js (for MJML templates)

### Development Environment

```bash
# Clone the repository
git clone <repository-url>
cd pdf-emailer

# Install dependencies
go mod tidy

# Set up environment
make setup-env

# Start development server
make dev
```

## Project Structure

```
pdf-emailer/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── env.example             # Example environment variables
├── README.md              # Main documentation
├── docs/                  # Detailed documentation
│   ├── installation.md    # Installation guide
│   ├── email-templates.md # Email templates guide
│   ├── development.md     # This file
│   └── configuration.md   # Configuration guide
├── internal/              # Internal application code
│   ├── config/           # Configuration management
│   ├── email/            # Email service with SMTP support
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── server/           # HTTP server setup
│   └── utils/            # Shared utilities (email validation, template management, etc.)
├── assets/               # Application assets
│   ├── templates/        # HTML and email templates
│   │   ├── html/         # Web page templates
│   │   │   ├── base.html # Base template with layout
│   │   │   └── home.html # Home page template
│   │   └── email/        # Email templates (MJML + HTML)
│   │       ├── template1.mjml  # MJML source files
│   │       ├── template1.html  # Compiled HTML templates
│   │       ├── template2.mjml
│   │       ├── template2.html
│   │       └── ...
│   └── js/               # JavaScript files
│       └── upload.js     # Upload functionality
└── uploads/              # Uploaded PDF files (created automatically)
```

## Development Workflow

### 1. Starting Development

```bash
# Start development environment
make dev
```

This will:
- Start Mailpit for email testing
- Start Air for auto-reload development
- Build and run the application

### 2. Making Changes

The application uses Air for hot reloading. When you make changes to Go files, the application will automatically restart.

### 3. Testing Changes

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./internal/handlers
```

### 4. Building

```bash
# Build the application
make build

# Run the built application
./bin/pdf-emailer
```

## Adding New Features

### Adding New Routes

1. **Add the handler method** in `internal/handlers/handlers.go`
2. **Register the route** in `internal/server/server.go`
3. **Add tests** in `internal/handlers/handlers_test.go`

Example:
```go
// In handlers.go
func (h *Handler) NewFeature(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// In server.go
router.Get("/new-feature", handlers.NewFeature)
```

### Adding New Web Templates

1. **Create the template file** in `assets/templates/html/`
2. **Update the handler** to use the template
3. **Add any necessary data structures**

Example:
```go
data := map[string]interface{}{
    "Title": "New Feature",
    "Data":  someData,
}
h.tmpl.ExecuteTemplate(w, "new-feature", data)
```

### Adding New Email Templates

See the [Email Templates Guide](email-templates.md) for detailed instructions.

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/handlers

# Run specific test
go test -run TestUpload ./internal/handlers
```

### Test Structure

- **Handlers**: Test HTTP handlers
- **Middleware**: Test middleware functionality
- **Configuration**: Test configuration loading and validation
- **Email Service**: Test email sending functionality
- **Utils**: Test shared utilities like email validation

### Writing Tests

Example test structure:
```go
func TestNewFeature(t *testing.T) {
    // Setup
    handler := &Handler{
        logger: logger,
        tmpl:   template,
    }
    
    // Test
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    
    handler.NewFeature(w, req)
    
    // Assertions
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

## Code Quality

### Linting

```bash
# Run linter
make lint

# Install linter (if not already installed)
make install-tools
```

### Code Formatting

```bash
# Format Go code
go fmt ./...

# Run go vet
go vet ./...
```

## Debugging

### Logging

The application uses structured logging with `log/slog`. Log levels can be configured via `LOG_LEVEL` environment variable:

```bash
# Set log level
export LOG_LEVEL=debug
```

### Email Testing

Use Mailpit for email testing:

1. **Start Mailpit**: `docker compose up mailpit -d`
2. **Upload a PDF** through the web interface
3. **Check emails** at `http://localhost:8025`

### Common Debugging Commands

```bash
# Check if application is running
curl http://localhost:8080/health

# Check available templates
ls -la assets/templates/email/

# Check logs
docker compose logs app

# Check Mailpit
curl http://localhost:8025/api/v1/messages
```

## Performance

### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./...

# Memory profiling
go test -memprofile=mem.prof ./...

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```
 
### Upload Concurrency (Client-Side Throttling)

To prevent server overload, file uploads are throttled on the client. The application uploads at most 5 files concurrently, proceeding in batches.

- Default limit: 5 concurrent uploads
- Location: `assets/js/upload.js`
- Identifier: `MAX_CONCURRENCY`

To change the limit, edit `MAX_CONCURRENCY` in `uploadFiles`:

```javascript
// assets/js/upload.js
const MAX_CONCURRENCY = 5; // Adjust as needed
for (let i = 0; i < filesToUpload.length; i += MAX_CONCURRENCY) {
  const chunk = filesToUpload.slice(i, i + MAX_CONCURRENCY);
  await Promise.all(chunk.map(file => uploadSingle(file)));
}
```

Note: This is a client-side throttle and does not change server-side limits.

### Benchmarks

```bash
# Run benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkUpload ./internal/handlers
```

## Deployment

### Local Build

```bash
# Build for local platform
make build

# Run locally
./bin/pdf-emailer
```

### Docker Build

```bash
# Build Docker image
docker build -t pdf-emailer .

# Run Docker container
docker run -p 8080:8080 pdf-emailer
```

### Production Considerations

1. **Environment Variables**: Set production environment variables
2. **Logging**: Configure appropriate log levels
3. **Security**: Review security headers and CORS settings
4. **Monitoring**: Add health checks and monitoring
5. **Backup**: Ensure uploads directory is backed up

## Troubleshooting

### Common Issues

1. **Port already in use**
   ```bash
   lsof -i :8080
   kill -9 <PID>
   ```

2. **Template not loading**
   ```bash
   # Check template files
   ls -la assets/templates/
   
   # Rebuild templates
   make mjml
   ```

3. **Email not sending**
   ```bash
   # Check SMTP configuration
   cat .env | grep SMTP
   
   # Test with Mailpit
   docker compose up mailpit -d
   ```

4. **Go version issues**
   ```bash
   # Use correct Go version
   ./scripts/use-go-1.24.sh
   
   # Check version
   go version
   ```

### Getting Help

- Check the [Installation Guide](installation.md)
- Review the [Email Templates Guide](email-templates.md)
- See the [Configuration Guide](configuration.md)
- Check GitHub Issues for known problems

## Contributing

### Code Style

- Follow Go conventions
- Use meaningful variable names
- Add comments for complex logic
- Write tests for new features

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

### Commit Messages

Use conventional commit messages:
```
feat: add new email template support
fix: resolve template loading issue
docs: update installation guide
test: add tests for upload handler
```
