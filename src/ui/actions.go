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

	questID := m.confirmQuest.ID

	// Optimistically remove quest from UI
	for i := range m.data.Journeys {
		newQuests := []models.Quest{}
		for _, q := range m.data.Journeys[i].Quests {
			if q.ID != questID {
				newQuests = append(newQuests, q)
			}
		}
		m.data.Journeys[i].Quests = newQuests
	}

	m.questList = newQuestList(m.data, m.width-4, m.height-8)
	m.message = "Deleting quest..."
	m.mode = QuestListView

	// Keep confirmQuest for rollback in case of error
	m.questDeleteCmd = deleteQuestCmd(m.storage, questID)

	return m, m.questDeleteCmd
}

func (m Model) confirmDeleteHabit() (Model, tea.Cmd) {
	if m.confirmHabit == nil {
		m.mode = QuestListView
		return m, nil
	}

	habitID := m.confirmHabit.ID

	// Optimistically remove habit from UI
	var newHabits []models.Habit
	for _, h := range m.data.Habits {
		if h.ID != habitID {
			newHabits = append(newHabits, h)
		}
	}
	m.data.Habits = newHabits

	m.habitList = newHabitList(m.data, m.width-4, m.height-10)
	m.message = "Deleting habit..."
	m.mode = QuestListView

	// Keep confirmHabit for rollback in case of error
	m.habitDeleteCmd = deleteHabitCmd(m.storage, habitID)

	return m, m.habitDeleteCmd
}

