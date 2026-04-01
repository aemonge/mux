package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/lunemis/mux/internal/tmux"
	"github.com/lunemis/mux/internal/ui"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "mux",
		Short:   "TUI tmux session manager",
		Version: version,
		RunE:    runTUI,
		// Suppress cobra's default completion and help subcommands
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	rootCmd.SetVersionTemplate("mux {{.Version}}\n")

	popupCmd := &cobra.Command{
		Use:   "popup",
		Short: "Open mux as a tmux popup overlay",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tmux.OpenPopup()
		},
	}

	setupKeybindCmd := &cobra.Command{
		Use:   "setup-keybind [key]",
		Short: fmt.Sprintf("Add popup keybinding to ~/.tmux.conf (default: %s)", tmux.DefaultBindKey),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := tmux.DefaultBindKey
			if len(args) > 0 {
				key = args[0]
			}
			return tmux.SetupKeybind(key)
		},
	}

	rootCmd.AddCommand(popupCmd, setupKeybindCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runTUI(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())

	result, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	if m, ok := result.(ui.Model); ok {
		if name := m.AttachName(); name != "" {
			if err := ui.AttachToSession(name); err != nil {
				return fmt.Errorf("failed to attach: %w", err)
			}
		}
	}
	return nil
}
