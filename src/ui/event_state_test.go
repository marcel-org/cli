package ui

import (
	"marcel-cli/models"
	"marcel-cli/ui/components"
	"testing"
)

func TestUpsertEventUpdatesCalendarDisplayData(t *testing.T) {
	calendar := components.NewCalendar()
	selectedDate := calendar.GetSelectedDate()
	initialEvent := models.Event{
		ID:    42,
		Title: "Old title",
		Date:  selectedDate,
	}

	model := Model{
		data: &models.AppData{
			Events: []models.Event{initialEvent},
		},
		calendar: calendar,
	}

	model.calendar.SetEvents(model.data.Events)
	model.calendar.FocusEventList()

	updatedEvent := initialEvent
	updatedEvent.Title = "New title"

	model.upsertEvent(initialEvent.ID, updatedEvent)

	selectedEvent := model.calendar.GetSelectedEvent()
	if selectedEvent == nil {
		t.Fatal("expected calendar to keep a selected event")
	}

	if selectedEvent.Title != "New title" {
		t.Fatalf("expected calendar to show updated title, got %q", selectedEvent.Title)
	}

	if model.data.Events[0].Title != "New title" {
		t.Fatalf("expected model data to update event title, got %q", model.data.Events[0].Title)
	}
}
