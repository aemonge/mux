package tmux

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	listFormat       = "#{session_name}|#{session_windows}|#{session_created}|#{session_attached}|#{pane_current_path}|#{session_last_attached}|#{pane_current_command}|#{pane_pid}"
	originSessionEnv = "MUX_ORIGIN_SESSION"
)

// ListSessions returns sessions in OS-switcher order: previous first and current last.
func ListSessions() ([]Session, error) {
	out, err := runner.Output("tmux", "list-sessions", "-F", listFormat)
	if err != nil {
		// tmux returns error when no server is running
		if strings.Contains(err.Error(), "exit status") {
			return nil, nil
		}
		return nil, fmt.Errorf("list sessions: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil, nil
	}

	sessions := make([]Session, 0, len(lines))
	for _, line := range lines {
		s, err := parseLine(line)
		if err != nil {
			continue
		}
		sessions = append(sessions, s)
	}

	sortSessionsForSwitcher(sessions, currentSessionName())

	return sessions, nil
}

func parseLine(line string) (Session, error) {
	parts := strings.SplitN(line, "|", 8)
	if len(parts) < 8 {
		return Session{}, fmt.Errorf("unexpected format: %s", line)
	}

	windows, _ := strconv.Atoi(parts[1])
	createdUnix, _ := strconv.ParseInt(parts[2], 10, 64)
	attached, _ := strconv.Atoi(parts[3])
	lastAttachedUnix, _ := strconv.ParseInt(parts[5], 10, 64)
	panePID, _ := strconv.Atoi(parts[7])

	var lastAttached time.Time
	if lastAttachedUnix > 0 {
		lastAttached = time.Unix(lastAttachedUnix, 0)
	}

	activeCommand := resolveCommand(panePID, parts[6])
	gitInfo := LookupGitInfo(parts[4])

	return Session{
		Name:          parts[0],
		WindowCount:   windows,
		Created:       time.Unix(createdUnix, 0),
		LastAttached:  lastAttached,
		Attached:      attached > 0,
		Directory:     parts[4],
		ActiveCommand: activeCommand,
		PanePID:       panePID,
		GitBranch:     gitInfo.Branch,
		IsWorktree:    gitInfo.IsWorktree,
	}, nil
}

func currentSessionName() string {
	if origin := os.Getenv(originSessionEnv); origin != "" {
		return origin
	}
	pane := os.Getenv("TMUX_PANE")
	if pane == "" {
		return ""
	}
	out, err := runner.Output("tmux", "display-message", "-p", "-t", pane, "#{session_name}")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func sortSessionsForSwitcher(sessions []Session, current string) {
	sort.Slice(sessions, func(i, j int) bool {
		if !sessions[i].LastAttached.Equal(sessions[j].LastAttached) {
			return sessions[i].LastAttached.After(sessions[j].LastAttached)
		}
		if !sessions[i].Created.Equal(sessions[j].Created) {
			return sessions[i].Created.After(sessions[j].Created)
		}
		return sessions[i].Name < sessions[j].Name
	})

	if current == "" {
		return
	}
	for i := range sessions {
		if sessions[i].Name != current {
			continue
		}
		active := sessions[i]
		copy(sessions[i:], sessions[i+1:])
		sessions[len(sessions)-1] = active
		return
	}
}

// CreateSession creates a new detached tmux session with the given name.
func CreateSession(name string) error {
	return runner.Run("tmux", "new-session", "-d", "-s", name)
}

// CreateSessionWithDir creates a new detached tmux session starting in the given directory.
func CreateSessionWithDir(name, dir string) error {
	return runner.Run("tmux", "new-session", "-d", "-s", name, "-c", dir)
}

// KillSession destroys the tmux session with the given name.
func KillSession(name string) error {
	return runner.Run("tmux", "kill-session", "-t", name)
}

// RenameSession renames a tmux session from oldName to newName.
func RenameSession(oldName, newName string) error {
	return runner.Run("tmux", "rename-session", "-t", oldName, newName)
}
