package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/lunemis/mux/theme"
	"github.com/lunemis/mux/tmux"
	"github.com/lunemis/mux/ui"
)

var version = "dev"

func main() {
	themeName := os.Getenv("MUX_THEME")
	if themeName == "" {
		themeName = "default"
	}

	rootCmd := &cobra.Command{
		Use:     "mux",
		Short:   "TUI tmux session manager",
		Version: version,
		RunE:    runTUI,
		// Suppress cobra's default completion and help subcommands
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	rootCmd.SetVersionTemplate("mux {{.Version}}\n")
	rootCmd.PersistentFlags().StringVar(&themeName, "theme", themeName,
		fmt.Sprintf("Color theme (%s)", joinWith(theme.Names(), ", ")))

	popupCmd := &cobra.Command{
		Use:   "popup",
		Short: "Open mux as a tmux popup overlay",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := theme.Get(themeName); err != nil {
				return err
			}
			return tmux.OpenPopup("--theme", themeName)
		},
	}

	setupKeybindCmd := &cobra.Command{
		Use:   "setup-keybind [key]",
		Short: fmt.Sprintf("Add popup keybinding to tmux config (default: %s)", tmux.DefaultBindKey),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := tmux.DefaultBindKey
			if len(args) > 0 {
				key = args[0]
			}
			return tmux.SetupKeybind(key)
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show AI session summary for tmux statusbar",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus()
		},
	}

	rootCmd.AddCommand(popupCmd, setupKeybindCmd, statusCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runStatus() error {
	sessions, err := tmux.ListSessions()
	if err != nil {
		return err
	}

	var parts []string
	for _, s := range sessions {
		tool, ok := tmux.LookupAITool(s.ActiveCommand)
		if !ok {
			continue
		}
		parts = append(parts, tool.Icon)
	}

	if len(parts) == 0 {
		return nil // no AI sessions, output nothing
	}

	fmt.Print(fmt.Sprintf(" %s ", joinWith(parts, " ")))
	return nil
}

func joinWith(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}

func configureTheme(name string) error {
	selected, err := theme.Get(name)
	if err != nil {
		return err
	}
	ui.UseTheme(selected)
	return nil
}

func runTUI(cmd *cobra.Command, args []string) error {
	name, err := cmd.Flags().GetString("theme")
	if err != nil {
		return err
	}
	if err := configureTheme(name); err != nil {
		return err
	}

	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())

	result, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	if m, ok := result.(ui.Model); ok {
		if name := m.AttachName(); name != "" {
			if err := ui.AttachToSession(name, m.AttachWindowIndex(), m.AttachPaneIndex()); err != nil {
				return fmt.Errorf("failed to attach: %w", err)
			}
		}
	}
	return nil
}
