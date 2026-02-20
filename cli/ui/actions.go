package ui

import (
	"fmt"
	"marcel-cli/api"
	"marcel-cli/models"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) toggleQuest(quest models.Quest) (Model, tea.Cmd) {
	newDone := !quest.Done

	_, err := m.storage.GetAPIClient().ToggleQuest(quest.ID, newDone)
	if err != nil {
		m.message = fmt.Sprintf("Failed to toggle quest: %v", err)
		return m, nil
	}

	for i := range m.data.Journeys {
		for j := range m.data.Journeys[i].Quests {
			if m.data.Journeys[i].Quests[j].ID == quest.ID {
				m.data.Journeys[i].Quests[j].Done = newDone
			}
		}
	}

	m.questList = newQuestList(m.data, m.width-4, m.height-8)

	if newDone {
		m.message = "✓ Quest completed!"
	} else {
		m.message = "Quest marked as incomplete"
	}

	return m, clearMessageAfter(1 * time.Second)
}

func (m Model) showDeleteConfirm(quest models.Quest) Model {
	m.mode = ConfirmDeleteView
	m.confirmQuest = &quest
	m.confirmSelected = false
	return m
}

func (m Model) confirmDeleteQuest() (Model, tea.Cmd) {
	if m.confirmQuest == nil {
		m.mode = QuestListView
		return m, nil
	}

	err := m.storage.GetAPIClient().DeleteQuest(m.confirmQuest.ID)
	if err != nil {
		m.message = fmt.Sprintf("Failed to delete quest: %v", err)
		m.mode = QuestListView
		m.confirmQuest = nil
		return m, nil
	}

	for i := range m.data.Journeys {
		newQuests := []models.Quest{}
		for _, q := range m.data.Journeys[i].Quests {
			if q.ID != m.confirmQuest.ID {
				newQuests = append(newQuests, q)
			}
		}
		m.data.Journeys[i].Quests = newQuests
	}

	m.questList = newQuestList(m.data, m.width-4, m.height-8)
	m.message = "Quest deleted successfully"
	m.mode = QuestListView
	m.confirmQuest = nil

	return m, clearMessageAfter(1 * time.Second)
}

func (m Model) confirmDeleteHabit() (Model, tea.Cmd) {
	if m.confirmHabit == nil {
		m.mode = QuestListView
		return m, nil
	}

	err := m.storage.GetAPIClient().DeleteHabit(m.confirmHabit.ID)
	if err != nil {
		m.message = fmt.Sprintf("Failed to delete habit: %v", err)
		m.mode = QuestListView
		m.confirmHabit = nil
		return m, nil
	}

	var newHabits []models.Habit
	for _, h := range m.data.Habits {
		if h.ID != m.confirmHabit.ID {
			newHabits = append(newHabits, h)
		}
	}
	m.data.Habits = newHabits

	m.habitList = newHabitList(m.data, m.width-4, m.height-10)
	m.message = "Habit deleted successfully"
	m.mode = QuestListView
	m.confirmHabit = nil

	return m, clearMessageAfter(1 * time.Second)
}

func (m Model) confirmDeleteJourney() (Model, tea.Cmd) {
	if m.confirmJourney == nil {
		m.mode = QuestListView
		return m, nil
	}

	err := m.storage.GetAPIClient().DeleteJourney(m.confirmJourney.ID)
	if err != nil {
		m.message = fmt.Sprintf("Failed to delete journey: %v", err)
		m.mode = QuestListView
		m.confirmJourney = nil
		return m, nil
	}

	var newJourneys []models.Journey
	for _, j := range m.data.Journeys {
		if j.ID != m.confirmJourney.ID {
			newJourneys = append(newJourneys, j)
		}
	}
	m.data.Journeys = newJourneys

	m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
	m.message = "Journey deleted successfully"
	m.mode = QuestListView
	m.confirmJourney = nil

	return m, clearMessageAfter(1 * time.Second)
}

func (m Model) cancelDelete() Model {
	returnToCalendar := m.confirmEvent != nil

	m.mode = QuestListView
	m.confirmQuest = nil
	m.confirmHabit = nil
	m.confirmJourney = nil
	m.confirmEvent = nil
	m.message = "Deletion cancelled"

	if returnToCalendar {
		m.currentSection = "calendar"
	}

	return m
}

func (m Model) createNewQuest() (Model, tea.Cmd) {
	m.questFormData = &QuestForm{
		Title:      "",
		Note:       "",
		Difficulty: "medium",
		JourneyID:  nil,
	}
	m.questForm = BuildQuestForm(m.questFormData, m.data.Journeys)
	m.mode = QuestFormView
	return m, m.questForm.Init()
}

func (m Model) createNewQuestInJourney() (Model, tea.Cmd) {
	if m.selectedJourney == nil {
		return m, nil
	}

	m.questFormData = &QuestForm{
		Title:      "",
		Note:       "",
		Difficulty: "medium",
		JourneyID:  nil,
	}
	m.questForm = BuildQuestForm(m.questFormData, m.data.Journeys)
	m.mode = QuestFormView
	return m, m.questForm.Init()
}

