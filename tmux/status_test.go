package tmux

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestReadAgentStatus(t *testing.T) {
	sessionID := fmt.Sprintf("test-%d", time.Now().UnixNano())
	path := fmt.Sprintf("/tmp/mux-status-%s.json", sessionID)
	defer os.Remove(path)

	tests := []struct {
		name   string
		data   hookStatus
		want   AgentStatus
	}{
		{"thinking", hookStatus{Status: "thinking", TS: time.Now().Unix()}, StatusThinking},
		{"permission", hookStatus{Status: "permission", TS: time.Now().Unix()}, StatusPermission},
		{"idle", hookStatus{Status: "idle", TS: time.Now().Unix()}, StatusIdle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := json.Marshal(tt.data)
			os.WriteFile(path, data, 0644)

			got := ReadAgentStatus(sessionID)
			if got != tt.want {
				t.Errorf("ReadAgentStatus() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestReadAgentStatusMissingFile(t *testing.T) {
	got := ReadAgentStatus("nonexistent-session-id")
	if got != StatusUnknown {
		t.Errorf("ReadAgentStatus(missing) = %d, want StatusUnknown", got)
	}
}

func TestReadAgentStatusEmptySessionID(t *testing.T) {
	got := ReadAgentStatus("")
	if got != StatusUnknown {
		t.Errorf("ReadAgentStatus('') = %d, want StatusUnknown", got)
	}
}

func TestReadAgentStatusStale(t *testing.T) {
	sessionID := fmt.Sprintf("test-stale-%d", time.Now().UnixNano())
	path := fmt.Sprintf("/tmp/mux-status-%s.json", sessionID)
	defer os.Remove(path)

	stale := hookStatus{Status: "thinking", TS: time.Now().Add(-10 * time.Minute).Unix()}
	data, _ := json.Marshal(stale)
	os.WriteFile(path, data, 0644)

	got := ReadAgentStatus(sessionID)
	if got != StatusUnknown {
		t.Errorf("ReadAgentStatus(stale) = %d, want StatusUnknown", got)
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
