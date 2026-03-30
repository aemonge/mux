package tmux

import "testing"

func TestIsAICommand(t *testing.T) {
	tests := []struct {
		cmd  string
		want bool
	}{
		{"claude", true},
		{"codex", true},
		{"aider", true},
		{"gemini", true},
		{"bash", false},
		{"vim", false},
		{"", false},
		{"Claude", false}, // case-sensitive
	}

	for _, tt := range tests {
		if got := IsAICommand(tt.cmd); got != tt.want {
			t.Errorf("IsAICommand(%q) = %v, want %v", tt.cmd, got, tt.want)
		}
	}
}
