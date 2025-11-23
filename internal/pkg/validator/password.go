package validator

import (
	"fmt"
	"unicode"
)

var (
	ErrPasswordTooShort      = fmt.Errorf("password must be at least 8 characters long")
	ErrPasswordTooLong       = fmt.Errorf("password must be at most 128 characters long")
	ErrPasswordNoUpper       = fmt.Errorf("password must contain at least one uppercase letter")
	ErrPasswordNoLower       = fmt.Errorf("password must contain at least one lowercase letter")
	ErrPasswordNoDigit       = fmt.Errorf("password must contain at least one digit")
	ErrPasswordNoSpecial     = fmt.Errorf("password must contain at least one special character")
	ErrPasswordCommon        = fmt.Errorf("password is too common or weak")
)

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	if len(password) > 128 {
		return ErrPasswordTooLong
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasDigit {
		return ErrPasswordNoDigit
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	// Check for common weak passwords
	if isCommonPassword(password) {
		return ErrPasswordCommon
	}

	return nil
}

// isCommonPassword checks if password is in common password list
func isCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "12345678", "123456789", "1234567890",
		"qwerty", "abc123", "password123", "admin123",
		"letmein", "welcome", "monkey", "1234567",
		"sunshine", "princess", "dragon", "passw0rd",
		"master", "hello", "freedom", "whatever",
		"qazwsx", "trustno1", "jordan23", "harley",
		"shadow", "superman", "michael", "football",
	}

	lowerPassword := password
	for _, common := range commonPasswords {
		if lowerPassword == common {
			return true
		}
	}

	return false
}

// ValidatePasswordStrength returns password strength level (weak, medium, strong)
func ValidatePasswordStrength(password string) string {
	if err := ValidatePassword(password); err != nil {
		return "weak"
	}

	// Additional checks for strength
	hasMultipleSpecial := 0
	hasMultipleDigits := 0
	hasMultipleUpper := 0
	hasMultipleLower := 0

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasMultipleUpper++
		case unicode.IsLower(char):
			hasMultipleLower++
		case unicode.IsDigit(char):
			hasMultipleDigits++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasMultipleSpecial++
		}
	}

	// Strong: length >= 12 and multiple character types
	if len(password) >= 12 && hasMultipleSpecial >= 2 && hasMultipleDigits >= 2 {
		return "strong"
	}

	// Medium: meets basic requirements
	return "medium"
}

