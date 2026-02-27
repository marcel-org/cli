package ui

import (
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
	QuestEditFormView
	JourneyEditFormView
	HabitEditFormView
	EventEditFormView
)

type clearMessageMsg struct{}

type clearSyncStatusMsg struct{}

type dataLoadedMsg struct {
	data *models.AppData
	err  error
}

type backgroundSyncMsg struct {
	data *models.AppData
	err  error
}

type authCheckMsg struct {
	err error
}

func loadDataCmd(s *storage.Storage) tea.Cmd {
	return func() tea.Msg {
		data, err := s.LoadFromCache()
		return dataLoadedMsg{data: data, err: err}
	}
}

func checkAuthCmd(s *storage.Storage) tea.Cmd {
	return func() tea.Msg {
		err := s.GetAPIClient().CheckAuth()
		return authCheckMsg{err: err}
	}
}

func loadFromAPICmd(s *storage.Storage) tea.Cmd {
	return func() tea.Msg {
		data, err := s.Load()
		return dataLoadedMsg{data: data, err: err}
	}
}

func backgroundSyncCmd(s *storage.Storage) tea.Cmd {
	return func() tea.Msg {
		data, err := s.LoadAll()
		return backgroundSyncMsg{data: data, err: err}
	}
}

func clearMessageAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return clearMessageMsg{}
	})
}

func clearSyncStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return clearSyncStatusMsg{}
	})
}

type SyncStatus int

const (
	SyncStatusNone SyncStatus = iota
	SyncStatusSyncing
	SyncStatusSynced
	SyncStatusError
)

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
	editingQuest      *models.Quest
	editingHabit      *models.Habit
	editingJourney    *models.Journey
	editingEvent      *models.Event
	syncStatus        SyncStatus
	syncSpinner       spinner.Model
}

func NewModel() (*Model, error) {
	s, err := storage.New()
	if err != nil {
		return nil, err
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = SpinnerStyle

	syncSp := spinner.New()
	syncSp.Spinner = spinner.Dot
	syncSp.Style = SpinnerStyle

	cal := components.NewCalendar()
	cal.SetWeekStartDay(s.GetConfig().WeekStartDay)

	cachedData, cacheErr := s.LoadFromCache()

	mode := QuestListView
	data := &models.AppData{}
	if cacheErr == nil && cachedData != nil {
		data = cachedData
		mode = QuestListView
	} else {
		mode = LoadingView
	}

	m := &Model{
		storage:        s,
		mode:           mode,
		spinner:        sp,
		syncSpinner:    syncSp,
		data:           data,
		currentSection: data.CurrentSection,
		calendar:       cal,
		syncStatus:     SyncStatusNone,
	}

	if data.CurrentSection == "" {
		m.currentSection = "quests"
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.syncSpinner.Tick}

	if m.mode == LoadingView {
		cmds = append(cmds, m.spinner.Tick, loadFromAPICmd(m.storage))
	} else {
		cmds = append(cmds, checkAuthCmd(m.storage))
	}

	return tea.Batch(cmds...)
}
