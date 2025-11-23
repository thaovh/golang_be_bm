package validator

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrEmailInvalid      = fmt.Errorf("invalid email format")
	ErrEmailTooLong      = fmt.Errorf("email must be at most 255 characters long")
	ErrEmailEmpty        = fmt.Errorf("email cannot be empty")
)

// Email regex pattern (RFC 5322 simplified)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmailEmpty
	}

	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	if len(email) > 255 {
		return ErrEmailTooLong
	}

	if !emailRegex.MatchString(email) {
		return ErrEmailInvalid
	}

	// Additional checks
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ErrEmailInvalid
	}

	localPart := parts[0]
	domain := parts[1]

	// Local part validation
	if len(localPart) == 0 || len(localPart) > 64 {
		return ErrEmailInvalid
	}

	// Domain validation
	if len(domain) == 0 || len(domain) > 255 {
		return ErrEmailInvalid
	}

	// Check for consecutive dots
	if strings.Contains(email, "..") {
		return ErrEmailInvalid
	}

	// Check for leading/trailing dots
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return ErrEmailInvalid
	}

	return nil
}

