package ui

import (
	"fmt"
	"strings"

	"marcel-cli/models"
	"marcel-cli/storage"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewMode int

const (
	QuestView ViewMode = iota
	HelpView
	ErrorView
)

type Model struct {
	storage       *storage.Storage
	data          *models.AppData
	mode          ViewMode
	questCursor   int
	message       string
	errorMessage  string
	width         int
	height        int
	allQuests     []models.Quest
	loading       bool
	lastKey       string
}

func NewModel() (*Model, error) {
	s, err := storage.New()
	if err != nil {
		return nil, err
	}

	if err := s.GetAPIClient().CheckAuth(); err != nil {
		return &Model{
			storage:      s,
			data:         &models.AppData{},
			mode:         ErrorView,
			errorMessage: fmt.Sprintf("Authentication failed: %v\n\nPlease set MARCEL_TOKEN environment variable or configure ~/.marcel.yml", err),
		}, nil
	}

	data, err := s.Load()
	if err != nil {
		return &Model{
			storage:      s,
			data:         &models.AppData{},
			mode:         ErrorView,
			errorMessage: fmt.Sprintf("Failed to load quests: %v", err),
		}, nil
	}

	allQuests := []models.Quest{}
	for _, journey := range data.Journeys {
		allQuests = append(allQuests, journey.Quests...)
	}

	return &Model{
		storage:     s,
		data:        data,
		mode:        QuestView,
		questCursor: 0,
		message:     "",
		allQuests:   allQuests,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		return m.handleKeypress(msg)
	}

	return m, nil
}

func (m Model) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case QuestView:
		return m.handleQuestViewKeys(msg)
	case HelpView:
		return m.handleHelpViewKeys(msg)
	case ErrorView:
		return m.handleErrorViewKeys(msg)
	}
	return m, nil
}

func (m Model) handleQuestViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.questCursor > 0 {
			m.questCursor--
		}
		m.lastKey = ""
	case "down", "j":
		if m.questCursor < len(m.allQuests)-1 {
			m.questCursor++
		}
		m.lastKey = ""
	case "g":
		if m.lastKey == "g" {
			m.questCursor = 0
			m.lastKey = ""
		} else {
			m.lastKey = "g"
		}
	case "G":
		if len(m.allQuests) > 0 {
			m.questCursor = len(m.allQuests) - 1
		}
		m.lastKey = ""
	case " ":
		m.lastKey = ""
		return m.toggleCurrentQuest(), nil
	case "r":
		m.lastKey = ""
		return m.refreshQuests(), nil
	case "?":
		m.lastKey = ""
		m.mode = HelpView
	default:
		m.lastKey = ""
	}
	return m, nil
}

func (m Model) handleHelpViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc", "?":
		m.mode = QuestView
	}
	return m, nil
}

func (m Model) handleErrorViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "r":
		return m.refreshQuests(), nil
	}
	return m, nil
}

func (m Model) toggleCurrentQuest() Model {
	if len(m.allQuests) == 0 || m.questCursor >= len(m.allQuests) {
		return m
	}

	quest := m.allQuests[m.questCursor]
	newDone := !quest.Done

	_, err := m.storage.GetAPIClient().ToggleQuest(quest.ID, newDone)
	if err != nil {
		m.message = fmt.Sprintf("Failed to toggle quest: %v", err)
		return m
	}

	m.allQuests[m.questCursor].Done = newDone

	for i := range m.data.Journeys {
		for j := range m.data.Journeys[i].Quests {
			if m.data.Journeys[i].Quests[j].ID == quest.ID {
				m.data.Journeys[i].Quests[j].Done = newDone
			}
		}
	}

	if newDone {
		m.message = "Quest completed!"
	} else {
		m.message = "Quest marked as incomplete"
	}

	return m
}

func (m Model) refreshQuests() Model {
	m.loading = true
	m.message = "Refreshing quests..."

	data, err := m.storage.Load()
	if err != nil {
		m.mode = ErrorView
		m.errorMessage = fmt.Sprintf("Failed to load quests: %v", err)
		m.loading = false
		return m
	}

	allQuests := []models.Quest{}
	for _, journey := range data.Journeys {
		allQuests = append(allQuests, journey.Quests...)
	}

	m.data = data
	m.allQuests = allQuests
	m.questCursor = 0
	m.message = "Quests refreshed!"
	m.loading = false
	m.mode = QuestView

	return m
}

func (m Model) View() string {
	switch m.mode {
	case QuestView:
		return m.renderQuestView()
	case HelpView:
		return m.renderHelpView()
	case ErrorView:
		return m.renderErrorView()
	}
	return ""
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#5856D6")).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5856D6")).
			Bold(true)

	completedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Strikethrough(true)

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF3B30")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#34C759")).
			Bold(true)
)

func (m Model) renderQuestView() string {
	title := titleStyle.Render("Marcel CLI - Your Quests")

	var content string
	if len(m.allQuests) == 0 {
		content = mutedStyle.Render("No quests found.\n\nCreate quests in the Marcel web app or use 'r' to refresh.")
	} else {
		var lines []string
		currentJourney := ""

		questIndex := 0
		for _, journey := range m.data.Journeys {
			if len(journey.Quests) == 0 {
				continue
			}

			if currentJourney != journey.Name {
				currentJourney = journey.Name
				lines = append(lines, "\n"+lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5856D6")).Render(journey.Name))
			}

			for _, quest := range journey.Quests {
				cursor := " "
				if questIndex == m.questCursor {
					cursor = ">"
				}

				checkbox := "[ ]"
				questText := quest.Title
				if quest.Done {
					checkbox = "[x]"
					questText = completedStyle.Render(questText)
				} else if questIndex == m.questCursor {
					questText = selectedStyle.Render(questText)
				}

				reward := mutedStyle.Render(fmt.Sprintf(" (+%d XP, +%d gold)", quest.XPReward, quest.GoldReward))
				line := fmt.Sprintf("%s %s %s%s", cursor, checkbox, questText, reward)
				lines = append(lines, line)

				questIndex++
			}
		}
		content = strings.Join(lines, "\n")
	}

	var statusLine string
	if m.message != "" {
		if strings.Contains(m.message, "completed") || strings.Contains(m.message, "refreshed") {
			statusLine = "\n" + successStyle.Render(m.message)
		} else {
			statusLine = "\n" + mutedStyle.Render(m.message)
		}
	}

	help := mutedStyle.Render("\n\nj/k or arrows navigate  • gg/G jump  • space toggle  • r refresh  • ? help  • q quit")

	return title + "\n" + content + statusLine + help
}

func (m Model) renderHelpView() string {
	title := titleStyle.Render("Help")

	help := `
Quest View:
  j/k, arrows Navigate quests
  gg          Jump to top
  G           Jump to bottom
  Space       Toggle quest completion
  r           Refresh quests from server
  ?           Show/hide help
  q, Ctrl+C   Quit

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

	footer := "\nPress any key to return..."

	return title + help + footer
}

func (m Model) renderErrorView() string {
	title := errorStyle.Render("Error")

	content := m.errorMessage

	help := mutedStyle.Render("\n\nr - retry  • q - quit")

	return title + "\n\n" + content + help
}
