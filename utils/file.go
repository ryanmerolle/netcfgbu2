package utils

import (
	"strings"
)

// Helper function to trim and ensure a single newline at the end
func ensureSingleNewline(s string) string {
	s = strings.TrimRight(s, "\n")
	return s + "\n"
}
