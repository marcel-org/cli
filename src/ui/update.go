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

	if m.mode == QuestFormView || m.mode == JourneyFormView || m.mode == HabitFormView || m.mode == EventFormView || m.mode == SpaceFormView ||
		m.mode == QuestEditFormView || m.mode == JourneyEditFormView || m.mode == HabitEditFormView || m.mode == EventEditFormView || m.mode == SpaceEditFormView {
		return m.handleFormUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.rebuildLists()
			m.calendar.SetSize(m.width-4, m.height-10)
			m.calendar.SetEvents(m.data.Events)
			m.ready = true
		} else {
			w, h := m.width-4, m.height-10
			m.questList.SetSize(w, h)
			m.habitList.SetSize(w, h)
			m.journeyList.SetSize(w, h)
			m.spaceList.SetSize(w, h)
			if m.selectedSpace != nil {
				m.spaceQuestList.SetSize(w, h)
				m.spaceJourneyList.SetSize(w, h)
				m.spaceMemberList.SetSize(w, h)
			}
			if m.selectedJourney != nil {
				m.journeyQuestList.SetSize(w, h)
			}
			m.calendar.SetSize(w, h)
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
			case "spaces":
				return m.handleSpaceListKeys(msg)
			case "calendar":
				return m.handleCalendarKeys(msg)
			}
		case JourneyDetailView:
			return m.handleJourneyDetailKeys(msg)
		case SpaceDetailView:
			return m.handleSpaceDetailKeys(msg)
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
			m.rebuildLists()
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
			m.rebuildLists()
			m.calendar.SetEvents(m.data.Events)

			var allQuests []models.Quest
			for _, journey := range m.data.Journeys {
				allQuests = append(allQuests, journey.Quests...)
			}
			m.storage.SaveToCache(m.data.Journeys, allQuests, m.data.Habits, m.data.Events, m.data.Spaces)
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
			m.rebuildLists()
			m.needsRedraw = true
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.quest != nil {
			for i := range m.data.Journeys {
				for j := range m.data.Journeys[i].Quests {
					if m.data.Journeys[i].Quests[j].ID == msg.tempID {
						m.data.Journeys[i].Quests = append(m.data.Journeys[i].Quests[:j], m.data.Journeys[i].Quests[j+1:]...)
						m.data.Journeys[i].Quests = append([]models.Quest{*msg.quest}, m.data.Journeys[i].Quests...)
						break
					}
				}
			}
			m.rebuildLists()
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
			m.rebuildLists()
			m.needsRedraw = true
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.journey != nil {
			for i := range m.data.Journeys {
				if m.data.Journeys[i].ID == msg.tempID {
					m.data.Journeys[i] = *msg.journey
					break
				}
			}
			m.rebuildLists()
		}

	case spaceCreatedMsg:
		if msg.err != nil {
			for i := range m.data.Spaces {
				if m.data.Spaces[i].ID == msg.tempID {
					m.data.Spaces = append(m.data.Spaces[:i], m.data.Spaces[i+1:]...)
					break
				}
			}
			m.message = fmt.Sprintf("Failed to create space: %v", msg.err)
			m.rebuildLists()
			m.needsRedraw = true
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.space != nil {
			for i := range m.data.Spaces {
				if m.data.Spaces[i].ID == msg.tempID {
					m.data.Spaces[i] = *msg.space
					break
				}
			}
			m.rebuildLists()
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
			m.rebuildLists()
			m.needsRedraw = true
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.habit != nil {
			for i := range m.data.Habits {
				if m.data.Habits[i].ID == msg.tempID {
					m.data.Habits[i] = *msg.habit
					break
				}
			}
			m.rebuildLists()
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
			m.syncCalendarEvents()
			m.needsRedraw = true
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.event != nil {
			m.upsertEvent(msg.tempID, *msg.event)
		}

	case questToggledMsg:
		if msg.err != nil {
			for i := range m.data.Journeys {
				for j := range m.data.Journeys[i].Quests {
					if m.data.Journeys[i].Quests[j].ID == msg.questID {
						m.data.Journeys[i].Quests[j].Done = !msg.done
					}
				}
			}
			m.rebuildLists()
			m.needsRedraw = true
			m.message = fmt.Sprintf("Failed to toggle quest: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		}

	case questDeletedMsg:
		if msg.err != nil {
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
				if !found && m.confirmQuest.JourneyID == nil {
					myQuests := models.Journey{
						ID:     0,
						Name:   "My Quests",
						Quests: []models.Quest{*m.confirmQuest},
					}
					m.data.Journeys = append([]models.Journey{myQuests}, m.data.Journeys...)
				}
				m.rebuildLists()
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
			if m.confirmHabit != nil {
				m.data.Habits = append(m.data.Habits, *m.confirmHabit)
				m.rebuildLists()
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
			if m.confirmJourney != nil {
				m.data.Journeys = append(m.data.Journeys, *m.confirmJourney)
				m.rebuildLists()
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
			if m.editingQuest != nil {
				for i := range m.data.Journeys {
					for j := range m.data.Journeys[i].Quests {
						if m.data.Journeys[i].Quests[j].ID == msg.questID {
							m.data.Journeys[i].Quests[j] = *m.editingQuest
							break
						}
					}
				}
				m.rebuildLists()
				m.needsRedraw = true
			}
			m.message = fmt.Sprintf("Failed to update quest: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.quest != nil {
			for i := range m.data.Journeys {
				for j := range m.data.Journeys[i].Quests {
					if m.data.Journeys[i].Quests[j].ID == msg.questID {
						m.data.Journeys[i].Quests[j] = *msg.quest
						break
					}
				}
			}
			m.rebuildLists()
			m.message = "Quest updated successfully"
			cmds = append(cmds, clearMessageAfter(1*time.Second))
		}
		m.editingQuest = nil

	case habitUpdatedMsg:
		if msg.err != nil {
			if m.editingHabit != nil {
				for i := range m.data.Habits {
					if m.data.Habits[i].ID == msg.habitID {
						m.data.Habits[i] = *m.editingHabit
						break
					}
				}
				m.rebuildLists()
				m.needsRedraw = true
			}
			m.message = fmt.Sprintf("Failed to update habit: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.habit != nil {
			for i := range m.data.Habits {
				if m.data.Habits[i].ID == msg.habitID {
					m.data.Habits[i] = *msg.habit
					break
				}
			}
			m.rebuildLists()
			m.message = "Habit updated successfully"
			cmds = append(cmds, clearMessageAfter(1*time.Second))
		}
		m.editingHabit = nil

	case journeyUpdatedMsg:
		if msg.err != nil {
			if m.editingJourney != nil {
				for i := range m.data.Journeys {
					if m.data.Journeys[i].ID == msg.journeyID {
						m.data.Journeys[i] = *m.editingJourney
						break
					}
				}
				m.rebuildLists()
				m.needsRedraw = true
			}
			m.message = fmt.Sprintf("Failed to update journey: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.journey != nil {
			for i := range m.data.Journeys {
				if m.data.Journeys[i].ID == msg.journeyID {
					m.data.Journeys[i] = *msg.journey
					break
				}
			}
			m.rebuildLists()
			m.message = "Journey updated successfully"
			cmds = append(cmds, clearMessageAfter(1*time.Second))
		}
		m.editingJourney = nil

	case spaceUpdatedMsg:
		if msg.err != nil {
			if m.editingSpace != nil {
				for i := range m.data.Spaces {
					if m.data.Spaces[i].ID == msg.spaceID {
						m.data.Spaces[i] = *m.editingSpace
						break
					}
				}
				m.rebuildLists()
				m.needsRedraw = true
			}
			m.message = fmt.Sprintf("Failed to update space: %v", msg.err)
			cmds = append(cmds, clearMessageAfter(3*time.Second))
		} else if msg.space != nil {
			for i := range m.data.Spaces {
				if m.data.Spaces[i].ID == msg.spaceID {
					m.data.Spaces[i] = *msg.space
					break
				}
			}
			m.rebuildLists()
			m.message = "Space updated successfully"
			cmds = append(cmds, clearMessageAfter(1*time.Second))
		}
		m.editingSpace = nil

	case clearMessageMsg:
		m.message = ""

	case clearSyncStatusMsg:
		m.syncStatus = SyncStatusNone
	}

	return m, tea.Batch(cmds...)
}
