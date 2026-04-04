package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/lunemis/mux/tmux"
)

func renderSessionList(sessions []tmux.Session, cursor int, filter string, width, height int) string {
	innerWidth := width - 2 // border chars
	innerHeight := height - 2

	if len(sessions) == 0 {
		msg := "No tmux sessions found"
		if filter != "" {
			msg = fmt.Sprintf("No match: \"%s\"", filter)
		}
		// Center the message
		lines := make([]string, innerHeight)
		mid := innerHeight / 2
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

	// Calculate scroll offset to keep cursor visible
	offset := 0
	if cursor >= innerHeight {
		offset = cursor - innerHeight + 1
	}

	lines := make([]string, innerHeight)
	for i := 0; i < innerHeight; i++ {
		idx := i + offset
		if idx < len(sessions) {
			lines[i] = formatSessionRow(sessions[idx], idx == cursor, innerWidth)
		} else {
			lines[i] = strings.Repeat(" ", innerWidth)
		}
	}

	content := strings.Join(lines, "\n")
	return drawBorder(content, width, innerHeight)
}

func formatSessionRow(s tmux.Session, selected bool, width int) string {
	status := "○"
	if s.Attached {
		status = "*"
	}

	name := s.Name
	if len(name) > maxSessionNameDisplay {
		name = name[:maxSessionNameDisplay-3] + "..."
	}

	ago := timeAgo(s.Created)

	// AI command icon — plain text, styled later to avoid ANSI truncation issues.
	// Ambiguous-width icons (✦ etc.) render as 2 cells in most terminals
	// but ansi.StringWidth reports 1, so we reserve extra space.
	icon, iconColor := commandIconPlain(s.ActiveCommand)
	iconReserved := 3 // icon(2 cells) + space(1), or 3 spaces for non-AI

	textWidth := width - iconReserved
	if textWidth < 0 {
		textWidth = 0
	}
	text := fmt.Sprintf(" %s %-18s %s", status, name, ago)
	paddedText := padOrTruncate(text, textWidth)

	// Build styled icon separately
	var styledIcon string
	if iconColor != "" {
		styledIcon = lipgloss.NewStyle().Foreground(lipgloss.Color(iconColor)).Render(icon) + " "
	} else {
		styledIcon = "   "
	}

	row := paddedText + styledIcon

	if selected {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCursor).
			Background(colorSelected).
			Render(row)
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Render(row)
}

// commandIconPlain returns the raw icon and its color for known AI CLIs.
// Returns empty strings for non-AI commands.
func commandIconPlain(cmd string) (icon string, color string) {
	if tool, ok := tmux.LookupAITool(cmd); ok {
		return tool.Icon, tool.Color
	}
	return "", ""
}

func centerText(s string, width int) string {
	pad := (width - len(s)) / 2
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat(" ", pad) + s
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%3ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%3dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%3dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%3dd", int(d.Hours()/24))
	}
}
