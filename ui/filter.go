package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type filterModel struct {
	input textinput.Model
}

type filterAppliedMsg struct {
	text    string
	cleared bool
}

func newFilterModel(currentText string) filterModel {
	fi := textinput.New()
	fi.Placeholder = "filter..."
	fi.CharLimit = filterCharLimit
	fi.Width = filterInputWidth
	fi.SetValue(currentText)
	fi.Focus()

	return filterModel{input: fi}
}

func (m filterModel) Update(msg tea.Msg) (filterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return filterAppliedMsg{text: "", cleared: true}
			}
		case "enter":
			return m, func() tea.Msg {
				return filterAppliedMsg{text: m.input.Value(), cleared: false}
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m filterModel) View() string {
	return inputLabelStyle.Render("/ ") + m.input.View()
}

func (m filterModel) LiveText() string {
	return m.input.Value()
}
