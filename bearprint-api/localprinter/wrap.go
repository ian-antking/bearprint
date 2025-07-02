package localprinter

import (
	"strings"
	"unicode/utf8"
)

func wrapText(text string, width int) []string {
	if text == "" {
		return []string{} // Return empty slice, not nil
	}

	var lines []string
	var current string
	for _, word := range strings.Fields(text) {
		if current == "" {
			current = word
			continue
		}
		if utf8.RuneCountInString(current+" "+word) > width {
			lines = append(lines, current)
			current = word
		} else {
			current += " " + word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
