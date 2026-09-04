package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lunemis/mux/tmux"
)

// renderPreview renders captured terminal cells across the exact canvas. It
// keeps the bottom-left visible when the source is larger and pads above/right
// when it is smaller; preview content is never wrapped.
func renderPreview(captured string, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	var capturedLines []string
	if captured != "" {
		capturedLines = strings.Split(captured, "\n")
	}
	if len(capturedLines) > height {
		capturedLines = capturedLines[len(capturedLines)-height:]
	}

	lines := make([]string, height)
	for i := range lines {
		lines[i] = strings.Repeat(" ", width)
	}
	start := height - len(capturedLines)
	for i, line := range capturedLines {
		lines[start+i] = padOrTruncate(line, width)
	}
	return strings.Join(lines, "\n")
}

func formatTokenLine(u *tmux.TokenUsage, width int) string {
	text := fmt.Sprintf("  %s in / %s out  ~$%.2f",
		tmux.FormatTokens(u.InputTokens),
		tmux.FormatTokens(u.OutputTokens),
		u.TotalCost)
	return lipgloss.NewStyle().Foreground(colorMuted).Render(
		padOrTruncate(text, width))
}