func (m Model) createNewHabit() (Model, tea.Cmd) {
	m.habitFormData = &HabitForm{
		Name:        "",
		CycleType:   "daily",
		CycleConfig: nil,
	}
	m.habitForm = BuildHabitForm(m.habitFormData)
	m.mode = HabitFormView
	return m, m.habitForm.Init()
}

func (m Model) createNewJourney() (Model, tea.Cmd) {
	m.journeyFormData = &JourneyForm{
		Name: "",
	}
	m.journeyForm = BuildJourneyForm(m.journeyFormData)
	m.mode = JourneyFormView
	return m, m.journeyForm.Init()
}

func (m Model) enterJourney(journey models.Journey) Model {
	m.selectedJourney = &journey
	m.journeyQuestList = newJourneyQuestList(&journey, m.width-4, m.height-10)
	m.mode = JourneyDetailView
	return m
}

func (m Model) refreshData() Model {
	m.mode = LoadingView
	m.message = "Refreshing data..."

	data, err := m.storage.Load()
	if err != nil {
		m.mode = ErrorView
		m.errorMessage = fmt.Sprintf("Failed to load data: %v", err)
		return m
	}

	m.data = data
	m.questList = newQuestList(m.data, m.width-4, m.height-10)
	m.habitList = newHabitList(m.data, m.width-4, m.height-10)
	m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
	m.calendar.SetEvents(m.data.Events)
	m.mode = QuestListView
	m.message = "✓ Data refreshed!"

	return m
}

func (m Model) toggleHabit(habit models.Habit) (Model, tea.Cmd) {
	completedToday := false
	today := fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	for _, d := range habit.Completed {
		if len(d) >= 10 && d[:10] == today {
			completedToday = true
			break
		}
	}

	newDone := !completedToday

	_, err := m.storage.GetAPIClient().ToggleHabit(habit.ID, newDone)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "not scheduled for today") {
			parts := strings.Split(errMsg, "It's configured for: ")
			if len(parts) > 1 {
				configPart := strings.Split(parts[1], ". Next due:")
				if len(configPart) > 1 {
					nextDue := strings.TrimSpace(configPart[1])
					nextDue = strings.TrimSuffix(nextDue, ".")
					nextDue = strings.TrimSuffix(nextDue, "\"}")
					nextDue = strings.TrimSuffix(nextDue, "}")
					m.message = fmt.Sprintf("Not due today. Next: %s", nextDue)
				} else {
					m.message = "This habit is not scheduled for today"
				}
			} else {
				m.message = "This habit is not scheduled for today"
			}
		} else {
			m.message = fmt.Sprintf("Failed to toggle habit: %v", err)
		}
		return m, nil
	}

	for i := range m.data.Habits {
		if m.data.Habits[i].ID == habit.ID {
			if newDone {
				m.data.Habits[i].Completed = append(m.data.Habits[i].Completed, today)
				m.data.Habits[i].CurrentStreak++
			} else {
				var newCompleted []string
				for _, d := range m.data.Habits[i].Completed {
					if len(d) < 10 || d[:10] != today {
						newCompleted = append(newCompleted, d)
					}
				}
				m.data.Habits[i].Completed = newCompleted
				m.data.Habits[i].CurrentStreak--
			}
		}
	}

	m.habitList = newHabitList(m.data, m.width-4, m.height-10)

	if newDone {
		m.message = "✓ Habit completed!"
	} else {
		m.message = "Habit marked as incomplete"
	}

	return m, clearMessageAfter(1 * time.Second)
}

func (m Model) showDeleteConfirmHabit(habit models.Habit) Model {
	m.mode = ConfirmDeleteView
	m.confirmHabit = &habit
	m.confirmSelected = false
	return m
}

func (m Model) showDeleteConfirmJourney(journey models.Journey) Model {
	m.mode = ConfirmDeleteView
	m.confirmJourney = &journey
	m.confirmSelected = false
	return m
}

func (m Model) createNewEvent() (Model, tea.Cmd) {
	m.eventFormData = &EventForm{
		Title:       "",
		Date:        m.calendar.GetSelectedDate().Format("2006-01-02"),
		Time:        "",
		EndTime:     "",
		Location:    "",
		Description: "",
	}
	m.eventForm = BuildEventForm(m.eventFormData)
	m.mode = EventFormView
	return m, m.eventForm.Init()
}

func (m Model) showDeleteConfirmEvent(event *models.Event) Model {
	m.mode = ConfirmDeleteView
	m.confirmEvent = event
	m.confirmSelected = false
	return m
}

