package utils

import (
	"strings"
	"unicode"
	"errors"

	"s2p_minichat/config"
)

var (
	ErrEmptyMessage   = errors.New("message is empty")
	ErrMessageTooLong = errors.New("message exceeds max length")
	ErrInvalidUsername = errors.New("invalid username")
	ErrWeakPassword    = errors.New("password too weak")
)

// Check if message is valid
func ValidateMessage(msg string) error {
	msg = strings.TrimSpace(msg)

	if msg == "" {
		return ErrEmptyMessage
	}

	if len(msg) > config.MaxMessageLength {
		return ErrMessageTooLong
	}

	return nil
}

// Optional sanitization (basic safe cleanup)
func SanitizeMessage(msg string) string {
	// Remove control characters (non-printable)
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, msg)
}

func ValidateCredentials(username, password string) error {
	username = strings.TrimSpace(username)

	if len(username) < 3 {
		return ErrInvalidUsername
	}

	if len(password) < 6 {
		return ErrWeakPassword
	}

	return nil
}