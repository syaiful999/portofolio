package utils

import (
	"regexp"
	"strings"
)

// Regex patterns for common fields
var (
	RegexNIK    = regexp.MustCompile(`^[0-9]{16}$`)
	RegexEmail  = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	RegexPhone  = regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	RegexAlpha  = regexp.MustCompile(`^[a-zA-Z\s]+$`)
	RegexAlphanum = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
)

// ValidateInput checks if the input matches a given regex pattern.
func ValidateInput(input string, pattern *regexp.Regexp) bool {
	return pattern.MatchString(input)
}

// SanitizeInput removes potentially dangerous characters from a string for basic XSS prevention.
func SanitizeInput(input string) string {
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#x27;")
	input = strings.ReplaceAll(input, "&", "&amp;")
	return input
}

// IsWhitelisted checks if a string is within a given slice of allowed values.
func IsWhitelisted(input string, whitelist []string) bool {
	for _, val := range whitelist {
		if input == val {
			return true
		}
	}
	return false
}
