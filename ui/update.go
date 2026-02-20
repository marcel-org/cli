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

	// Handle forms first - they need all messages
	if m.mode == QuestFormView || m.mode == JourneyFormView || m.mode == HabitFormView || m.mode == EventFormView {
		return m.handleFormUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.questList = newQuestList(m.data, m.width-4, m.height-10)
			m.habitList = newHabitList(m.data, m.width-4, m.height-10)
			m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
			m.calendar.SetSize(m.width-4, m.height-10)
			m.calendar.SetEvents(m.data.Events)
			m.ready = true
		} else {
			m.questList.SetSize(m.width-4, m.height-10)
			m.habitList.SetSize(m.width-4, m.height-10)
			m.journeyList.SetSize(m.width-4, m.height-10)
			m.calendar.SetSize(m.width-4, m.height-10)
		}

	case tea.KeyMsg:
		switch m.mode {
		case QuestListView:
			switch m.currentSection {
			case "quests":
				return m.handleQuestListKeys(msg)
			case "habits":
				return m.handleHabitListKeys(msg)
			case "journeys":
				return m.handleJourneyListKeys(msg)
			case "calendar":
				return m.handleCalendarKeys(msg)
			}
		case JourneyDetailView:
			return m.handleJourneyDetailKeys(msg)
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

	case clearMessageMsg:
		m.message = ""
	}

	return m, tea.Batch(cmds...)
}
