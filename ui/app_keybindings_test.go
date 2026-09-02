package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/lunemis/mux/tmux"
)

func TestConfiguredListBindingsDriveNavigationAndModes(t *testing.T) {
	keyMap := mustKeyMap(t, map[string]map[string][]string{
		"list": {
			"down":   {"s"},
			"up":     {"w"},
			"create": {"c"},
			"rename": {"e"},
			"kill":   {"d"},
			"filter": {"f"},
		},
	})
	m := NewModelWithKeyMap(keyMap)
	m.sessions = []tmux.Session{{Name: "one"}, {Name: "two"}}
	m.applyFilter()

	m = updateModel(t, m, runeKey("s"))
	if m.cursor != 1 {
		t.Fatalf("custom down cursor = %d, want 1", m.cursor)
	}
	m = updateModel(t, m, runeKey("w"))
	if m.cursor != 0 {
		t.Fatalf("custom up cursor = %d, want 0", m.cursor)
	}

	m = updateModel(t, m, runeKey("c"))
	if m.mode != modeCreate {
		t.Fatalf("custom create mode = %v, want modeCreate", m.mode)
	}
	m.mode = modeList

	m = updateModel(t, m, runeKey("e"))
	if m.mode != modeRename {
		t.Fatalf("custom rename mode = %v, want modeRename", m.mode)
	}
	m.mode = modeList

	m = updateModel(t, m, runeKey("d"))
	if m.mode != modeConfirmKill {
		t.Fatalf("custom kill mode = %v, want modeConfirmKill", m.mode)
	}
	m.mode = modeList

	m = updateModel(t, m, runeKey("f"))
	if m.mode != modeFilter {
		t.Fatalf("custom filter mode = %v, want modeFilter", m.mode)
	}
}

func TestReplacedListBindingStopsUsingItsDefault(t *testing.T) {
	keyMap := mustKeyMap(t, map[string]map[string][]string{
		"list": {"create": {"c"}},
	})
	m := NewModelWithKeyMap(keyMap)

	m = updateModel(t, m, runeKey("n"))
	if m.mode != modeList {
		t.Fatalf("replaced default n changed mode to %v", m.mode)
	}
}

func TestConfiguredModalBindingsDriveEveryMode(t *testing.T) {
	keyMap := mustKeyMap(t, map[string]map[string][]string{
		"create": {"switch_field": {"v"}, "cancel": {"b"}},
		"rename": {"cancel": {"b"}},
		"filter": {"clear": {"u"}},
		"kill":   {"cancel": {"n"}},
	})

	t.Run("create switch and cancel", func(t *testing.T) {
		m := NewModelWithKeyMap(keyMap)
		m.mode = modeCreate
		m.createModel = newCreateModel()
		m = updateModel(t, m, runeKey("v"))
		if m.createModel.focused != 1 {
			t.Fatalf("custom switch field focus = %d, want 1", m.createModel.focused)
		}
		m = updateModel(t, m, runeKey("b"))
		if m.mode != modeList {
			t.Fatalf("custom create cancel mode = %v, want modeList", m.mode)
		}
	})

	t.Run("rename cancel", func(t *testing.T) {
		m := NewModelWithKeyMap(keyMap)
		m.mode = modeRename
		m.renameModel = newRenameModel("old")
		m = updateModel(t, m, runeKey("b"))
		if m.mode != modeList {
			t.Fatalf("custom rename cancel mode = %v, want modeList", m.mode)
		}
	})

	t.Run("filter clear", func(t *testing.T) {
		m := NewModelWithKeyMap(keyMap)
		m.mode = modeFilter
		m.filterMod = newFilterModel("needle")
		next, cmd := m.Update(runeKey("u"))
		m = next.(Model)
		if cmd == nil {
			t.Fatal("custom filter clear returned no command")
		}
		m = updateModel(t, m, cmd())
		if m.mode != modeList || m.filterText != "" {
			t.Fatalf("custom filter clear mode/text = %v/%q, want modeList/empty", m.mode, m.filterText)
		}
	})

	t.Run("kill cancel", func(t *testing.T) {
		m := NewModelWithKeyMap(keyMap)
		m.mode = modeConfirmKill
		m.confirmKillMod = newConfirmKillModel("doomed")
		next, cmd := m.Update(runeKey("n"))
		m = next.(Model)
		if cmd == nil {
			t.Fatal("custom kill cancel returned no command")
		}
		m = updateModel(t, m, cmd())
		if m.mode != modeList {
			t.Fatalf("custom kill cancel mode = %v, want modeList", m.mode)
		}
	})
}

func TestConfiguredGlobalQuitWorksFromModalMode(t *testing.T) {
	keyMap := mustKeyMap(t, map[string]map[string][]string{
		"global": {"quit": {"z"}},
	})
	m := NewModelWithKeyMap(keyMap)
	m.mode = modeCreate
	m.createModel = newCreateModel()

	_, cmd := m.Update(runeKey("z"))
	if cmd == nil {
		t.Fatal("custom global quit returned no command")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("custom global quit command returned %T, want tea.QuitMsg", cmd())
	}
}

func TestRenderedHelpAndPromptsUseConfiguredBindings(t *testing.T) {
	keyMap := mustKeyMap(t, map[string]map[string][]string{
		"list":   {"up": {"w"}, "down": {"s"}, "create": {"c"}},
		"create": {"switch_field": {"ctrl+n"}, "submit": {"ctrl+s"}, "cancel": {"ctrl+x"}},
		"rename": {"submit": {"ctrl+s"}, "cancel": {"ctrl+x"}},
		"filter": {"apply": {"ctrl+s"}, "clear": {"ctrl+x"}},
		"kill":   {"confirm": {"enter"}, "cancel": {"esc"}},
	})

	help := renderHelp(keyMap, 120)
	assertContainsAll(t, help, "w/s", "navigate", "c", "new")
	if lines := strings.Split(help, "\n"); len(lines) != 2 {
		t.Fatalf("help lines = %d, want exactly 2", len(lines))
	} else {
		for i, line := range lines {
			if width := ansi.StringWidth(line); width != 120 {
				t.Errorf("help line %d width = %d, want 120", i, width)
			}
		}
	}
	assertContainsAll(t, newCreateModel().View(keyMap), "ctrl+n", "ctrl+s", "ctrl+x")
	assertContainsAll(t, newRenameModel("old").View(keyMap), "ctrl+s", "ctrl+x")
	assertContainsAll(t, newFilterModel("").View(keyMap), "ctrl+s", "ctrl+x")
	assertContainsAll(t, newConfirmKillModel("old").View(keyMap), "enter", "esc")
}

func mustKeyMap(t *testing.T, overrides map[string]map[string][]string) KeyMap {
	t.Helper()
	keyMap, err := NewKeyMap(overrides)
	if err != nil {
		t.Fatalf("NewKeyMap() error = %v", err)
	}
	return keyMap
}

func runeKey(value string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)}
}

func updateModel(t *testing.T, model Model, msg tea.Msg) Model {
	t.Helper()
	next, _ := model.Update(msg)
	got, ok := next.(Model)
	if !ok {
		t.Fatalf("Update() model type = %T, want ui.Model", next)
	}
	return got
}

func assertContainsAll(t *testing.T, got string, want ...string) {
	t.Helper()
	for _, value := range want {
		if !strings.Contains(got, value) {
			t.Errorf("rendered text %q does not contain %q", got, value)
		}
	}
}
