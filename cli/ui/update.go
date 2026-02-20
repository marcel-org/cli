package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.needsRedraw {
		m.needsRedraw = false
		cmds = append(cmds, tea.ClearScreen)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.questList = newQuestList(m.data, m.width-4, m.height-8)
			m.ready = true
		} else {
			m.questList.SetSize(m.width-4, m.height-8)
		}

	case tea.KeyMsg:
		switch m.mode {
		case QuestListView:
			return m.handleQuestListKeys(msg)
		case ErrorView:
			return m.handleErrorKeys(msg)
		case HelpView:
			return m.handleHelpKeys(msg)
		case ConfirmDeleteView:
			return m.handleConfirmDeleteKeys(msg)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
