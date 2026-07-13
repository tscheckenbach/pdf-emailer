package email

import (
	"testing"
)

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		// Valid emails
		{"simple valid", "test@example.com", false},
		{"with subdomain", "test@sub.example.com", false},
		{"with plus sign", "test+tag@example.com", false},
		{"with dots", "test.name@example.com", false},
		{"with numbers", "test123@example.com", false},
		{"with hyphens", "test-name@example.com", false},
		{"with underscores", "test_name@example.com", false},
		{"short but valid", "a@b.c", false},
		{"long domain", "test@very-long-domain-name-with-many-subdomains.example.com", false},

		// Invalid emails
		{"empty", "", true},
		{"with display name", "Test User <test@example.com>", true}, // bare addresses only
		{"missing @", "testexample.com", true},
		{"missing domain", "test@", true},
		{"missing local part", "@example.com", true},
		{"too short", "a@b", true},
		{"too long", "very-long-email-address-that-exceeds-the-maximum-length-allowed-by-rfc-5321@very-long-domain-name-with-many-subdomains.example.com", true},
		{"invalid domain chars", "test@example!com", true},
		{"invalid domain chars 2", "test@example#com", true},
		{"no dot in domain", "test@example", true},
		{"single char tld", "test@example.c", false}, // Single char TLDs can be valid
		{"multiple @", "test@@example.com", true},
		{"space in local", "test name@example.com", true},
		{"space in domain", "test@example .com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddress(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAddress(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}
