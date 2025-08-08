# PDF Emailer

A modern web application built with Go, featuring HTMX for dynamic interactions, TailwindCSS for styling, and PostgreSQL for data storage. The application is designed with testability in mind and uses minimal external dependencies.

## Features

- **Modern Web Stack**: Go + HTMX + TailwindCSS
- **HTTP Router**: Chi router with middleware
- **Logging**: Structured logging with slog
- **Testing**: Comprehensive unit tests with mocking
- **Configuration**: Environment-based configuration
- **Security**: Built-in security headers and CORS support
- **Email Service**: Robust email sending with attachments using gomail library

## Project Structure

```
pdf-emailer/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── env.example             # Example environment variables
├── README.md              # This file
├── internal/              # Internal application code
│   ├── config/           # Configuration management
│   ├── email/            # Email service with SMTP support
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── server/           # HTTP server setup
│   └── utils/            # Shared utilities (email validation, etc.)
├── templates/            # HTML templates
│   ├── base.html         # Base template with layout
│   ├── home.html         # Home page template
│   └── hello.html        # Hello World example template
└── static/               # Static assets (CSS, JS, images)
```

## Prerequisites

- Go 1.24.5 or later (recommended: use Homebrew version)
- Docker and Docker Compose (for containerized development)
- Git

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd pdf-emailer
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   # Quick setup (recommended)
   make setup-env
   
   # Or manually:
   # For local development
   cp env.example .env
   # Edit .env with your configuration
   
   # For Docker development
   cp env.example .env.docker
   # Edit .env.docker with your configuration
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

The application will be available at `http://localhost:8080`

## Email Testing with Mailpit

For local development and testing, the application includes [Mailpit](https://github.com/axllent/mailpit) - a modern email testing tool that captures all outgoing emails and provides a web interface to view them.

### Mailpit Features

- **SMTP Server**: Captures all outgoing emails on port 1025
- **Web Interface**: View emails at `http://localhost:8025`
- **No Configuration**: Works out of the box with the provided Docker setup
- **Email Preview**: View HTML, text, and attachments
- **Search**: Find specific emails quickly

### Using Mailpit

1. **Start the application with Docker:**
   ```bash
   docker compose up
   ```

2. **Upload a PDF file** through the web interface at `http://localhost:8080`

3. **View the sent email** at `http://localhost:8025`

4. **Test different email templates** by selecting different templates during upload

### Mailpit Configuration

The Docker Compose setup automatically configures:
- SMTP server on port 1025
- Web interface on port 8025
- Accepts any authentication (for development)
- No TLS required (for local testing)

## Configuration

The application uses environment variables for configuration:

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `ENVIRONMENT` | `development` | Application environment |
| `LOG_LEVEL` | `info` | Logging level |

### Environment Files

- `.env` - Local development environment variables
- `.env.docker` - Docker environment variables (not committed to git)
- `.env.docker.example` - Example Docker environment file (committed to git)

**Security Note**: Never commit `.env` or `.env.docker` files to version control as they may contain sensitive information like passwords and API keys.

### Quick Setup

Use the provided setup script to quickly configure your environment:

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

## API Endpoints

### Web Routes
- `GET /` - Home page
- `GET /hello` - Hello World example with HTMX

### API Routes
- `GET /api/health` - Health check (JSON)
- `GET /api/health/json` - Health check (JSON format)
- `GET /api/counter` - Get counter value
- `POST /api/counter/increment` - Increment counter
- `POST /api/counter/decrement` - Decrement counter
- `GET /api/dynamic-content` - Get dynamic content for HTMX

## Testing

The application includes comprehensive unit tests:

```bash
# Run all tests (uses Go 1.24.5 automatically)
make test

# Run tests with coverage
make test-coverage

# Run tests manually (ensure Go 1.24.5 is in PATH)
go test ./...

# Switch to Go 1.24.5 for this session
./scripts/use-go-1.24.sh
```

**Note**: The application requires Go 1.24.5+ for the `log/slog` package. If you have multiple Go versions installed, use the provided scripts to ensure the correct version is used.

### Test Structure

- **Handlers**: Test HTTP handlers
- **Middleware**: Test middleware functionality
- **Configuration**: Test configuration loading and validation
- **Email Service**: Test email sending functionality
- **Utils**: Test shared utilities like email validation

## Development

### Adding New Routes

1. Add the handler method in `internal/handlers/handlers.go`
2. Register the route in `internal/server/server.go`
3. Add tests in `internal/handlers/handlers_test.go`

### Adding New Templates

1. Create the template file in `templates/`
2. Update the handler to use the template
3. Add any necessary data structures



## Architecture

### Design Principles

- **Testability**: All components are designed for easy testing
- **Minimal Dependencies**: Uses standard library where possible
- **Structured Logging**: Uses slog for consistent logging
- **Security First**: Built-in security headers and CORS

### Key Components

- **Chi Router**: Fast, lightweight HTTP router
- **HTMX**: Dynamic HTML without JavaScript
- **TailwindCSS**: Utility-first CSS framework
- **slog**: Structured logging from Go standard library
- **gomail**: Robust email sending library



## Deployment

### Docker (Recommended)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Environment Variables for Production

```bash
ENVIRONMENT=production
SERVER_PORT=8080
LOG_LEVEL=info
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 