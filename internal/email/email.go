package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log/slog"
	"path/filepath"
	"strconv"
	"time"

	"pdf-emailer/internal/config"
	"pdf-emailer/internal/utils"

	"gopkg.in/mail.v2"
)

type EmailService struct {
	config *config.Config
	logger *slog.Logger
}

type EmailData struct {
	RecipientName string
	SenderName    string
	Message       string
	PDFName       string
	Date          string
}

func NewEmailService(cfg *config.Config, logger *slog.Logger) *EmailService {
	return &EmailService{
		config: cfg,
		logger: logger,
	}
}

func (e *EmailService) SendEmailWithAttachment(recipientEmail, recipientName, pdfPath, templateName string) error {
	// Load email template
	templatePath := filepath.Join("assets", "templates", "email", templateName+".html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Prepare email data
	emailData := EmailData{
		RecipientName: recipientName,
		SenderName:    e.config.SMTPConfig.FromName,
		PDFName:       filepath.Base(pdfPath),
		Date:          time.Now().Format("2006-01-02"),
	}

	// Execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, emailData); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Create email message using gomail
	m := mail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", e.config.SMTPConfig.FromName, e.config.SMTPConfig.FromEmail))
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "PDF Dokument")
	m.SetBody("text/html", body.String())

	// Add PDF attachment
	m.Attach(pdfPath)

	// Convert port string to int
	port, err := strconv.Atoi(e.config.SMTPConfig.Port)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	// Create dialer
	d := mail.NewDialer(e.config.SMTPConfig.Host, port, e.config.SMTPConfig.Username, e.config.SMTPConfig.Password)

	// Configure TLS if needed
	if e.config.SMTPConfig.UseTLS {
		d.TLSConfig = &tls.Config{ServerName: e.config.SMTPConfig.Host}
	}

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	e.logger.Info("Email sent successfully", "recipient", recipientEmail, "template", templateName)
	return nil
}

// ValidateEmailAddress validates if the provided email address is properly formatted
func (e *EmailService) ValidateEmailAddress(email string) error {
	return utils.ValidateEmail(email)
}

// SendTestEmail sends a test email to verify SMTP configuration
func (e *EmailService) SendTestEmail(recipientEmail string) error {
	m := mail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", e.config.SMTPConfig.FromName, e.config.SMTPConfig.FromEmail))
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "Test Email - PDF Emailer")
	m.SetBody("text/plain", "This is a test email to verify SMTP configuration is working correctly.")

	// Convert port string to int
	port, err := strconv.Atoi(e.config.SMTPConfig.Port)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	// Create dialer
	d := mail.NewDialer(e.config.SMTPConfig.Host, port, e.config.SMTPConfig.Username, e.config.SMTPConfig.Password)

	// Configure TLS if needed
	if e.config.SMTPConfig.UseTLS {
		d.TLSConfig = &tls.Config{ServerName: e.config.SMTPConfig.Host}
	}

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send test email: %w", err)
	}

	e.logger.Info("Test email sent successfully", "recipient", recipientEmail)
	return nil
}
