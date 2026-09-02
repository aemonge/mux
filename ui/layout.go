package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

const (
	stackedHelpRows      = 2
	stackedSeparatorRows = 1
	minimumStackedHeight = 12
)

type stackedHeights struct {
	preview   int
	separator int
	list      int
	help      int
}

// calculateStackedHeights gives the upper half to preview plus its separator,
// reserves exactly two help rows, and gives the remainder to selection.
func calculateStackedHeights(total int) stackedHeights {
	if total <= 0 {
		return stackedHeights{}
	}

	help := min(stackedHelpRows, total)
	upper := total / 2
	separator := 0
	if upper > 0 {
		separator = stackedSeparatorRows
	}
	preview := max(0, upper-separator)
	list := max(0, total-upper-help)

	return stackedHeights{
		preview:   preview,
		separator: separator,
		list:      list,
		help:      help,
	}
}

func renderSeparator(width int) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().
		Foreground(colorSeparator).
		Render(strings.Repeat("━", width))
}

// padOrTruncate ensures a string is exactly `width` visible characters
func padOrTruncate(s string, width int) string {
	w := ansi.StringWidth(s)
	if w > width {
		return ansi.Truncate(s, width, "")
	}
	if w < width {
		return s + strings.Repeat(" ", width-w)
	}
	return s
}

// fixedBox takes rendered content and forces it to exactly width x height visible area.
// It splits by newlines, truncates/pads each line to width, and truncates/pads to height lines.
func fixedBox(content string, width, height int) string {
	lines := strings.Split(content, "\n")

	result := make([]string, height)
	for i := 0; i < height; i++ {
		if i < len(lines) {
			result[i] = padOrTruncate(lines[i], width)
		} else {
			result[i] = strings.Repeat(" ", width)
		}
	}
	return strings.Join(result, "\n")
}

// drawBorder wraps content lines with a rounded border
func drawBorder(content string, width, height int) string {
	innerWidth := width - 2
	lines := strings.Split(content, "\n")

	// Build bordered output
	result := make([]string, 0, height+2)

	// Top border
	result = append(result, "╭"+strings.Repeat("─", innerWidth)+"╮")

	// Content lines (pad/truncate to exactly height)
	for i := 0; i < height; i++ {
		line := ""
		if i < len(lines) {
			line = lines[i]
		}
		line = padOrTruncate(line, innerWidth)
		result = append(result, "│"+line+"│")
	}

	// Bottom border
	result = append(result, "╰"+strings.Repeat("─", innerWidth)+"╯")

	return lipgloss.NewStyle().
		Foreground(colorBorder).
		Render(strings.Join(result, "\n"))
}
