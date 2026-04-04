package tmux

import "strings"

const statusScanLines = 15

// Permission patterns — checked first (highest priority).
var permissionPatterns = []string{
	"Allow",
	"Deny",
	"(Y)es",
	"(y/N)",
	"(Y/n)",
	"approve",
	"Do you want to",
	"permission",
}

// Thinking/working patterns.
var thinkingPatterns = []string{
	// Braille spinner characters used by many CLI tools
	"\u280b", "\u2819", "\u2839", "\u2838", "\u283c", "\u2834", "\u2826", "\u2827", "\u2807", "\u280f",
	"Thinking",
	"Reading",
	"Writing",
	"Searching",
	"Analyzing",
	"Generating",
	"Planning",
	"Editing",
	"Running",
}

// DetectAgentStatus infers the state of an AI agent from captured pane output.
// Returns StatusUnknown for non-AI commands.
func DetectAgentStatus(activeCommand, capturedPane string) AgentStatus {
	if !IsAICommand(activeCommand) {
		return StatusUnknown
	}

	lines := lastNLines(capturedPane, statusScanLines)

	for _, line := range lines {
		for _, p := range permissionPatterns {
			if strings.Contains(line, p) {
				return StatusPermission
			}
		}
	}

	for _, line := range lines {
		for _, p := range thinkingPatterns {
			if strings.Contains(line, p) {
				return StatusThinking
			}
		}
	}

	return StatusIdle
}

// StatusIcon returns a short icon string for the given status.
func StatusIcon(s AgentStatus) string {
	switch s {
	case StatusThinking:
		return "\u27f3"
	case StatusPermission:
		return "\u26a0"
	default:
		return ""
	}
}

// StatusLabel returns a human-readable label for the given status.
func StatusLabel(s AgentStatus) string {
	switch s {
	case StatusThinking:
		return "thinking"
	case StatusPermission:
		return "permission"
	case StatusIdle:
		return "idle"
	default:
		return ""
	}
}

func lastNLines(s string, n int) []string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines
}
