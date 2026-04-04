// Package tmux provides functions for managing tmux sessions,
// capturing pane output, and detecting running processes.
package tmux

import "time"

// AgentStatus represents the detected state of an AI coding agent.
type AgentStatus int

const (
	StatusUnknown    AgentStatus = iota // not an AI command or unable to detect
	StatusIdle                          // AI agent is waiting for user input
	StatusThinking                      // AI agent is actively working
	StatusPermission                    // AI agent is waiting for permission approval
)

// Session represents a tmux session with its metadata and state.
type Session struct {
	Name          string
	Windows       int
	Created       time.Time
	Activity      time.Time
	Attached      bool
	Directory     string
	ActiveCommand string
	PanePID       int
}
