package ui

import "marcel-cli/models"

func (m *Model) syncCalendarEvents() {
	if m.calendar != nil {
		m.calendar.SetEvents(m.data.Events)
	}
}

func (m *Model) upsertEvent(matchID int, event models.Event) {
	for i := range m.data.Events {
		if m.data.Events[i].ID == matchID {
			m.data.Events[i] = event
			m.syncCalendarEvents()
			return
		}
	}

	m.data.Events = append(m.data.Events, event)
	m.syncCalendarEvents()
}
