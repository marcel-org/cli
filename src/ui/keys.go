package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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
			return m.toggleQuest(item.quest)
		}

	case "d":
		if item, ok := m.questList.SelectedItem().(questItem); ok {
			return m.showDeleteConfirm(item.quest), nil
		}

	case "e":
		if item, ok := m.questList.SelectedItem().(questItem); ok {
			return m.editQuest(item.quest)
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
			return m.toggleHabit(item.habit)
		}

	case "d":
		if item, ok := m.habitList.SelectedItem().(habitItem); ok {
			return m.showDeleteConfirmHabit(item.habit), nil
		}

	case "e":
		if item, ok := m.habitList.SelectedItem().(habitItem); ok {
			return m.editHabit(item.habit)
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
		m.currentSection = "spaces"
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

	case "e":
		if item, ok := m.journeyList.SelectedItem().(journeyItem); ok {
			return m.editJourney(item.journey)
		}

	default:
		var cmd tea.Cmd
		m.journeyList, cmd = m.journeyList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleSpaceListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m.currentSection = "calendar"
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
		return m.createNewSpace()

	case " ", "enter":
		if item, ok := m.spaceList.SelectedItem().(spaceItem); ok {
			return m.enterSpace(item.space), nil
		}

	case "e":
		if item, ok := m.spaceList.SelectedItem().(spaceItem); ok {
			return m.editSpace(item.space)
		}

	default:
		var cmd tea.Cmd
		m.spaceList, cmd = m.spaceList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleCalendarKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.calendar.IsFocusedOnEventList() {
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			m.calendar.FocusMonthView()
			return m, nil

		case "up", "k":
			m.calendar.NavigateEventListUp()
			return m, nil

		case "down", "j":
			m.calendar.NavigateEventListDown()
			return m, nil

		case "d":
			event := m.calendar.GetSelectedEvent()
			if event != nil {
				return m.showDeleteConfirmEvent(event), nil
			}
			return m, nil

		case "e":
			event := m.calendar.GetSelectedEvent()
			if event != nil {
				return m.editEvent(event)
			}
			return m, nil
		}
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m.currentSection = "quests"
		return m, nil

	case "shift+tab":
		m.currentSection = "spaces"
		return m, nil

	case "?":
		m.mode = HelpView
		return m, nil

	case "r":
		return m.refreshData(), nil

	case "n":
		return m.createNewEvent()

	case "enter":
		m.calendar.FocusEventList()
		return m, nil

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
		m.currentSection = "journeys"
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
			return m.toggleQuest(item.quest)
		}

	case "d":
		if item, ok := m.journeyQuestList.SelectedItem().(questItem); ok {
			return m.showDeleteConfirm(item.quest), nil
		}

	case "e":
		if item, ok := m.journeyQuestList.SelectedItem().(questItem); ok {
			return m.editQuest(item.quest)
		}

	default:
		var cmd tea.Cmd
		m.journeyQuestList, cmd = m.journeyQuestList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleSpaceDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		m.mode = QuestListView
		m.currentSection = "spaces"
		m.selectedSpace = nil
		m.spaceSection = "quests"
		return m, nil

	case "tab":
		switch m.spaceSection {
		case "quests":
			m.spaceSection = "journeys"
		case "journeys":
			m.spaceSection = "members"
		default:
			m.spaceSection = "quests"
		}
		return m, nil

	case "shift+tab":
		switch m.spaceSection {
		case "quests":
			m.spaceSection = "members"
		case "journeys":
			m.spaceSection = "quests"
		default:
			m.spaceSection = "journeys"
		}
		return m, nil

	case "?":
		m.mode = HelpView
		return m, nil

	case "r":
		return m.refreshData(), nil
	}

	switch m.spaceSection {
	case "quests":
		switch msg.String() {
		case "n":
			return m.createNewQuestInSpace()
		case " ", "enter":
			if item, ok := m.spaceQuestList.SelectedItem().(questItem); ok {
				return m.toggleQuest(item.quest)
			}
		case "d":
			if item, ok := m.spaceQuestList.SelectedItem().(questItem); ok {
				return m.showDeleteConfirm(item.quest), nil
			}
		case "e":
			if item, ok := m.spaceQuestList.SelectedItem().(questItem); ok {
				return m.editQuest(item.quest)
			}
		default:
			var cmd tea.Cmd
			m.spaceQuestList, cmd = m.spaceQuestList.Update(msg)
			return m, cmd
		}
	case "journeys":
		switch msg.String() {
		case "n":
			return m.createNewJourneyInSpace()
		case " ", "enter":
			if item, ok := m.spaceJourneyList.SelectedItem().(journeyItem); ok {
				return m.enterJourney(item.journey), nil
			}
		case "d":
			if item, ok := m.spaceJourneyList.SelectedItem().(journeyItem); ok {
				return m.showDeleteConfirmJourney(item.journey), nil
			}
		case "e":
			if item, ok := m.spaceJourneyList.SelectedItem().(journeyItem); ok {
				return m.editJourney(item.journey)
			}
		default:
			var cmd tea.Cmd
			m.spaceJourneyList, cmd = m.spaceJourneyList.Update(msg)
			return m, cmd
		}
	case "members":
		var cmd tea.Cmd
		m.spaceMemberList, cmd = m.spaceMemberList.Update(msg)
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
				return m.confirmDeleteQuest()
			} else if m.confirmHabit != nil {
				return m.confirmDeleteHabit()
			} else if m.confirmJourney != nil {
				return m.confirmDeleteJourney()
			} else if m.confirmEvent != nil {
				return m.confirmDeleteEvent()
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
		} else if m.selectedSpace != nil {
			returnMode = SpaceDetailView
		} else {
			returnMode = QuestListView
		}
	case HabitFormView:
		form = m.habitForm
		returnMode = QuestListView
	case EventFormView:
		form = m.eventForm
		returnMode = QuestListView
		m.currentSection = "calendar"
	case SpaceFormView:
		form = m.spaceForm
		returnMode = QuestListView
		m.currentSection = "spaces"
	case QuestEditFormView:
		form = m.questForm
		if m.selectedJourney != nil {
			returnMode = JourneyDetailView
		} else if m.selectedSpace != nil {
			returnMode = SpaceDetailView
		} else {
			returnMode = QuestListView
		}
	case JourneyEditFormView:
		form = m.journeyForm
		if m.selectedSpace != nil {
			returnMode = SpaceDetailView
		} else {
			returnMode = QuestListView
		}
	case HabitEditFormView:
		form = m.habitForm
		returnMode = QuestListView
	case EventEditFormView:
		form = m.eventForm
		returnMode = QuestListView
		m.currentSection = "calendar"
	case SpaceEditFormView:
		form = m.spaceForm
		if m.selectedSpace != nil {
			returnMode = SpaceDetailView
		} else {
			returnMode = QuestListView
		}
		m.currentSection = "spaces"
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
		case EventFormView:
			m.eventForm = f
		case SpaceFormView:
			m.spaceForm = f
		case QuestEditFormView:
			m.questForm = f
		case JourneyEditFormView:
			m.journeyForm = f
		case HabitEditFormView:
			m.habitForm = f
		case EventEditFormView:
			m.eventForm = f
		case SpaceEditFormView:
			m.spaceForm = f
		}

		if f.State == huh.StateCompleted {
			return m.handleFormCompletion()
		}

		if f.State == huh.StateAborted {
			m.mode = returnMode
			m.message = "Form cancelled"
			return m, nil
		}
	}

	return m, cmd
}
