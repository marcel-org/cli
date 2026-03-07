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

	if m.mode == QuestFormView || m.mode == JourneyFormView || m.mode == HabitFormView || m.mode == EventFormView ||
		m.mode == QuestEditFormView || m.mode == JourneyEditFormView || m.mode == HabitEditFormView || m.mode == EventEditFormView {
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
			m.questList = newQuestList(m.data, m.width-4, m.height-10)
			m.habitList = newHabitList(m.data, m.width-4, m.height-10)
			m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
			m.calendar.SetEvents(m.data.Events)
			m.syncStatus = SyncStatusSyncing
			cmds = append(cmds, backgroundSyncCmd(m.storage), m.syncSpinner.Tick)
		}

	case authCheckMsg:
		if msg.err != nil {
			m.mode = ErrorView
			m.errorMessage = fmt.Sprintf("Authentication failed: %v\n\nSet your token in ~/.marcel.token or MARCEL_TOKEN environment variable", msg.err)
		} else {
			m.syncStatus = SyncStatusSyncing
			cmds = append(cmds, backgroundSyncCmd(m.storage), m.syncSpinner.Tick)
		}

	case backgroundSyncMsg:
		if msg.err != nil {
			m.syncStatus = SyncStatusError
			if m.mode == LoadingView {
				m.mode = ErrorView
				m.errorMessage = fmt.Sprintf("Failed to load data: %v", msg.err)
			}
		} else {
			m.data = msg.data
			m.syncStatus = SyncStatusSynced
			if m.mode == LoadingView {
				m.mode = QuestListView
				m.currentSection = msg.data.CurrentSection
			}
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

	case questCreatedMsg:
		if msg.err != nil {
			for i := range m.data.Journeys {
				for j := range m.data.Journeys[i].Quests {
					if m.data.Journeys[i].Quests[j].ID == msg.tempID {
						m.data.Journeys[i].Quests = append(
							m.data.Journeys[i].Quests[:j],
							m.data.Journeys[i].Quests[j+1:]...,
						)
						if m.data.Journeys[i].ID == 0 && len(m.data.Journeys[i].Quests) == 0 {
							m.data.Journeys = append(m.data.Journeys[:i], m.data.Journeys[i+1:]...)
						}
						break
					}
				}
			}
			m.message = fmt.Sprintf("Failed to create quest: %v", msg.err)
			m.questList = newQuestList(m.data, m.width-4, m.height-10)
			m.needsRedraw = true
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.quest != nil {
			for i := range m.data.Journeys {
				for j := range m.data.Journeys[i].Quests {
					if m.data.Journeys[i].Quests[j].ID == msg.tempID {
						m.data.Journeys[i].Quests[j] = *msg.quest
						break
					}
				}
			}
			m.questList = newQuestList(m.data, m.width-4, m.height-10)
			if m.selectedJourney != nil {
				for _, j := range m.data.Journeys {
					if j.ID == m.selectedJourney.ID {
						m.selectedJourney = &j
						m.journeyQuestList = newJourneyQuestList(&j, m.width-4, m.height-10)
						break
					}
				}
			}
		}

	case journeyCreatedMsg:
		if msg.err != nil {
			for i := range m.data.Journeys {
				if m.data.Journeys[i].ID == msg.tempID {
					m.data.Journeys = append(m.data.Journeys[:i], m.data.Journeys[i+1:]...)
					break
				}
			}
			m.message = fmt.Sprintf("Failed to create journey: %v", msg.err)
			m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
			m.needsRedraw = true
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.journey != nil {
			for i := range m.data.Journeys {
				if m.data.Journeys[i].ID == msg.tempID {
					m.data.Journeys[i] = *msg.journey
					break
				}
			}
			m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
		}

	case habitCreatedMsg:
		if msg.err != nil {
			for i := range m.data.Habits {
				if m.data.Habits[i].ID == msg.tempID {
					m.data.Habits = append(m.data.Habits[:i], m.data.Habits[i+1:]...)
					break
				}
			}
			m.message = fmt.Sprintf("Failed to create habit: %v", msg.err)
			m.habitList = newHabitList(m.data, m.width-4, m.height-10)
			m.needsRedraw = true
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.habit != nil {
			for i := range m.data.Habits {
				if m.data.Habits[i].ID == msg.tempID {
					m.data.Habits[i] = *msg.habit
					break
				}
			}
			m.habitList = newHabitList(m.data, m.width-4, m.height-10)
		}

	case eventCreatedMsg:
		if msg.err != nil {
			for i := range m.data.Events {
				if m.data.Events[i].ID == msg.tempID {
					m.data.Events = append(m.data.Events[:i], m.data.Events[i+1:]...)
					break
				}
			}
			m.message = fmt.Sprintf("Failed to create event: %v", msg.err)
			m.calendar.SetEvents(m.data.Events)
			m.needsRedraw = true
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.event != nil {
			for i := range m.data.Events {
				if m.data.Events[i].ID == msg.tempID {
					m.data.Events[i] = *msg.event
					break
				}
			}
			m.calendar.SetEvents(m.data.Events)
		}

	case questDeletedMsg:
		if msg.err != nil {
			// Rollback: restore the deleted quest
			if m.confirmQuest != nil {
				found := false
				for i := range m.data.Journeys {
					if m.confirmQuest.JourneyID != nil && m.data.Journeys[i].ID == *m.confirmQuest.JourneyID {
						m.data.Journeys[i].Quests = append(m.data.Journeys[i].Quests, *m.confirmQuest)
						found = true
						break
					} else if m.confirmQuest.JourneyID == nil && m.data.Journeys[i].ID == 0 {
						m.data.Journeys[i].Quests = append(m.data.Journeys[i].Quests, *m.confirmQuest)
						found = true
						break
					}
				}
				// If journey not found, create "My Quests" journey
				if !found && m.confirmQuest.JourneyID == nil {
					myQuests := models.Journey{
						ID:     0,
						Name:   "My Quests",
						Quests: []models.Quest{*m.confirmQuest},
					}
					m.data.Journeys = append([]models.Journey{myQuests}, m.data.Journeys...)
				}
				m.questList = newQuestList(m.data, m.width-4, m.height-10)
				m.needsRedraw = true
			}
			m.message = fmt.Sprintf("Failed to delete quest: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else {
			m.message = "Quest deleted successfully"
			cmds = append(cmds, clearMessageAfter(1*time.Second))
		}
		m.confirmQuest = nil

	case habitDeletedMsg:
		if msg.err != nil {
			// Rollback: restore the deleted habit
			if m.confirmHabit != nil {
				m.data.Habits = append(m.data.Habits, *m.confirmHabit)
				m.habitList = newHabitList(m.data, m.width-4, m.height-10)
				m.needsRedraw = true
			}
			m.message = fmt.Sprintf("Failed to delete habit: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else {
			m.message = "Habit deleted successfully"
			cmds = append(cmds, clearMessageAfter(1*time.Second))
		}
		m.confirmHabit = nil

	case journeyDeletedMsg:
		if msg.err != nil {
			// Rollback: restore the deleted journey
			if m.confirmJourney != nil {
				m.data.Journeys = append(m.data.Journeys, *m.confirmJourney)
				m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
				m.needsRedraw = true
			}
			m.message = fmt.Sprintf("Failed to delete journey: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else {
			m.message = "Journey deleted successfully"
			cmds = append(cmds, clearMessageAfter(1*time.Second))
		}
		m.confirmJourney = nil

	case questUpdatedMsg:
		if msg.err != nil {
			// Rollback: restore the original quest
			if m.editingQuest != nil {
				for i := range m.data.Journeys {
					for j := range m.data.Journeys[i].Quests {
						if m.data.Journeys[i].Quests[j].ID == msg.questID {
							m.data.Journeys[i].Quests[j] = *m.editingQuest
							break
						}
					}
				}
				m.questList = newQuestList(m.data, m.width-4, m.height-10)
				if m.selectedJourney != nil {
					for _, j := range m.data.Journeys {
						if j.ID == m.selectedJourney.ID {
							m.selectedJourney = &j
							m.journeyQuestList = newJourneyQuestList(&j, m.width-4, m.height-10)
							break
						}
					}
				}
				m.needsRedraw = true
			}
			m.message = fmt.Sprintf("Failed to update quest: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.quest != nil {
			// Replace with server data
			for i := range m.data.Journeys {
				for j := range m.data.Journeys[i].Quests {
					if m.data.Journeys[i].Quests[j].ID == msg.questID {
						m.data.Journeys[i].Quests[j] = *msg.quest
						break
					}
				}
			}
			m.questList = newQuestList(m.data, m.width-4, m.height-10)
			if m.selectedJourney != nil {
				for _, j := range m.data.Journeys {
					if j.ID == m.selectedJourney.ID {
						m.selectedJourney = &j
						m.journeyQuestList = newJourneyQuestList(&j, m.width-4, m.height-10)
						break
					}
				}
			}
			m.message = "Quest updated successfully"
			cmds = append(cmds, clearMessageAfter(1*time.Second))
		}
		m.editingQuest = nil

	case habitUpdatedMsg:
		if msg.err != nil {
			// Rollback: restore the original habit
			if m.editingHabit != nil {
				for i := range m.data.Habits {
					if m.data.Habits[i].ID == msg.habitID {
						m.data.Habits[i] = *m.editingHabit
						break
					}
				}
				m.habitList = newHabitList(m.data, m.width-4, m.height-10)
				m.needsRedraw = true
			}
			m.message = fmt.Sprintf("Failed to update habit: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.habit != nil {
			// Replace with server data
			for i := range m.data.Habits {
				if m.data.Habits[i].ID == msg.habitID {
					m.data.Habits[i] = *msg.habit
					break
				}
			}
			m.habitList = newHabitList(m.data, m.width-4, m.height-10)
			m.message = "Habit updated successfully"
			cmds = append(cmds, clearMessageAfter(1*time.Second))
		}
		m.editingHabit = nil

	case journeyUpdatedMsg:
		if msg.err != nil {
			// Rollback: restore the original journey
			if m.editingJourney != nil {
				for i := range m.data.Journeys {
					if m.data.Journeys[i].ID == msg.journeyID {
						m.data.Journeys[i] = *m.editingJourney
						break
					}
				}
				m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
				m.needsRedraw = true
			}
			m.message = fmt.Sprintf("Failed to update journey: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.journey != nil {
			// Replace with server data
			for i := range m.data.Journeys {
				if m.data.Journeys[i].ID == msg.journeyID {
					m.data.Journeys[i] = *msg.journey
					break
				}
			}
			m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
			m.message = "Journey updated successfully"
			cmds = append(cmds, clearMessageAfter(1*time.Second))
		}
		m.editingJourney = nil

	case clearMessageMsg:
		m.message = ""

	case clearSyncStatusMsg:
		m.syncStatus = SyncStatusNone
	}

	return m, tea.Batch(cmds...)
}
