package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleQuestListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m.currentSection = "habits"
		return m, nil

	case "shift+tab":
		m.currentSection = "journeys"
		return m, nil

	case "?":
		m.mode = HelpView
		return m, nil

	case "r":
		return m.refreshData(), nil

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

func (m Model) handleHabitListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m.currentSection = "journeys"
		return m, nil

	case "shift+tab":
		m.currentSection = "quests"
		return m, nil

	case "?":
		m.mode = HelpView
		return m, nil

	case "r":
		return m.refreshData(), nil

	case "n":
		return m.createNewHabit(), nil

	case " ", "enter":
		if item, ok := m.habitList.SelectedItem().(habitItem); ok {
			return m.toggleHabit(item.habit), nil
		}

	case "d":
		if item, ok := m.habitList.SelectedItem().(habitItem); ok {
			return m.showDeleteConfirmHabit(item.habit), nil
		}

	default:
		var cmd tea.Cmd
		m.habitList, cmd = m.habitList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleJourneyListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m.currentSection = "quests"
		return m, nil

	case "shift+tab":
		m.currentSection = "habits"
		return m, nil

	case "?":
		m.mode = HelpView
		return m, nil

	case "r":
		return m.refreshData(), nil

	case "n":
		return m.createNewJourney(), nil

	case " ", "enter":
		if item, ok := m.journeyList.SelectedItem().(journeyItem); ok {
			return m.enterJourney(item.journey), nil
		}

	case "d":
		if item, ok := m.journeyList.SelectedItem().(journeyItem); ok {
			return m.showDeleteConfirmJourney(item.journey), nil
		}

	default:
		var cmd tea.Cmd
		m.journeyList, cmd = m.journeyList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleJourneyDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		m.mode = QuestListView
		m.selectedJourney = nil
		return m, nil

	case "?":
		m.mode = HelpView
		return m, nil

	case "r":
		return m.refreshData(), nil

	case "n":
		return m.createNewQuestInJourney(), nil

	case " ", "enter":
		if item, ok := m.journeyQuestList.SelectedItem().(questItem); ok {
			return m.toggleQuest(item.quest), nil
		}

	case "d":
		if item, ok := m.journeyQuestList.SelectedItem().(questItem); ok {
			return m.showDeleteConfirm(item.quest), nil
		}

	default:
		var cmd tea.Cmd
		m.journeyQuestList, cmd = m.journeyQuestList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleErrorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "r":
		return m.refreshData(), nil
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
			if m.confirmQuest != nil {
				return m.confirmDeleteQuest(), nil
			} else if m.confirmHabit != nil {
				return m.confirmDeleteHabit(), nil
			} else if m.confirmJourney != nil {
				return m.confirmDeleteJourney(), nil
			}
		}
		return m.cancelDelete(), nil
	}
	return m, nil
}
