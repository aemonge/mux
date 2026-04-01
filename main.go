package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunemis/mux/internal/tmux"
	"github.com/lunemis/mux/internal/ui"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Println("mux " + version)
			return
		case "popup":
			if err := tmux.OpenPopup(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "setup-keybind":
			key := tmux.DefaultBindKey
			if len(os.Args) > 2 {
				key = os.Args[2]
			}
			if err := tmux.SetupKeybind(key); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		case "--help", "-h":
			printUsage()
			return
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
			printUsage()
			os.Exit(1)
		}
	}

	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())

	result, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// If user selected a session to attach, exec into tmux
	if m, ok := result.(ui.Model); ok {
		if name := m.AttachName(); name != "" {
			if err := ui.AttachToSession(name); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to attach: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func printUsage() {
	fmt.Println(`mux - TUI tmux session manager

Usage:
  mux                         Launch session manager
  mux popup                   Open mux as a tmux popup overlay
  mux setup-keybind [key]     Add popup keybinding to ~/.tmux.conf (default: m)
  mux --version               Show version
  mux --help                  Show this help

Popup Mode:
  Run "mux popup" inside tmux to open mux as a floating overlay.
  Or set up a keybinding for instant access:

    mux setup-keybind          # binds prefix + m
    mux setup-keybind Space    # binds prefix + Space`)
}
