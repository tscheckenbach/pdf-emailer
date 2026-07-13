package email

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

// ValidateAddress checks that the given string is a bare, RFC 5322 compliant
// email address (no display name, no surrounding whitespace).
func ValidateAddress(email string) error {
	if email == "" {
		return errors.New("email address must not be empty")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}

	// Reject display names, quoted forms and surrounding whitespace: the
	// parsed address must be exactly the input.
	if addr.Address != email {
		return errors.New("email address must be a bare address without display name")
	}

	// RFC 5321 length limits
	if len(email) > 254 {
		return errors.New("email address exceeds 254 characters")
	}

	local, domain, _ := strings.Cut(email, "@")
	if len(local) > 64 {
		return errors.New("local part of email exceeds 64 characters")
	}

	if !strings.Contains(domain, ".") {
		return errors.New("email domain must contain a dot")
	}

	for _, r := range domain {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '.':
		default:
			return fmt.Errorf("email domain contains invalid character %q", r)
		}
	}

	return nil
}

// IsValidAddress is a convenience wrapper returning a boolean instead of an error.
func IsValidAddress(email string) bool {
	return ValidateAddress(email) == nil
}
