package email

import (
	"testing"

	"pdf-emailer/internal/config"
)

func TestNewEmailService(t *testing.T) {
	cfg := &config.Config{
		SMTPConfig: config.SMTPConfig{
			Host:      "localhost",
			Port:      "1025",
			Username:  "test",
			Password:  "test",
			FromEmail: "test@example.com",
			FromName:  "Test Sender",
			UseTLS:    false,
		},
	}

	service := NewEmailService(cfg, nil)
	if service == nil {
		t.Fatal("Expected email service to be created")
	}

	if service.config != cfg {
		t.Error("Expected config to be set correctly")
	}
}

func TestValidateEmailAddress(t *testing.T) {
	cfg := &config.Config{}
	service := NewEmailService(cfg, nil)

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"valid email with subdomain", "test@sub.example.com", false},
		{"valid email with plus", "test+tag@example.com", false},
		{"empty email", "", true},
		{"missing @", "testexample.com", true},
		{"missing domain", "test@", true},
		{"missing local part", "@example.com", true},
		{"invalid domain chars", "test@example!com", true},
		{"too short", "a@b", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateEmailAddress(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmailAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
