package ui

import (
	"github.com/charmbracelet/huh"
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
		m.currentSection = "calendar"
		return m, nil

	case "?":
		m.mode = HelpView
		return m, nil

	case "r":
		return m.refreshData(), nil

	case "n":
		return m.createNewQuest()

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
		return m.createNewHabit()

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
		m.currentSection = "calendar"
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
		return m.createNewJourney()

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

func (m Model) handleCalendarKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m.currentSection = "quests"
		return m, nil

	case "shift+tab":
		m.currentSection = "journeys"
		return m, nil

	case "?":
		m.mode = HelpView
		return m, nil

	case "r":
		return m.refreshData(), nil

	case "left", "h":
		m.calendar.NavigateLeft()
		return m, nil

	case "right", "l":
		m.calendar.NavigateRight()
		return m, nil

	case "up", "k":
		m.calendar.NavigateUp()
		return m, nil

	case "down", "j":
		m.calendar.NavigateDown()
		return m, nil

	case "ctrl+left":
		m.calendar.NavigatePrevMonth()
		return m, nil

	case "ctrl+right":
		m.calendar.NavigateNextMonth()
		return m, nil

	case "t":
		m.calendar.GoToToday()
		return m, nil
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
		return m.createNewQuestInJourney()

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

func (m Model) handleFormUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var form *huh.Form
	var returnMode ViewMode = QuestListView

	switch m.mode {
	case QuestFormView:
		form = m.questForm
		returnMode = QuestListView
	case JourneyFormView:
		form = m.journeyForm
		if m.selectedJourney != nil {
			returnMode = JourneyDetailView
		} else {
			returnMode = QuestListView
		}
	case HabitFormView:
		form = m.habitForm
		returnMode = QuestListView
	}

	if form == nil {
		m.mode = returnMode
		return m, nil
	}

	newForm, cmd := form.Update(msg)
	if f, ok := newForm.(*huh.Form); ok {
		switch m.mode {
		case QuestFormView:
			m.questForm = f
		case JourneyFormView:
			m.journeyForm = f
		case HabitFormView:
			m.habitForm = f
		}

		if f.State == huh.StateCompleted {
			return m.handleFormCompletion()
		}

		if f.State == huh.StateAborted {
			m.mode = returnMode
			return m, nil
		}
	}

	return m, cmd
}
