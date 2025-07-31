// Package internal provides utility functions and types used internally by CupSleeve.
// These are implementation details and should not be used directly by external code.
package internal

import "strings"

// WrapText wraps the given text to fit within the specified width.
// It breaks text at word boundaries and returns a slice of lines.
// If width is 0 or negative, returns the original text as a single line.
func WrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			// First word on the line
			currentLine.WriteString(word)
		} else if currentLine.Len()+1+len(word) <= width {
			// Word fits on current line with a space
			currentLine.WriteString(" " + word)
		} else {
			// Word doesn't fit, start a new line
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	// Add the last line if it has content
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// TruncateString truncates a string to the specified length.
// If the string is longer than the specified length, it's cut off and "..." is appended.
// If length is 3 or less, no ellipsis is added to avoid making the string longer.
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	if length <= 3 {
		return s[:length]
	}
	return s[:length-3] + "..."
}
