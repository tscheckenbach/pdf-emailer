package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log/slog"
	"path/filepath"
	"time"

	"pdf-emailer/internal/config"

	"gopkg.in/mail.v2"
)

type EmailService struct {
	config          *config.Config
	logger          *slog.Logger
	templateManager *TemplateManager
}

type EmailData struct {
	RecipientName string
	SenderName    string
	Message       string
	PDFName       string
	Date          string
}

func NewEmailService(cfg *config.Config, logger *slog.Logger, templateManager *TemplateManager) *EmailService {
	return &EmailService{
		config:          cfg,
		logger:          logger,
		templateManager: templateManager,
	}
}

// dialer builds an SMTP dialer from the service configuration.
func (e *EmailService) dialer() *mail.Dialer {
	d := mail.NewDialer(e.config.SMTPConfig.Host, e.config.SMTPConfig.Port, e.config.SMTPConfig.Username, e.config.SMTPConfig.Password)
	if e.config.SMTPConfig.UseTLS {
		d.TLSConfig = &tls.Config{ServerName: e.config.SMTPConfig.Host}
	}
	return d
}

func (e *EmailService) SendEmailWithAttachment(recipientEmail, recipientName, pdfPath, templateName, attachmentName string) error {
	// Validate template using template manager
	if err := e.templateManager.ValidateTemplate(templateName); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	// Load email template using template manager path
	templatePath := e.templateManager.GetTemplatePath(templateName)
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

	// Add PDF attachment with custom name
	m.Attach(pdfPath, mail.Rename(attachmentName+".pdf"))

	// Send email
	if err := e.dialer().DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	e.logger.Info("Email sent successfully", "recipient", recipientEmail, "template", templateName)
	return nil
}

// ValidateEmailAddress validates if the provided email address is properly formatted
func (e *EmailService) ValidateEmailAddress(email string) error {
	return ValidateAddress(email)
}

// SendTestEmail sends a test email to verify SMTP configuration
func (e *EmailService) SendTestEmail(recipientEmail string) error {
	m := mail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", e.config.SMTPConfig.FromName, e.config.SMTPConfig.FromEmail))
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "Test Email - PDF Emailer")
	m.SetBody("text/plain", "This is a test email to verify SMTP configuration is working correctly.")

	// Send email
	if err := e.dialer().DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send test email: %w", err)
	}

	e.logger.Info("Test email sent successfully", "recipient", recipientEmail)
	return nil
}
