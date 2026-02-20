package components

import (
	"fmt"
	"marcel-cli/models"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Calendar struct {
	currentDate    time.Time
	selectedDate   time.Time
	events         []models.Event
	selectedEvent  int
	width          int
	height         int
	weekStartDay   string
}

func NewCalendar() *Calendar {
	now := time.Now()
	return &Calendar{
		currentDate:   now,
		selectedDate:  now,
		events:        []models.Event{},
		selectedEvent: 0,
		width:         80,
		height:        24,
		weekStartDay:  "sunday",
	}
}

func (c *Calendar) SetSize(width, height int) {
	c.width = width
	c.height = height
}

func (c *Calendar) SetWeekStartDay(day string) {
	c.weekStartDay = day
}

func (c *Calendar) SetEvents(events []models.Event) {
	c.events = events
}

func (c *Calendar) GetSelectedDate() time.Time {
	return c.selectedDate
}

func (c *Calendar) GetSelectedEvent() *models.Event {
	eventsOnDate := c.getEventsForDate(c.selectedDate)
	if len(eventsOnDate) > 0 && c.selectedEvent < len(eventsOnDate) {
		return &eventsOnDate[c.selectedEvent]
	}
	return nil
}

func (c *Calendar) NavigateLeft() {
	c.selectedDate = c.selectedDate.AddDate(0, 0, -1)
	c.selectedEvent = 0
	c.updateCurrentDateIfNeeded()
}

func (c *Calendar) NavigateRight() {
	c.selectedDate = c.selectedDate.AddDate(0, 0, 1)
	c.selectedEvent = 0
	c.updateCurrentDateIfNeeded()
}

func (c *Calendar) NavigateUp() {
	c.selectedDate = c.selectedDate.AddDate(0, 0, -7)
	c.selectedEvent = 0
	c.updateCurrentDateIfNeeded()
}

func (c *Calendar) NavigateDown() {
	c.selectedDate = c.selectedDate.AddDate(0, 0, 7)
	c.selectedEvent = 0
	c.updateCurrentDateIfNeeded()
}

func (c *Calendar) NavigateNextMonth() {
	c.currentDate = c.currentDate.AddDate(0, 1, 0)
	c.selectedDate = c.selectedDate.AddDate(0, 1, 0)
}

func (c *Calendar) NavigatePrevMonth() {
	c.currentDate = c.currentDate.AddDate(0, -1, 0)
	c.selectedDate = c.selectedDate.AddDate(0, -1, 0)
}

func (c *Calendar) GoToToday() {
	c.currentDate = time.Now()
	c.selectedDate = time.Now()
	c.selectedEvent = 0
}

func (c *Calendar) NextEvent() {
	eventsOnDate := c.getEventsForDate(c.selectedDate)
	if len(eventsOnDate) > 0 {
		c.selectedEvent = (c.selectedEvent + 1) % len(eventsOnDate)
	}
}

func (c *Calendar) PrevEvent() {
	eventsOnDate := c.getEventsForDate(c.selectedDate)
	if len(eventsOnDate) > 0 {
		c.selectedEvent = (c.selectedEvent - 1 + len(eventsOnDate)) % len(eventsOnDate)
	}
}

func (c *Calendar) updateCurrentDateIfNeeded() {
	if c.selectedDate.Month() != c.currentDate.Month() ||
		c.selectedDate.Year() != c.currentDate.Year() {
		c.currentDate = c.selectedDate
	}
}

func (c *Calendar) View() string {
	var sb strings.Builder

	header := c.renderHeader()
	sb.WriteString(header)
	sb.WriteString("\n\n")

	weekdaysSunday := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	weekdaysMonday := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

	var weekdays []string
	if c.weekStartDay == "monday" {
		weekdays = weekdaysMonday
	} else {
		weekdays = weekdaysSunday
	}

	for i, day := range weekdays {
		if i > 0 {
			sb.WriteString(" ")
		}
		style := lipgloss.NewStyle().Width(8).Align(lipgloss.Center).Foreground(lipgloss.Color("#6c7086"))
		sb.WriteString(style.Render(day))
	}
	sb.WriteString("\n")

	firstDay := time.Date(c.currentDate.Year(), c.currentDate.Month(), 1, 0, 0, 0, 0, c.currentDate.Location())

	var weekStartOffset int
	if c.weekStartDay == "monday" {
		weekStartOffset = int(firstDay.Weekday()) - 1
		if weekStartOffset < 0 {
			weekStartOffset = 6
		}
	} else {
		weekStartOffset = int(firstDay.Weekday())
	}

	startOfWeek := firstDay.AddDate(0, 0, -weekStartOffset)

	for week := 0; week < 6; week++ {
		weekHasCurrentMonth := false
		for day := 0; day < 7; day++ {
			date := startOfWeek.AddDate(0, 0, week*7+day)
			if date.Month() == c.currentDate.Month() {
				weekHasCurrentMonth = true
			}
			cellContent := c.renderMonthCell(date)
			sb.WriteString(cellContent)
			if day < 6 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")

		if !weekHasCurrentMonth && week > 0 {
			break
		}
	}

	sb.WriteString("\n")
	sb.WriteString(c.renderSelectedDateEvents())

	return sb.String()
}

func (c *Calendar) renderHeader() string {
	title := c.currentDate.Format("January 2006")

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(c.width - 4).
		Foreground(lipgloss.Color("#fab387"))

	return headerStyle.Render(title)
}

func (c *Calendar) renderMonthCell(date time.Time) string {
	dayStr := fmt.Sprintf("%2d", date.Day())

	events := c.getEventsForDate(date)
	eventIndicator := ""
	if len(events) > 0 {
		eventIndicator = fmt.Sprintf("\n●%d", len(events))
	}

	cellContent := dayStr + eventIndicator

	style := lipgloss.NewStyle().Width(8).Align(lipgloss.Center)

	if c.isSameDate(date, c.selectedDate) {
		style = style.Background(lipgloss.Color("#45475a")).Foreground(lipgloss.Color("#fab387")).Bold(true)
	} else if c.isToday(date) {
		style = style.Foreground(lipgloss.Color("#fab387")).Bold(true)
	} else if date.Month() != c.currentDate.Month() {
		style = style.Foreground(lipgloss.Color("#45475a"))
	}

	return style.Render(cellContent)
}

func (c *Calendar) renderSelectedDateEvents() string {
	var sb strings.Builder

	dateStr := c.selectedDate.Format("Monday, January 2, 2006")
	if c.isToday(c.selectedDate) {
		dateStr += " (Today)"
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cdd6f4"))
	sb.WriteString(headerStyle.Render(dateStr))
	sb.WriteString("\n\n")

	events := c.getEventsForDate(c.selectedDate)
	if len(events) == 0 {
		mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Italic(true)
		sb.WriteString(mutedStyle.Render("No events scheduled"))
	} else {
		for i, event := range events {
			eventStyle := lipgloss.NewStyle().Padding(0, 1)

			if i == c.selectedEvent {
				eventStyle = eventStyle.Background(lipgloss.Color("#fab387")).Foreground(lipgloss.Color("#1e1e2e"))
			}

			timeStr := ""
			if event.Time != nil && *event.Time != "" {
				timeStr = *event.Time
				if event.EndTime != nil && *event.EndTime != "" {
					timeStr += " - " + *event.EndTime
				}
			} else {
				timeStr = "All Day"
			}

			eventText := fmt.Sprintf("[%s] %s", timeStr, event.Title)
			if event.Location != nil && *event.Location != "" {
				eventText += fmt.Sprintf(" @ %s", *event.Location)
			}

			sb.WriteString(eventStyle.Render(eventText))
			sb.WriteString("\n")
		}

		if len(events) > 1 {
			sb.WriteString("\n")
			helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
			sb.WriteString(helpStyle.Render("Tab/Shift+Tab to cycle events"))
		}
	}

	return sb.String()
}

func (c *Calendar) getEventsForDate(date time.Time) []models.Event {
	var dayEvents []models.Event

	for _, event := range c.events {
		if c.isSameDate(event.Date, date) {
			dayEvents = append(dayEvents, event)
		}
	}

	return dayEvents
}

func (c *Calendar) isToday(date time.Time) bool {
	now := time.Now()
	return c.isSameDate(date, now)
}

func (c *Calendar) isSameDate(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
