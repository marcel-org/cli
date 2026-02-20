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
			return m.deleteQuest(item.quest), nil
		}
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
