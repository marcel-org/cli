package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.ready && m.mode != ErrorView {
		return m.renderLoadingView()
	}

	switch m.mode {
	case QuestListView:
		return m.renderQuestListView()
	case JourneyDetailView:
		return m.renderJourneyDetailView()
	case LoadingView:
		return m.renderLoadingView()
	case ErrorView:
		return m.renderErrorView()
	case HelpView:
		return m.renderHelpView()
	case ConfirmDeleteView:
		return m.renderConfirmDeleteView()
	case QuestFormView:
		return m.renderFormView(m.questForm)
	case JourneyFormView:
		return m.renderFormView(m.journeyForm)
	case HabitFormView:
		return m.renderFormView(m.habitForm)
	case EventFormView:
		return m.renderFormView(m.eventForm)
	}

	return ""
}

func (m Model) renderQuestListView() string {
	var headerText string
	var content string

	tabs := renderTabs(m.currentSection, m.width)

	switch m.currentSection {
	case "quests":
		headerText = "Marcel - Quests"
		content = m.questList.View()
	case "habits":
		headerText = "Marcel - Habits"
		content = m.habitList.View()
	case "journeys":
		headerText = "Marcel - Journeys"
		content = m.journeyList.View()
	case "calendar":
		headerText = "Marcel - Calendar"
		content = m.calendar.View()
	default:
		headerText = "Marcel - Quests"
		content = m.questList.View()
	}

	header := HeaderStyle.Width(m.width).Render(headerText)

	statusBar := ""
	if m.message != "" {
		var msgStyle lipgloss.Style
		if strings.Contains(m.message, "✓") {
			msgStyle = SuccessStyle
		} else if strings.Contains(m.message, "Failed") || strings.Contains(m.message, "Error") {
			msgStyle = ErrorStyle
		} else {
			msgStyle = MutedStyle
		}
		statusBar = "\n" + StatusBarStyle.Width(m.width).Render(msgStyle.Render(m.message))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tabs,
		content,
		statusBar,
	)
}

func renderTabs(currentSection string, width int) string {
	questStyle := lipgloss.NewStyle().Foreground(lightGray).Padding(0, 2)
	habitStyle := lipgloss.NewStyle().Foreground(lightGray).Padding(0, 2)
	journeyStyle := lipgloss.NewStyle().Foreground(lightGray).Padding(0, 2)
	calendarStyle := lipgloss.NewStyle().Foreground(lightGray).Padding(0, 2)

	switch currentSection {
	case "quests":
		questStyle = questStyle.Foreground(brandOrange).Bold(true).Background(darkGray)
	case "habits":
		habitStyle = habitStyle.Foreground(brandOrange).Bold(true).Background(darkGray)
	case "journeys":
		journeyStyle = journeyStyle.Foreground(brandOrange).Bold(true).Background(darkGray)
	case "calendar":
		calendarStyle = calendarStyle.Foreground(brandOrange).Bold(true).Background(darkGray)
	}

	tabs := lipgloss.JoinHorizontal(
		lipgloss.Left,
		questStyle.Render("Quests"),
		habitStyle.Render("Habits"),
		journeyStyle.Render("Journeys"),
		calendarStyle.Render("Calendar"),
	)

	return lipgloss.NewStyle().Width(width).Render(tabs)
}

func (m Model) renderLoadingView() string {
	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.spinner.View(),
		" Loading quests...",
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		SpinnerStyle.Render(content),
	)
}

func (m Model) renderErrorView() string {
	title := ErrorStyle.Render("⚠ Error")

	content := lipgloss.NewStyle().
		Width(m.width - 8).
		Padding(1, 2).
		Render(m.errorMessage)

	help := HelpStyle.Render("r - retry  •  q - quit")

	box := BoxStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			content,
			"",
			help,
		),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func (m Model) renderHelpView() string {
	title := TitleStyle.Render("Help & Keyboard Shortcuts")

	helpContent := `
Navigation:
  ↑/k, ↓/j     Navigate quests
  /            Filter quests
  gg           Jump to top
  G            Jump to bottom

Actions:
  space/enter  Toggle quest completion
  n            Create new quest
  d            Delete quest
  r            Refresh quests from server

Other:
  ?            Show/hide help
  q, Ctrl+C    Quit

Configuration:
  Environment variable:
    MARCEL_TOKEN           - Your Marcel CLI token

  Config file (~/.marcel.yml):
    week_start_day: sunday

Authentication:
  1. Go to Marcel web app settings
  2. Generate a Marcel CLI token
  3. Copy the token
  4. Set MARCEL_TOKEN environment variable

Note: API endpoint is hardcoded to https://api.marcel.my
`

	footer := HelpStyle.Render("\nPress ? or Esc to return")

	box := BoxStyle.Width(m.width - 8).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			helpContent,
			footer,
		),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func (m Model) renderConfirmDeleteView() string {
	var title, itemName string

	if m.confirmQuest != nil {
		title = "Delete Quest?"
		itemName = m.confirmQuest.Title
	} else if m.confirmHabit != nil {
		title = "Delete Habit?"
		itemName = m.confirmHabit.Name
	} else if m.confirmJourney != nil {
		title = "Delete Journey?"
		itemName = m.confirmJourney.Name
	} else if m.confirmEvent != nil {
		title = "Delete Event?"
		itemName = m.confirmEvent.Title
	} else {
		return ""
	}

	titleStyled := lipgloss.NewStyle().
		Foreground(red).
		Bold(true).
		Render(title)

	itemNameStyled := lipgloss.NewStyle().
		Foreground(white).
		Render(itemName)

	noStyle := lipgloss.NewStyle().
		Foreground(white).
		Padding(0, 2)

	yesStyle := lipgloss.NewStyle().
		Foreground(white).
		Padding(0, 2)

	if !m.confirmSelected {
		noStyle = noStyle.Background(brandOrange).Bold(true)
	} else {
		yesStyle = yesStyle.Background(red).Bold(true)
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Left,
		noStyle.Render("No"),
		"  ",
		yesStyle.Render("Yes"),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyled,
		"",
		itemNameStyled,
		"",
		buttons,
		"",
		MutedStyle.Render("←/h: No  →/l: Yes  Enter: Confirm  Esc: Cancel"),
	)

	box := BoxStyle.Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func (m Model) renderFormView(form *huh.Form) string {
	if form == nil {
		return ""
	}

	formView := form.View()

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		formView,
	)
}

func (m Model) renderJourneyDetailView() string {
	if m.selectedJourney == nil {
		return ""
	}

	headerText := fmt.Sprintf("Marcel - %s", m.selectedJourney.Name)
	header := HeaderStyle.Width(m.width).Render(headerText)

	content := m.journeyQuestList.View()

	statusBar := ""
	if m.message != "" {
		var msgStyle lipgloss.Style
		if strings.Contains(m.message, "✓") {
			msgStyle = SuccessStyle
		} else if strings.Contains(m.message, "Failed") || strings.Contains(m.message, "Error") {
			msgStyle = ErrorStyle
		} else {
			msgStyle = MutedStyle
		}
		statusBar = "\n" + StatusBarStyle.Width(m.width).Render(msgStyle.Render(m.message))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		statusBar,
	)
}
