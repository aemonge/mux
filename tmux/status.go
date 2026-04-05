package tmux

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	// statusStaleThreshold is how old a status file can be before we ignore it.
	statusStaleThreshold = 5 * time.Minute
)

// hookStatus represents the JSON written by the mux-status.sh hook script.
type hookStatus struct {
	Status string `json:"status"`
	Tool   string `json:"tool,omitempty"`
	TS     int64  `json:"ts"`
}

// ReadAgentStatus reads the agent status from the hook-written file.
// Falls back to StatusUnknown if no file exists or it's stale.
func ReadAgentStatus(sessionID string) AgentStatus {
	if sessionID == "" {
		return StatusUnknown
	}

	path := fmt.Sprintf("/tmp/mux-status-%s.json", sessionID)
	data, err := os.ReadFile(path)
	if err != nil {
		return StatusUnknown
	}

	var hs hookStatus
	if err := json.Unmarshal(data, &hs); err != nil {
		return StatusUnknown
	}

	// Ignore stale status
	if time.Since(time.Unix(hs.TS, 0)) > statusStaleThreshold {
		return StatusUnknown
	}

	switch hs.Status {
	case "thinking":
		return StatusThinking
	case "permission":
		return StatusPermission
	case "idle":
		return StatusIdle
	default:
		return StatusUnknown
	}
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