func (m Model) confirmDeleteEvent() (Model, tea.Cmd) {
	if m.confirmEvent == nil {
		m.mode = QuestListView
		m.currentSection = "calendar"
		return m, nil
	}

	err := m.storage.GetAPIClient().DeleteEvent(m.confirmEvent.ID)
	if err != nil {
		m.message = fmt.Sprintf("Failed to delete event: %v", err)
		m.mode = QuestListView
		m.currentSection = "calendar"
		m.confirmEvent = nil
		return m, nil
	}

	var newEvents []models.Event
	for _, e := range m.data.Events {
		if e.ID != m.confirmEvent.ID {
			newEvents = append(newEvents, e)
		}
	}
	m.data.Events = newEvents
	m.calendar.SetEvents(newEvents)

	m.message = "Event deleted successfully"
	m.mode = QuestListView
	m.currentSection = "calendar"
	m.confirmEvent = nil
	return m, clearMessageAfter(1 * time.Second)
}

func (m Model) handleFormCompletion() (tea.Model, tea.Cmd) {
	var returnMode ViewMode = QuestListView
	var message string

	switch m.mode {
	case QuestFormView:
		if m.questFormData.Title == "" {
			message = "Quest title cannot be empty"
			m.mode = returnMode
			m.message = message
			return m, nil
		}

		var journeyID *int
		if m.questFormData.JourneyID != nil && *m.questFormData.JourneyID != 0 {
			journeyID = m.questFormData.JourneyID
		} else if m.selectedJourney != nil {
			journeyID = &m.selectedJourney.ID
		}

		quest, err := m.storage.GetAPIClient().CreateQuest(
			m.questFormData.Title,
			m.questFormData.Note,
			m.questFormData.Difficulty,
			journeyID,
		)

		if err != nil {
			message = fmt.Sprintf("Failed to create quest: %v", err)
			m.mode = returnMode
			m.message = message
			return m, nil
		}

		message = fmt.Sprintf("✓ Quest created: %s", quest.Title)

		if m.selectedJourney != nil {
			returnMode = JourneyDetailView
		}

	case JourneyFormView:
		if m.journeyFormData.Name == "" {
			message = "Journey name cannot be empty"
			m.mode = returnMode
			return m, nil
		}

		journey, err := m.storage.GetAPIClient().CreateJourney(m.journeyFormData.Name)
		if err != nil {
			message = fmt.Sprintf("Failed to create journey: %v", err)
			m.mode = returnMode
			m.message = message
			return m, nil
		}

		message = fmt.Sprintf("✓ Journey created: %s", journey.Name)

		if m.selectedJourney != nil {
			returnMode = JourneyDetailView
		}

	case HabitFormView:
		if m.habitFormData.Name == "" {
			message = "Habit name cannot be empty"
			m.mode = returnMode
			return m, nil
		}

		habit, err := m.storage.GetAPIClient().CreateHabit(
			m.habitFormData.Name,
			m.habitFormData.CycleType,
			m.habitFormData.CycleConfig,
		)

		if err != nil {
			message = fmt.Sprintf("Failed to create habit: %v", err)
			m.mode = returnMode
			m.message = message
			return m, nil
		}

		message = fmt.Sprintf("✓ Habit created: %s", habit.Name)

	case EventFormView:
		if m.eventFormData.Title == "" {
			message = "Event title cannot be empty"
			m.mode = returnMode
			m.currentSection = "calendar"
			return m, nil
		}

		var timePtr, endTimePtr, locationPtr, descriptionPtr *string
		if m.eventFormData.Time != "" {
			timePtr = &m.eventFormData.Time
		}
		if m.eventFormData.EndTime != "" {
			endTimePtr = &m.eventFormData.EndTime
		}
		if m.eventFormData.Location != "" {
			locationPtr = &m.eventFormData.Location
		}
		if m.eventFormData.Description != "" {
			descriptionPtr = &m.eventFormData.Description
		}

		event, err := m.storage.GetAPIClient().CreateEvent(api.CreateEventRequest{
			Title:       m.eventFormData.Title,
			Date:        m.eventFormData.Date,
			Time:        timePtr,
			EndTime:     endTimePtr,
			Location:    locationPtr,
			Description: descriptionPtr,
		})

		if err != nil {
			message = fmt.Sprintf("Failed to create event: %v", err)
			m.mode = returnMode
			m.currentSection = "calendar"
			m.message = message
			return m, nil
		}

		message = fmt.Sprintf("✓ Event created: %s", event.Title)

		m.currentSection = "calendar"
		returnMode = QuestListView
	}

	m = m.refreshData()

	if m.mode == ErrorView {
		return m, nil
	}

	if m.selectedJourney != nil && returnMode == JourneyDetailView {
		for _, j := range m.data.Journeys {
			if j.ID == m.selectedJourney.ID {
				m.selectedJourney = &j
				m.journeyQuestList = newJourneyQuestList(&j, m.width-4, m.height-10)
				break
			}
		}
	}

	m.mode = returnMode
	m.message = message
	m.needsRedraw = true

	return m, clearMessageAfter(1 * time.Second)
}
