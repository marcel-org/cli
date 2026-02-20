package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleQuestListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "?":
		m.mode = HelpView
		return m, nil

	case "r":
		return m.refreshQuests(), nil

	case "n":
		return m.createNewQuest(), nil

	case " ", "enter":
		if item, ok := m.questList.SelectedItem().(questItem); ok {
			return m.toggleQuest(item.quest), nil
		}

	case "d":
		if item, ok := m.questList.SelectedItem().(questItem); ok {
			return m.showDeleteConfirm(item.quest), nil
		}

	default:
		var cmd tea.Cmd
		m.questList, cmd = m.questList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleErrorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "r":
		return m.refreshQuests(), nil
	}
	return m, nil
}

func (m Model) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc", "?":
		m.mode = QuestListView
	}
	return m, nil
}

func (m Model) handleConfirmDeleteKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return m.cancelDelete(), nil
	case "left", "h":
		m.confirmSelected = false
	case "right", "l":
		m.confirmSelected = true
	case "enter", " ":
		if m.confirmSelected {
			return m.confirmDeleteQuest(), nil
		}
		return m.cancelDelete(), nil
	}
	return m, nil
}
