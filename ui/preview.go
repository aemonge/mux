package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lunemis/mux/tmux"
)

func renderPreview(session *tmux.Session, captured string, width, height int, status tmux.AgentStatus, tokenUsage *tmux.TokenUsage) string {
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

	// Header: build plain text first, then append styled badges after padding
	// to prevent ANSI codes and ambiguous-width icons from being clipped.
	badge := aiLabelPlain(session.ActiveCommand)
	statusBadge := statusLabelPlain(status)
	suffixPlain := badge.text + statusBadge.text
	// Each ambiguous-width icon takes 2 cells but ansi.StringWidth reports 1
	suffixExtra := badge.extraWidth + statusBadge.extraWidth

	headerText := fmt.Sprintf("[ %s ]  %s", session.Name, shortenPath(session.Directory))
	headerWidth := innerWidth - len(suffixPlain) - suffixExtra
	if headerWidth < 0 {
		headerWidth = 0
	}
	headerPadded := padOrTruncate(headerText, headerWidth)

	// Now build the styled suffix
	styledSuffix := badge.styled + statusBadge.styled
	headerStyled := lipgloss.NewStyle().Foreground(colorAccent).Bold(true).Render(headerPadded) + styledSuffix
	separator := lipgloss.NewStyle().Foreground(colorBorder).Render(
		strings.Repeat("─", innerWidth))

	// Token usage line (optional)
	var tokenLine string
	if tokenUsage != nil {
		tokenLine = formatTokenLine(tokenUsage, innerWidth)
	}

	// Available lines for content (minus header + separator + optional token line)
	headerLines := 2
	if tokenLine != "" {
		headerLines = 3
	}
	contentLines := innerHeight - headerLines
	if contentLines < 1 {
		contentLines = 1
	}

	capLines := strings.Split(captured, "\n")
	// Keep last N lines (most recent output)
	if len(capLines) > contentLines {
		capLines = capLines[len(capLines)-contentLines:]
	}

	// Build all lines: header, [token], separator, then content
	allLines := make([]string, innerHeight)
	lineIdx := 0
	allLines[lineIdx] = headerStyled
	lineIdx++
	if tokenLine != "" {
		allLines[lineIdx] = tokenLine
		lineIdx++
	}
	allLines[lineIdx] = separator
	lineIdx++
	for i := 0; i < contentLines; i++ {
		if i < len(capLines) {
			allLines[lineIdx+i] = padOrTruncate(capLines[i], innerWidth)
		} else {
			allLines[lineIdx+i] = strings.Repeat(" ", innerWidth)
		}
	}

	content := strings.Join(allLines, "\n")
	return drawBorder(content, width, innerHeight)
}

// labelInfo holds both the styled and plain-text versions of a badge,
// plus extra width to compensate for ambiguous-width Unicode characters
// that terminals render as 2 cells but ansi.StringWidth measures as 1.
type labelInfo struct {
	text       string // plain text for width calculation (e.g. "  ✦ claude")
	styled     string // ANSI-styled version for display
	extraWidth int    // extra cells for ambiguous-width chars (1 per icon)
}

func aiLabelPlain(cmd string) labelInfo {
	tool, ok := tmux.LookupAITool(cmd)
	if !ok {
		return labelInfo{}
	}
	text := "  " + tool.Icon + " " + tool.Name
	styled := "  " + lipgloss.NewStyle().Foreground(lipgloss.Color(tool.Color)).Bold(true).Render(tool.Icon+" "+tool.Name)
	return labelInfo{text: text, styled: styled, extraWidth: 1}
}

func statusLabelPlain(s tmux.AgentStatus) labelInfo {
	icon := tmux.StatusIcon(s)
	if icon == "" {
		return labelInfo{}
	}
	label := tmux.StatusLabel(s)
	var color lipgloss.Color
	switch s {
	case tmux.StatusThinking:
		color = lipgloss.Color("#FBBF24") // yellow
	case tmux.StatusPermission:
		color = lipgloss.Color("#EF4444") // red
	default:
		return labelInfo{}
	}
	text := "  " + icon + " " + label
	styled := "  " + lipgloss.NewStyle().Foreground(color).Bold(true).Render(icon+" "+label)
	return labelInfo{text: text, styled: styled, extraWidth: 1}
}

func formatTokenLine(u *tmux.TokenUsage, width int) string {
	text := fmt.Sprintf("  %s in / %s out  ~$%.2f",
		tmux.FormatTokens(u.InputTokens),
		tmux.FormatTokens(u.OutputTokens),
		u.TotalCost)
	styled := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
		padOrTruncate(text, width))
	return styled
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
