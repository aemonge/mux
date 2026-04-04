package tmux

import (
	"strings"
	"testing"
)

func TestDetectAgentStatus(t *testing.T) {
	tests := []struct {
		name    string
		command string
		pane    string
		want    AgentStatus
	}{
		{
			name:    "non-AI command returns unknown",
			command: "bash",
			pane:    "$ ",
			want:    StatusUnknown,
		},
		{
			name:    "empty pane with AI command returns idle",
			command: "claude",
			pane:    "",
			want:    StatusIdle,
		},
		{
			name:    "claude permission prompt with Allow",
			command: "claude",
			pane:    "some output\nAllow claude to edit file.go?",
			want:    StatusPermission,
		},
		{
			name:    "claude permission prompt with Y/n",
			command: "claude",
			pane:    "Do you want to proceed? (Y/n)",
			want:    StatusPermission,
		},
		{
			name:    "claude thinking with spinner",
			command: "claude",
			pane:    "⠹ Processing request...",
			want:    StatusThinking,
		},
		{
			name:    "claude thinking with keyword",
			command: "claude",
			pane:    "Thinking about the best approach...",
			want:    StatusThinking,
		},
		{
			name:    "claude reading",
			command: "claude",
			pane:    "Reading file src/main.go",
			want:    StatusThinking,
		},
		{
			name:    "claude idle with shell prompt",
			command: "claude",
			pane:    "Done. Changes applied.\n> ",
			want:    StatusIdle,
		},
		{
			name:    "codex thinking",
			command: "codex",
			pane:    "Searching codebase for references...",
			want:    StatusThinking,
		},
		{
			name:    "permission takes priority over thinking",
			command: "claude",
			pane:    "Writing changes...\nAllow edit? (Y)es/(N)o",
			want:    StatusPermission,
		},
		{
			name:    "only scans last N lines",
			command: "claude",
			pane:    "Allow edit?\n" + strings.Repeat("normal output\n", 20) + "idle prompt",
			want:    StatusIdle,
		},
		{
			name:    "aider detected",
			command: "aider",
			pane:    "Editing files...",
			want:    StatusThinking,
		},
		{
			name:    "gemini detected",
			command: "gemini",
			pane:    "Generating response...",
			want:    StatusThinking,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectAgentStatus(tt.command, tt.pane)
			if got != tt.want {
				t.Errorf("DetectAgentStatus(%q, ...) = %d, want %d", tt.command, got, tt.want)
			}
		})
	}
}

func TestStatusIcon(t *testing.T) {
	tests := []struct {
		status AgentStatus
		want   string
	}{
		{StatusUnknown, ""},
		{StatusIdle, ""},
		{StatusThinking, "⟳"},
		{StatusPermission, "⚠"},
	}
	for _, tt := range tests {
		got := StatusIcon(tt.status)
		if got != tt.want {
			t.Errorf("StatusIcon(%d) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestStatusLabel(t *testing.T) {
	tests := []struct {
		status AgentStatus
		want   string
	}{
		{StatusUnknown, ""},
		{StatusIdle, "idle"},
		{StatusThinking, "thinking"},
		{StatusPermission, "permission"},
	}
	for _, tt := range tests {
		got := StatusLabel(tt.status)
		if got != tt.want {
			t.Errorf("StatusLabel(%d) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestLastNLines(t *testing.T) {
	input := "a\nb\nc\nd\ne"
	got := lastNLines(input, 3)
	if len(got) != 3 || got[0] != "c" || got[1] != "d" || got[2] != "e" {
		t.Errorf("lastNLines(%q, 3) = %v, want [c d e]", input, got)
	}

	short := "a\nb"
	got2 := lastNLines(short, 5)
	if len(got2) != 2 {
		t.Errorf("lastNLines(%q, 5) = %v, want 2 lines", short, got2)
	}
}
