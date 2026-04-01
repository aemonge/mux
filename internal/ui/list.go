package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/lunemis/mux/internal/tmux"
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

	lines := make([]string, innerHeight)
	for i := 0; i < innerHeight; i++ {
		if i < len(sessions) {
			lines[i] = formatSessionRow(sessions[i], i == cursor, innerWidth)
		} else {
			lines[i] = strings.Repeat(" ", innerWidth)
		}
	}

	content := strings.Join(lines, "\n")
	return drawBorder(content, width, innerHeight)
}

func formatSessionRow(s tmux.Session, selected bool, width int) string {
	var status string
	if s.Attached {
		status = "●"
	} else {
		status = "○"
	}

	name := s.Name
	if len(name) > maxSessionNameDisplay {
		name = name[:maxSessionNameDisplay-3] + "..."
	}

	ago := timeAgo(s.Created)

	// AI command icon
	icon := commandIcon(s.ActiveCommand)

	raw := fmt.Sprintf(" %s %-18s %s %s%dw", status, name, ago, icon, s.Windows)

	if selected {
		styled := lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCursor).
			Background(colorSelected).
			Render(padOrTruncate(raw, width))
		return styled
	}

	styled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Render(padOrTruncate(raw, width))
	return styled
}

// commandIcon returns a short icon string (with trailing space) for known AI CLIs,
// or two spaces for anything else, keeping column widths consistent.
func commandIcon(cmd string) string {
	switch cmd {
	case "claude":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B")).Render("✦") + " "
	case "codex":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#60A5FA")).Render("◈") + " "
	case "aider":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#34D399")).Render("⬡") + " "
	case "gemini":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA")).Render("✧") + " "
	default:
		return "  "
	}
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
