package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	ErrUsernameTooShort    = fmt.Errorf("username must be at least 3 characters long")
	ErrUsernameTooLong     = fmt.Errorf("username must be at most 30 characters long")
	ErrUsernameInvalidChar = fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	ErrUsernameStartEnd    = fmt.Errorf("username must start and end with a letter or number")
	ErrUsernameEmpty       = fmt.Errorf("username cannot be empty")
)

// Username regex: alphanumeric, underscore, hyphen, 3-30 chars
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{1,28}[a-zA-Z0-9]$`)

// ValidateUsername validates username format
func ValidateUsername(username string) error {
	if username == "" {
		return ErrUsernameEmpty
	}

	username = strings.TrimSpace(username)

	if len(username) < 3 {
		return ErrUsernameTooShort
	}

	if len(username) > 30 {
		return ErrUsernameTooLong
	}

	// Check regex pattern
	if !usernameRegex.MatchString(username) {
		return ErrUsernameInvalidChar
	}

	// Additional validation: must start and end with alphanumeric
	firstChar := rune(username[0])
	lastChar := rune(username[len(username)-1])

	if !unicode.IsLetter(firstChar) && !unicode.IsDigit(firstChar) {
		return ErrUsernameStartEnd
	}

	if !unicode.IsLetter(lastChar) && !unicode.IsDigit(lastChar) {
		return ErrUsernameStartEnd
	}

	// Check for consecutive special characters
	if strings.Contains(username, "__") || strings.Contains(username, "--") {
		return ErrUsernameInvalidChar
	}

	return nil
}

