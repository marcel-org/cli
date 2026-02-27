package ui

import (
	"fmt"
	"marcel-cli/models"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.needsRedraw {
		m.needsRedraw = false
		cmds = append(cmds, tea.ClearScreen)
	}

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

		if m.syncStatus == SyncStatusSyncing {
			var syncCmd tea.Cmd
			m.syncSpinner, syncCmd = m.syncSpinner.Update(msg)
			cmds = append(cmds, syncCmd)
		}

	case dataLoadedMsg:
		if msg.err != nil {
			cmds = append(cmds, checkAuthCmd(m.storage))
		} else {
			m.data = msg.data
			m.mode = QuestListView
			m.currentSection = msg.data.CurrentSection
			m.syncStatus = SyncStatusSyncing
			cmds = append(cmds, backgroundSyncCmd(m.storage), m.syncSpinner.Tick)
		}

	case authCheckMsg:
		if msg.err != nil {
			m.mode = ErrorView
			m.errorMessage = fmt.Sprintf("Authentication failed: %v\n\nSet your token in ~/.marcel.token or MARCEL_TOKEN environment variable", msg.err)
		} else {
			cmds = append(cmds, loadFromAPICmd(m.storage))
		}

	case backgroundSyncMsg:
		if msg.err != nil {
			m.syncStatus = SyncStatusError
		} else {
			m.data = msg.data
			m.syncStatus = SyncStatusSynced
			m.questList = newQuestList(m.data, m.width-4, m.height-10)
			m.habitList = newHabitList(m.data, m.width-4, m.height-10)
			m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
			m.calendar.SetEvents(m.data.Events)

			var allQuests []models.Quest
			for _, journey := range m.data.Journeys {
				allQuests = append(allQuests, journey.Quests...)
			}
			m.storage.SaveToCache(m.data.Journeys, allQuests, m.data.Habits, m.data.Events)
			cmds = append(cmds, clearSyncStatusAfter(3*time.Second))
		}

	case clearMessageMsg:
		m.message = ""

	case clearSyncStatusMsg:
		m.syncStatus = SyncStatusNone
	}

	return m, tea.Batch(cmds...)
}
