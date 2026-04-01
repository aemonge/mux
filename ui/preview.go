package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lunemis/mux/tmux"
)

func renderPreview(session *tmux.Session, captured string, width, height int) string {
	innerWidth := width - 2
	innerHeight := height - 2

	if session == nil {
		lines := make([]string, innerHeight)
		mid := innerHeight / 2
		msg := "No session selected"
		for i := range lines {
			if i == mid {
				lines[i] = padOrTruncate(centerText(msg, innerWidth), innerWidth)
			} else {
				lines[i] = strings.Repeat(" ", innerWidth)
			}
		}
		content := strings.Join(lines, "\n")
		return drawBorder(content, width, innerHeight)
	}

	// Header
	badge := aiLabel(session.ActiveCommand)
	header := fmt.Sprintf("[ %s ]  %s%s", session.Name, shortenPath(session.Directory), badge)
	headerStyled := lipgloss.NewStyle().Foreground(colorAccent).Bold(true).Render(
		padOrTruncate(header, innerWidth))
	separator := lipgloss.NewStyle().Foreground(colorBorder).Render(
		strings.Repeat("─", innerWidth))

	// Available lines for content (minus header + separator)
	contentLines := innerHeight - 2
	if contentLines < 1 {
		contentLines = 1
	}

	capLines := strings.Split(captured, "\n")
	// Keep last N lines (most recent output)
	if len(capLines) > contentLines {
		capLines = capLines[len(capLines)-contentLines:]
	}

	// Build all lines: header, separator, then content
	allLines := make([]string, innerHeight)
	allLines[0] = headerStyled
	allLines[1] = separator
	for i := 0; i < contentLines; i++ {
		if i < len(capLines) {
			allLines[i+2] = padOrTruncate(capLines[i], innerWidth)
		} else {
			allLines[i+2] = strings.Repeat(" ", innerWidth)
		}
	}

	content := strings.Join(allLines, "\n")
	return drawBorder(content, width, innerHeight)
}

func aiLabel(cmd string) string {
	tool, ok := tmux.LookupAITool(cmd)
	if !ok {
		return ""
	}
	label := tool.Icon + " " + tool.Name
	return "  " + lipgloss.NewStyle().Foreground(lipgloss.Color(tool.Color)).Bold(true).Render(label)
}

func shortenPath(path string) string {
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		if strings.HasPrefix(path, home) {
			path = "~" + path[len(home):]
		}
	}
	if len(path) > maxPathDisplay {
		path = "..." + path[len(path)-maxPathDisplay+3:]
	}
	return path
}
