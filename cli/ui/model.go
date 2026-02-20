package ui

import (
	"fmt"
	"marcel-cli/models"
	"marcel-cli/storage"
	"marcel-cli/ui/components"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type ViewMode int

const (
	QuestListView ViewMode = iota
	LoadingView
	ErrorView
	HelpView
	ConfirmDeleteView
	JourneyDetailView
	QuestFormView
	JourneyFormView
	HabitFormView
	EventFormView
)

type clearMessageMsg struct{}

func clearMessageAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return clearMessageMsg{}
	})
}

type Model struct {
	storage           *storage.Storage
	data              *models.AppData
	mode              ViewMode
	currentSection    string
	questList         list.Model
	habitList         list.Model
	journeyList       list.Model
	journeyQuestList  list.Model
	calendar          *components.Calendar
	spinner           spinner.Model
	message           string
	errorMessage      string
	width             int
	height            int
	ready             bool
	needsRedraw       bool
	confirmQuest      *models.Quest
	confirmHabit      *models.Habit
	confirmJourney    *models.Journey
	confirmEvent      *models.Event
	confirmSelected   bool
	selectedJourney   *models.Journey
	questForm         *huh.Form
	journeyForm       *huh.Form
	habitForm         *huh.Form
	eventForm         *huh.Form
	questFormData     *QuestForm
	journeyFormData   *JourneyForm
	habitFormData     *HabitForm
	eventFormData     *EventForm
}

func NewModel() (*Model, error) {
	s, err := storage.New()
	if err != nil {
		return nil, err
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = SpinnerStyle

	cal := components.NewCalendar()
	cal.SetWeekStartDay(s.GetConfig().WeekStartDay)

	m := &Model{
		storage:        s,
		mode:           LoadingView,
		spinner:        sp,
		data:           &models.AppData{},
		currentSection: "quests",
		calendar:       cal,
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
	m.currentSection = data.CurrentSection

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tea.EnterAltScreen)
}
