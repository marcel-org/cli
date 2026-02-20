package ui

import (
	"fmt"
	"marcel-cli/models"
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

func (m Model) deleteQuest(quest models.Quest) Model {
	dialog, err := NewConfirmDialog(
		"Delete Quest?",
		fmt.Sprintf("Are you sure you want to delete: %s", quest.Title),
	)

	if err != nil {
		m.message = fmt.Sprintf("Error: %v", err)
		return m
	}

	if !dialog.Confirmed {
		m.message = "Deletion cancelled"
		return m
	}

	err = m.storage.GetAPIClient().DeleteQuest(quest.ID)
	if err != nil {
		m.message = fmt.Sprintf("Failed to delete quest: %v", err)
		return m
	}

	for i := range m.data.Journeys {
		newQuests := []models.Quest{}
		for _, q := range m.data.Journeys[i].Quests {
			if q.ID != quest.ID {
				newQuests = append(newQuests, q)
			}
		}
		m.data.Journeys[i].Quests = newQuests
	}

	m.questList = newQuestList(m.data, m.width-4, m.height-8)
	m.message = "Quest deleted successfully"

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

	m = m.refreshQuests()
	m.message = fmt.Sprintf("✓ Quest created: %s", quest.Title)

	return m
}

func (m Model) refreshQuests() Model {
	m.mode = LoadingView
	m.message = "Refreshing quests..."

	data, err := m.storage.Load()
	if err != nil {
		m.mode = ErrorView
		m.errorMessage = fmt.Sprintf("Failed to load quests: %v", err)
		return m
	}

	m.data = data
	m.questList = newQuestList(m.data, m.width-4, m.height-8)
	m.mode = QuestListView
	m.message = "✓ Quests refreshed!"

	return m
}
