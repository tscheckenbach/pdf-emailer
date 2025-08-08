package utils

import (
	"net/mail"
	"strings"
)

// ValidateEmail validates email format using Go's net/mail package for RFC 5322 compliant validation
func ValidateEmail(email string) error {
	if email == "" {
		return &EmailValidationError{Message: "email address cannot be empty"}
	}

	// Use Go's built-in mail.ParseAddress for RFC 5322 compliant validation
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return &EmailValidationError{Message: "invalid email format", OriginalError: err}
	}

	// Extract the email address (remove display name if present)
	emailAddr := addr.Address
	if emailAddr == "" {
		return &EmailValidationError{Message: "email address is empty after parsing"}
	}

	// Basic sanity checks that mail.ParseAddress might not catch
	// Check for minimum length (a@b.c minimum)
	if len(emailAddr) < 5 {
		return &EmailValidationError{Message: "email address too short"}
	}

	// Check for maximum length (RFC 5321 limit)
	if len(emailAddr) > 254 {
		return &EmailValidationError{Message: "email address too long"}
	}

	// Split by @ and validate parts
	parts := strings.Split(emailAddr, "@")
	if len(parts) != 2 {
		return &EmailValidationError{Message: "email must contain exactly one @ symbol"}
	}

	localPart := parts[0]
	domain := parts[1]

	// Validate local part length (RFC 5321 limit)
	if len(localPart) == 0 || len(localPart) > 64 {
		return &EmailValidationError{Message: "local part of email is invalid"}
	}

	// Validate domain length (RFC 5321 limit)
	if len(domain) == 0 || len(domain) > 253 {
		return &EmailValidationError{Message: "domain part of email is invalid"}
	}

	// Check for valid domain structure (must have at least one dot)
	if !strings.Contains(domain, ".") {
		return &EmailValidationError{Message: "domain must contain a dot"}
	}

	// Check for valid TLD (at least 1 character - some valid TLDs can be single character)
	domainParts := strings.Split(domain, ".")
	tld := domainParts[len(domainParts)-1]
	if len(tld) < 1 {
		return &EmailValidationError{Message: "top-level domain cannot be empty"}
	}

	// Check for invalid characters in domain
	invalidDomainChars := []string{"!", "#", "$", "%", "&", "*", "+", "=", "?", "^", "{", "|", "}", "~"}
	for _, char := range invalidDomainChars {
		if strings.Contains(domain, char) {
			return &EmailValidationError{Message: "domain contains invalid characters"}
		}
	}

	return nil
}

// IsValidEmail is a convenience function that returns a boolean instead of an error
func IsValidEmail(email string) bool {
	return ValidateEmail(email) == nil
}

// EmailValidationError represents an email validation error
type EmailValidationError struct {
	Message       string
	OriginalError error
}

func (e *EmailValidationError) Error() string {
	if e.OriginalError != nil {
		return e.Message + ": " + e.OriginalError.Error()
	}
	return e.Message
}

func (e *EmailValidationError) Unwrap() error {
	return e.OriginalError
}
