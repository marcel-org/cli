package ui

import (
	"fmt"
	"marcel-cli/models"
	"time"
)

func (m Model) toggleQuest(quest models.Quest) Model {
	newDone := !quest.Done

	_, err := m.storage.GetAPIClient().ToggleQuest(quest.ID, newDone)
	if err != nil {
		m.message = fmt.Sprintf("Failed to toggle quest: %v", err)
		return m
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

	return m
}

func (m Model) showDeleteConfirm(quest models.Quest) Model {
	m.mode = ConfirmDeleteView
	m.confirmQuest = &quest
	m.confirmSelected = false
	return m
}

func (m Model) confirmDeleteQuest() Model {
	if m.confirmQuest == nil {
		m.mode = QuestListView
		return m
	}

	err := m.storage.GetAPIClient().DeleteQuest(m.confirmQuest.ID)
	if err != nil {
		m.message = fmt.Sprintf("Failed to delete quest: %v", err)
		m.mode = QuestListView
		m.confirmQuest = nil
		return m
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

	return m
}

func (m Model) confirmDeleteHabit() Model {
	if m.confirmHabit == nil {
		m.mode = QuestListView
		return m
	}

	err := m.storage.GetAPIClient().DeleteHabit(m.confirmHabit.ID)
	if err != nil {
		m.message = fmt.Sprintf("Failed to delete habit: %v", err)
		m.mode = QuestListView
		m.confirmHabit = nil
		return m
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

	return m
}

func (m Model) confirmDeleteJourney() Model {
	if m.confirmJourney == nil {
		m.mode = QuestListView
		return m
	}

	err := m.storage.GetAPIClient().DeleteJourney(m.confirmJourney.ID)
	if err != nil {
		m.message = fmt.Sprintf("Failed to delete journey: %v", err)
		m.mode = QuestListView
		m.confirmJourney = nil
		return m
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

	return m
}

func (m Model) cancelDelete() Model {
	m.mode = QuestListView
	m.confirmQuest = nil
	m.confirmHabit = nil
	m.confirmJourney = nil
	m.message = "Deletion cancelled"
	return m
}

func (m Model) createNewQuest() Model {
	form, err := NewQuestForm(m.data.Journeys)
	if err != nil {
		m.message = fmt.Sprintf("Error: %v", err)
		return m
	}

	var journeyID *int
	if form.JourneyID != nil && *form.JourneyID != 0 {
		journeyID = form.JourneyID
	}

	quest, err := m.storage.GetAPIClient().CreateQuest(
		form.Title,
		form.Note,
		form.Difficulty,
		journeyID,
	)

	if err != nil {
		m.message = fmt.Sprintf("Failed to create quest: %v", err)
		return m
	}

	m = m.refreshData()
	m.message = fmt.Sprintf("✓ Quest created: %s", quest.Title)

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
	m.mode = QuestListView
	m.message = "✓ Data refreshed!"

	return m
}

func (m Model) toggleHabit(habit models.Habit) Model {
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
		m.message = fmt.Sprintf("Failed to toggle habit: %v", err)
		return m
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

	return m
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
