package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// HashPassword hashes a password using SHA-256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// ParseCommand parses a command string into a command and arguments
func ParseCommand(input string) (string, string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", ""
	}

	parts := strings.SplitN(input, " ", 2)
	command := strings.ToLower(parts[0])
	
	var args string
	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}
	
	return command, args
}

// FormatMessage formats a message with substitutions
// $n = character name, $N = character name with title
// $m = victim name, $M = victim name with title
// $s = he/she, $S = He/She, $e = his/her, $E = His/Her
func FormatMessage(message string, actor, victim, object interface{}) string {
	// TODO: Implement message formatting with substitutions
	return message
}

// Colorize adds ANSI color codes to a string
func Colorize(text string, color string) string {
	colors := map[string]string{
		"reset":     "\033[0m",
		"black":     "\033[30m",
		"red":       "\033[31m",
		"green":     "\033[32m",
		"yellow":    "\033[33m",
		"blue":      "\033[34m",
		"magenta":   "\033[35m",
		"cyan":      "\033[36m",
		"white":     "\033[37m",
		"bold":      "\033[1m",
		"underline": "\033[4m",
	}
	
	if code, ok := colors[color]; ok {
		return fmt.Sprintf("%s%s%s", code, text, colors["reset"])
	}
	
	return text
}
