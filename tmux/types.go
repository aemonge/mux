// Package tmux provides functions for managing tmux sessions,
// capturing pane output, and detecting running processes.
package tmux

import "time"

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
