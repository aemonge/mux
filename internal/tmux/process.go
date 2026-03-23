package tmux

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsAICommand reports whether cmd is a known AI CLI process.
func IsAICommand(cmd string) bool {
	switch cmd {
	case "claude", "codex", "aider", "gemini":
		return true
	}
	return false
}

// interpreters are process names that wrap scripts (e.g. node runs codex).
var interpreters = map[string]bool{
	"node": true, "python": true, "python3": true,
}

// resolveCommand returns the logical command name for a pane.
// If rawCmd is an interpreter, it inspects child processes to find the
// actual AI CLI being run (e.g. "node" → "codex").
func resolveCommand(panePID int, rawCmd string) string {
	if panePID <= 0 {
		return rawCmd
	}
	if !interpreters[rawCmd] {
		return rawCmd
	}
	// Look for child processes of the pane's shell that use this interpreter
	children, _ := filepath.Glob(fmt.Sprintf("/proc/[0-9]*/stat"))
	for _, statPath := range children {
		data, err := os.ReadFile(statPath)
		if err != nil {
			continue
		}
		// stat format: pid (comm) state ppid ...
		fields := strings.Fields(string(data))
		if len(fields) < 4 {
			continue
		}
		ppid := fields[3]
		if ppid != fmt.Sprintf("%d", panePID) {
			continue
		}
		// Read cmdline of this child
		pid := fields[0]
		cmdline, err := os.ReadFile("/proc/" + pid + "/cmdline")
		if err != nil {
			continue
		}
		args := strings.Split(string(cmdline), "\x00")
		// args[0] is the interpreter, args[1] is the script path
		for _, arg := range args[1:] {
			base := filepath.Base(arg)
			if IsAICommand(base) {
				return base
			}
		}
	}
	return rawCmd
}
