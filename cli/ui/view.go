package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.ready && m.mode != ErrorView {
		return m.renderLoadingView()
	}

	switch m.mode {
	case QuestListView:
		return m.renderQuestListView()
	case LoadingView:
		return m.renderLoadingView()
	case ErrorView:
		return m.renderErrorView()
	case HelpView:
		return m.renderHelpView()
	}

	return ""
}

func (m Model) renderQuestListView() string {
	header := HeaderStyle.Width(m.width).Render("Marcel CLI - Your Quests")

	content := m.questList.View()

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

	help := "\n" + questListHelp()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		statusBar,
		help,
	)
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
  Environment variables (highest priority):
    MARCEL_TOKEN           - Your Marcel CLI token
    MARCEL_API_ENDPOINT    - API endpoint URL

  Config file (~/.marcel.yml):
    api_endpoint: http://localhost:3000
    auth_token: marcel_your_token_here

Authentication:
  1. Go to Marcel web app settings
  2. Generate a Marcel CLI token
  3. Copy the token
  4. Set MARCEL_TOKEN environment variable or add to config file
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
