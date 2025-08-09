# PDF Emailer

A modern web application built with Go, featuring HTMX for dynamic interactions, TailwindCSS for styling, and dynamic email template management. The application is designed with testability in mind and uses minimal external dependencies.

## 🚀 Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd pdf-emailer
make setup-env
make dev
```

The application will be available at `http://localhost:8080`

## ✨ Features

- **Modern Web Stack**: Go + HTMX + TailwindCSS
- **Dynamic Email Templates**: Automatically discover and validate email templates
- **File Upload**: Drag & drop PDF upload with email extraction
- **Email Service**: Robust email sending with attachments using gomail library
- **Template Management**: MJML support for responsive email templates
- **Testing**: Comprehensive unit tests with mocking
- **Security**: Built-in security headers and CORS support
- **Development**: Hot reload with Air, email testing with Mailpit

## 📚 Documentation

- **[Installation Guide](docs/installation.md)** - Complete setup instructions
- **[Email Templates Guide](docs/email-templates.md)** - Creating and managing email templates
- **[Development Guide](docs/development.md)** - Development workflow and best practices
- **[Configuration Guide](docs/configuration.md)** - Environment variables and configuration

## 🎯 Core Functionality

### Email Template System

- **Dynamic Discovery**: Automatically finds templates in `assets/templates/email/`
- **MJML Support**: Write responsive email templates using MJML
- **Template Variables**: Use `{{.RecipientName}}`, `{{.PDFName}}`, etc.
- **Validation**: Clear error messages with available template options

### File Upload & Email

- **Smart Filename Parsing**: Extract email from `filename_email@domain.com.pdf`
- **Template Selection**: Choose from available email templates
- **Email Validation**: Robust email format validation
- **PDF Attachment**: Automatically attach PDFs to emails
- **Client-side Throttling**: Uploads are limited to 5 concurrent files by default. Change `MAX_CONCURRENCY` in `assets/js/upload.js` to adjust.

### Development Experience

- **Hot Reload**: Automatic restart on file changes with Air
- **Email Testing**: Mailpit integration for local email testing
- **Comprehensive Testing**: Unit tests for all components
- **Template Management**: Easy MJML to HTML conversion

## 🛠️ Available Commands

```bash
make dev           # Start development server with hot reload
make build         # Build the application
make test          # Run all tests
make mjml          # Convert MJML templates to HTML
make setup-env     # Set up environment files
```

## 📁 Project Structure

```
pdf-emailer/
├── docs/                  # 📚 Detailed documentation
├── internal/              # 🔧 Application code
│   ├── config/           # Configuration management
│   ├── email/            # Email service
│   ├── handlers/         # HTTP handlers
│   ├── utils/            # Shared utilities
│   └── ...
├── assets/               # 🎨 Application assets
│   ├── templates/        # HTML and email templates
│   └── js/              # JavaScript files
└── uploads/             # 📁 Uploaded files (auto-created)
```

## 🔧 Configuration

The application uses environment variables for configuration. See the [Configuration Guide](docs/configuration.md) for details.

**Quick setup:**
```bash
make setup-env  # Creates .env and .env.docker files
```

## 📧 Email Templates

### Adding New Templates

1. **Create MJML file**: `assets/templates/email/template4.mjml`
2. **Convert to HTML**: `make mjml`
3. **Template is automatically available** - no code changes needed!

### Template Variables

- `{{.RecipientName}}` - Email recipient name
- `{{.SenderName}}` - Sender name
- `{{.PDFName}}` - Uploaded PDF filename
- `{{.Date}}` - Current date

See the [Email Templates Guide](docs/email-templates.md) for complete documentation.

## 🧪 Testing

```bash
make test          # Run all tests
make test-coverage # Run tests with coverage
```

## 🐳 Docker

```bash
docker compose up  # Start app + Mailpit
```

- **Application**: `http://localhost:8080`
- **Mailpit (Email Testing)**: `http://localhost:8025`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Need help?** Check the [documentation](docs/) or [create an issue](https://github.com/your-repo/pdf-emailer/issues). 