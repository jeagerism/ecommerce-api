package util

import (
	"fmt"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type ValidationError struct {
	Fields []string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(v.Fields, "; "))
}

func ValidateCreateUser(username, email, password string) error {
	var errs []string

	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if len(username) < 3 || len(username) > 20 {
		errs = append(errs, "username must be between 3 to 20 characters")
	}
	if !emailRegex.MatchString(email) {
		errs = append(errs, "invalid email format")
	}
	if len(password) < 6 {
		errs = append(errs, "password must be at least 6 characters")
	}

	if len(errs) > 0 {
		return &ValidationError{Fields: errs}
	}
	return nil
}

func ValidateLoginInput(email, password string) error {
	var errs []string

	email = strings.TrimSpace(email)

	if !emailRegex.MatchString(email) {
		errs = append(errs, "invalid email format")
	}
	if len(password) < 6 {
		errs = append(errs, "password must be at least 6 characters")
	}

	if len(errs) > 0 {
		return &ValidationError{Fields: errs}
	}
	return nil
}
