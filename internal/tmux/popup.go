package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// OpenPopup opens mux inside a tmux display-popup overlay.
// Must be called from inside a tmux session.
func OpenPopup() error {
	if os.Getenv("TMUX") == "" {
		return fmt.Errorf("mux popup must be run inside a tmux session")
	}

	version, err := getTmuxVersion()
	if err != nil {
		return fmt.Errorf("failed to detect tmux version: %w", err)
	}
	if version < 3.2 {
		return fmt.Errorf("tmux 3.2+ required for popup (current: %.1f)", version)
	}

	muxPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find mux executable: %w", err)
	}

	cmd := exec.Command("tmux", "display-popup",
		"-E",
		"-w", "85%",
		"-h", "80%",
		muxPath,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SetupKeybind adds a popup keybinding to ~/.tmux.conf.
func SetupKeybind(key string) error {
	muxPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find mux executable: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to find home directory: %w", err)
	}

	confPath := filepath.Join(home, ".tmux.conf")
	bindLine := fmt.Sprintf(`bind %s display-popup -E -w 85%% -h 80%% "%s"`, key, muxPath)
	marker := "# mux popup keybinding"

	// Read existing config
	content, err := os.ReadFile(confPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %w", confPath, err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	replaced := false

	for _, line := range lines {
		if strings.Contains(line, marker) {
			// Replace existing mux keybinding
			newLines = append(newLines, bindLine+"  "+marker)
			replaced = true
			continue
		}
		newLines = append(newLines, line)
	}

	if !replaced {
		// Ensure trailing newline before appending
		if len(newLines) > 0 && newLines[len(newLines)-1] != "" {
			newLines = append(newLines, "")
		}
		newLines = append(newLines, bindLine+"  "+marker)
	}

	result := strings.Join(newLines, "\n")
	// Ensure file ends with newline
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	if err := os.WriteFile(confPath, []byte(result), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", confPath, err)
	}

	fmt.Printf("Added to %s:\n  %s\n\n", confPath, bindLine)
	fmt.Printf("Reload tmux config:\n  tmux source-file %s\n\n", confPath)
	fmt.Printf("Then press: prefix + %s (default prefix: Ctrl+b)\n", key)

	return nil
}

func getTmuxVersion() (float64, error) {
	out, err := exec.Command("tmux", "-V").Output()
	if err != nil {
		return 0, err
	}
	// Output: "tmux 3.4" or "tmux 3.2a"
	s := strings.TrimSpace(string(out))
	s = strings.TrimPrefix(s, "tmux ")
	// Strip trailing letter (e.g. "3.2a" -> "3.2")
	var version float64
	fmt.Sscanf(s, "%f", &version)
	return version, nil
}