func (m Model) confirmDeleteJourney() (Model, tea.Cmd) {
	if m.confirmJourney == nil {
		m.mode = QuestListView
		return m, nil
	}

	journeyID := m.confirmJourney.ID

	// Optimistically remove journey from UI
	var newJourneys []models.Journey
	for _, j := range m.data.Journeys {
		if j.ID != journeyID {
			newJourneys = append(newJourneys, j)
		}
	}
	m.data.Journeys = newJourneys

	m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
	m.message = "Deleting journey..."
	m.mode = QuestListView

	// Keep confirmJourney for rollback in case of error
	m.journeyDeleteCmd = deleteJourneyCmd(m.storage, journeyID)

	return m, m.journeyDeleteCmd
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

func (m Model) editQuest(quest models.Quest) (Model, tea.Cmd) {
	m.editingQuest = &quest
	m.questFormData = &QuestForm{
		Title:      quest.Title,
		Note:       quest.Note,
		Difficulty: quest.Difficulty,
		JourneyID:  quest.JourneyID,
	}
	m.questForm = BuildQuestForm(m.questFormData, m.data.Journeys)
	m.mode = QuestEditFormView
	return m, m.questForm.Init()
}

func (m Model) editHabit(habit models.Habit) (Model, tea.Cmd) {
	m.editingHabit = &habit
	m.habitFormData = &HabitForm{
		Name:        habit.Name,
		CycleType:   habit.CycleType,
		CycleConfig: habit.CycleConfig,
	}
	m.habitForm = BuildHabitForm(m.habitFormData)
	m.mode = HabitEditFormView
	return m, m.habitForm.Init()
}

func (m Model) editJourney(journey models.Journey) (Model, tea.Cmd) {
	m.editingJourney = &journey
	m.journeyFormData = &JourneyForm{
		Name: journey.Name,
	}
	m.journeyForm = BuildJourneyForm(m.journeyFormData)
	m.mode = JourneyEditFormView
	return m, m.journeyForm.Init()
}

func (m Model) editEvent(event *models.Event) (Model, tea.Cmd) {
	m.editingEvent = event

	timeStr := ""
	if event.Time != nil {
		timeStr = *event.Time
	}
	endTimeStr := ""
	if event.EndTime != nil {
		endTimeStr = *event.EndTime
	}
	locationStr := ""
	if event.Location != nil {
		locationStr = *event.Location
	}
	descriptionStr := ""
	if event.Description != nil {
		descriptionStr = *event.Description
	}

	m.eventFormData = &EventForm{
		Title:       event.Title,
		Date:        event.Date.Format("2006-01-02"),
		Time:        timeStr,
		EndTime:     endTimeStr,
		Location:    locationStr,
		Description: descriptionStr,
	}
	m.eventForm = BuildEventForm(m.eventFormData)
	m.mode = EventEditFormView
	return m, m.eventForm.Init()
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

		tempID := -int(time.Now().UnixNano())
		optimisticQuest := &models.Quest{
			ID:         tempID,
			Title:      m.questFormData.Title,
			Note:       m.questFormData.Note,
			Done:       false,
			Difficulty: m.questFormData.Difficulty,
			AuthorID:   0,
			XPReward:   0,
			GoldReward: 0,
			JourneyID:  journeyID,
			Status:     "todo",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if journeyID != nil {
			found := false
			for i := range m.data.Journeys {
				if m.data.Journeys[i].ID == *journeyID {
					m.data.Journeys[i].Quests = append(m.data.Journeys[i].Quests, *optimisticQuest)
					found = true
					break
				}
			}
			if !found {
				message = "Journey not found"
				m.mode = returnMode
				m.message = message
				return m, nil
			}
		} else {
			found := false
			for i := range m.data.Journeys {
				if m.data.Journeys[i].ID == 0 {
					m.data.Journeys[i].Quests = append(m.data.Journeys[i].Quests, *optimisticQuest)
					found = true
					break
				}
			}
			if !found {
				myQuests := models.Journey{
					ID:     0,
					Name:   "My Quests",
					Quests: []models.Quest{*optimisticQuest},
				}
				m.data.Journeys = append([]models.Journey{myQuests}, m.data.Journeys...)
			}
		}

		m.questList = newQuestList(m.data, m.width-4, m.height-10)

		if m.selectedJourney != nil {
			returnMode = JourneyDetailView
			for _, j := range m.data.Journeys {
				if j.ID == m.selectedJourney.ID {
					m.selectedJourney = &j
					m.journeyQuestList = newJourneyQuestList(&j, m.width-4, m.height-10)
					break
				}
			}
		}

		message = fmt.Sprintf("✓ Quest created: %s", m.questFormData.Title)

		m.questCreateCmd = createQuestCmd(
			m.storage,
			tempID,
			m.questFormData.Title,
			m.questFormData.Note,
			m.questFormData.Difficulty,
			journeyID,
		)

	case JourneyFormView:
		if m.journeyFormData.Name == "" {
			message = "Journey name cannot be empty"
			m.mode = returnMode
			return m, nil
		}

		tempID := -int(time.Now().UnixNano())
		optimisticJourney := models.Journey{
			ID:        tempID,
			Name:      m.journeyFormData.Name,
			AuthorID:  0,
			Quests:    []models.Quest{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		m.data.Journeys = append(m.data.Journeys, optimisticJourney)

		m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)

		message = fmt.Sprintf("✓ Journey created: %s", m.journeyFormData.Name)
		m.currentSection = "journeys"

		if m.selectedJourney != nil {
			returnMode = JourneyDetailView
		}

		m.journeyCreateCmd = createJourneyCmd(m.storage, tempID, m.journeyFormData.Name)

	case HabitFormView:
		if m.habitFormData.Name == "" {
			message = "Habit name cannot be empty"
			m.mode = returnMode
			return m, nil
		}

		tempID := -int(time.Now().UnixNano())
		optimisticHabit := models.Habit{
			ID:            tempID,
			Name:          m.habitFormData.Name,
			AuthorID:      0,
			XPReward:      10,
			GoldReward:    10,
			CycleType:     m.habitFormData.CycleType,
			CycleConfig:   m.habitFormData.CycleConfig,
			Completed:     []string{},
			CurrentStreak: 0,
			MaxStreak:     0,
			StartDate:     time.Now(),
			EndDate:       nil,
			IsDueToday:    true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		m.data.Habits = append(m.data.Habits, optimisticHabit)

		m.habitList = newHabitList(m.data, m.width-4, m.height-10)

		message = fmt.Sprintf("✓ Habit created: %s", m.habitFormData.Name)

		m.habitCreateCmd = createHabitCmd(
			m.storage,
			tempID,
			m.habitFormData.Name,
			m.habitFormData.CycleType,
			m.habitFormData.CycleConfig,
		)

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

		tempID := -int(time.Now().UnixNano())

		eventDate, err := time.Parse("2006-01-02", m.eventFormData.Date)
		if err != nil {
			message = fmt.Sprintf("Invalid date format: %v", err)
			m.mode = returnMode
			m.currentSection = "calendar"
			m.message = message
			return m, nil
		}

		optimisticEvent := models.Event{
			ID:          tempID,
			Title:       m.eventFormData.Title,
			Date:        eventDate,
			EndDate:     nil,
			Time:        timePtr,
			EndTime:     endTimePtr,
			Location:    locationPtr,
			Description: descriptionPtr,
			AuthorID:    0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		m.data.Events = append(m.data.Events, optimisticEvent)
		m.calendar.SetEvents(m.data.Events)

		message = fmt.Sprintf("✓ Event created: %s", m.eventFormData.Title)

		m.currentSection = "calendar"
		returnMode = QuestListView

		m.eventCreateCmd = createEventCmd(m.storage, tempID, api.CreateEventRequest{
			Title:       m.eventFormData.Title,
			Date:        m.eventFormData.Date,
			Time:        timePtr,
			EndTime:     endTimePtr,
			Location:    locationPtr,
			Description: descriptionPtr,
		})

	case QuestEditFormView:
		if m.editingQuest == nil {
			message = "No quest being edited"
			m.mode = QuestListView
			m.message = message
			return m, nil
		}

		if m.questFormData.Title == "" {
			message = "Quest title cannot be empty"
			m.mode = QuestListView
			m.message = message
			return m, nil
		}

		// Optimistically update quest in data
		questID := m.editingQuest.ID
		for i := range m.data.Journeys {
			for j := range m.data.Journeys[i].Quests {
				if m.data.Journeys[i].Quests[j].ID == questID {
					m.data.Journeys[i].Quests[j].Title = m.questFormData.Title
					m.data.Journeys[i].Quests[j].Note = m.questFormData.Note
					m.data.Journeys[i].Quests[j].Difficulty = m.questFormData.Difficulty
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
			returnMode = JourneyDetailView
		}

		message = fmt.Sprintf("Updating quest...")

		// Keep editingQuest for rollback, fire async update
		m.questUpdateCmd = updateQuestCmd(m.storage, questID, api.UpdateQuestRequest{
			Title:      &m.questFormData.Title,
			Note:       &m.questFormData.Note,
			Difficulty: &m.questFormData.Difficulty,
		})

	case HabitEditFormView:
		if m.editingHabit == nil {
			message = "No habit being edited"
			m.mode = QuestListView
			m.message = message
			return m, nil
		}

		if m.habitFormData.Name == "" {
			message = "Habit name cannot be empty"
			m.mode = QuestListView
			m.message = message
			return m, nil
		}

		// Optimistically update habit in data
		habitID := m.editingHabit.ID
		for i := range m.data.Habits {
			if m.data.Habits[i].ID == habitID {
				m.data.Habits[i].Name = m.habitFormData.Name
				m.data.Habits[i].CycleType = m.habitFormData.CycleType
				m.data.Habits[i].CycleConfig = m.habitFormData.CycleConfig
				break
			}
		}

		m.habitList = newHabitList(m.data, m.width-4, m.height-10)
		message = fmt.Sprintf("Updating habit...")

		// Keep editingHabit for rollback, fire async update
		m.habitUpdateCmd = updateHabitCmd(m.storage, habitID, api.UpdateHabitRequest{
			Name:        &m.habitFormData.Name,
			CycleType:   &m.habitFormData.CycleType,
			CycleConfig: m.habitFormData.CycleConfig,
		})

	case JourneyEditFormView:
		if m.editingJourney == nil {
			message = "No journey being edited"
			m.mode = QuestListView
			m.message = message
			return m, nil
		}

		if m.journeyFormData.Name == "" {
			message = "Journey name cannot be empty"
			m.mode = QuestListView
			m.message = message
			return m, nil
		}

		// Optimistically update journey in data
		journeyID := m.editingJourney.ID
		for i := range m.data.Journeys {
			if m.data.Journeys[i].ID == journeyID {
				m.data.Journeys[i].Name = m.journeyFormData.Name
				break
			}
		}

		m.journeyList = newJourneyList(m.data, m.width-4, m.height-10)
		m.currentSection = "journeys"
		message = fmt.Sprintf("Updating journey...")

		// Keep editingJourney for rollback, fire async update
		m.journeyUpdateCmd = updateJourneyCmd(m.storage, journeyID, api.UpdateJourneyRequest{
			Name: &m.journeyFormData.Name,
		})

	case EventEditFormView:
		if m.editingEvent == nil {
			message = "No event being edited"
			m.mode = QuestListView
			m.currentSection = "calendar"
			m.message = message
			return m, nil
		}

		if m.eventFormData.Title == "" {
			message = "Event title cannot be empty"
			m.mode = QuestListView
			m.currentSection = "calendar"
			m.message = message
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

		event, err := m.storage.GetAPIClient().UpdateEvent(m.editingEvent.ID, api.UpdateEventRequest{
			Title:       &m.eventFormData.Title,
			Date:        &m.eventFormData.Date,
			Time:        timePtr,
			EndTime:     endTimePtr,
			Location:    locationPtr,
			Description: descriptionPtr,
		})

		if err != nil {
			message = fmt.Sprintf("Failed to update event: %v", err)
			m.mode = QuestListView
			m.currentSection = "calendar"
			m.message = message
			m.editingEvent = nil
			return m, nil
		}

		message = fmt.Sprintf("✓ Event updated: %s", event.Title)
		m.editingEvent = nil
		m.currentSection = "calendar"
		returnMode = QuestListView
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

	var cmd tea.Cmd
	if m.questCreateCmd != nil {
		cmd = m.questCreateCmd
		m.questCreateCmd = nil
	} else if m.journeyCreateCmd != nil {
		cmd = m.journeyCreateCmd
		m.journeyCreateCmd = nil
	} else if m.habitCreateCmd != nil {
		cmd = m.habitCreateCmd
		m.habitCreateCmd = nil
	} else if m.eventCreateCmd != nil {
		cmd = m.eventCreateCmd
		m.eventCreateCmd = nil
	} else if m.questUpdateCmd != nil {
		cmd = m.questUpdateCmd
		m.questUpdateCmd = nil
	} else if m.journeyUpdateCmd != nil {
		cmd = m.journeyUpdateCmd
		m.journeyUpdateCmd = nil
	} else if m.habitUpdateCmd != nil {
		cmd = m.habitUpdateCmd
		m.habitUpdateCmd = nil
	} else if m.questDeleteCmd != nil {
		cmd = m.questDeleteCmd
		m.questDeleteCmd = nil
	} else if m.journeyDeleteCmd != nil {
		cmd = m.journeyDeleteCmd
		m.journeyDeleteCmd = nil
	} else if m.habitDeleteCmd != nil {
		cmd = m.habitDeleteCmd
		m.habitDeleteCmd = nil
	}

	if cmd != nil {
		return m, tea.Batch(cmd, clearMessageAfter(1*time.Second))
	}

	return m, clearMessageAfter(1 * time.Second)
}
