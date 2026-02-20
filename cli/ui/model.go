package ui

import (
	"fmt"
	"marcel-cli/models"
	"marcel-cli/storage"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewMode int

const (
	QuestListView ViewMode = iota
	LoadingView
	ErrorView
	HelpView
	ConfirmDeleteView
)

type Model struct {
	storage         *storage.Storage
	data            *models.AppData
	mode            ViewMode
	questList       list.Model
	spinner         spinner.Model
	message         string
	errorMessage    string
	width           int
	height          int
	ready           bool
	needsRedraw     bool
	confirmQuest    *models.Quest
	confirmSelected bool
}

func NewModel() (*Model, error) {
	s, err := storage.New()
	if err != nil {
		return nil, err
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = SpinnerStyle

	m := &Model{
		storage: s,
		mode:    LoadingView,
		spinner: sp,
		data:    &models.AppData{},
	}

	if err := s.GetAPIClient().CheckAuth(); err != nil {
		m.mode = ErrorView
		m.errorMessage = fmt.Sprintf("Authentication failed: %v\n\nPlease set MARCEL_TOKEN environment variable or configure ~/.marcel.yml", err)
		return m, nil
	}

	data, err := s.Load()
	if err != nil {
		m.mode = ErrorView
		m.errorMessage = fmt.Sprintf("Failed to load quests: %v", err)
		return m, nil
	}

	m.data = data
	m.mode = QuestListView

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tea.EnterAltScreen)
}
