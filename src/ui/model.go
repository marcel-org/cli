package ui

import (
	"marcel-cli/api"
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
	SpaceDetailView
	QuestFormView
	JourneyFormView
	HabitFormView
	EventFormView
	SpaceFormView
	QuestEditFormView
	JourneyEditFormView
	HabitEditFormView
	EventEditFormView
	SpaceEditFormView
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

type questCreatedMsg struct {
	tempID int
	quest  *models.Quest
	err    error
}

type journeyCreatedMsg struct {
	tempID  int
	journey *models.Journey
	err     error
}

type habitCreatedMsg struct {
	tempID int
	habit  *models.Habit
	err    error
}

type eventCreatedMsg struct {
	tempID int
	event  *models.Event
	err    error
}

type questDeletedMsg struct {
	questID int
	err     error
}

type journeyDeletedMsg struct {
	journeyID int
	err       error
}

type habitDeletedMsg struct {
	habitID int
	err     error
}

type questUpdatedMsg struct {
	questID int
	quest   *models.Quest
	err     error
}

type journeyUpdatedMsg struct {
	journeyID int
	journey   *models.Journey
	err       error
}

type habitUpdatedMsg struct {
	habitID int
	habit   *models.Habit
	err     error
}

type spaceCreatedMsg struct {
	tempID int
	space  *models.Space
	err    error
}

type spaceUpdatedMsg struct {
	spaceID int
	space   *models.Space
	err     error
}

func createQuestCmd(s *storage.Storage, tempID int, title, note, difficulty string, journeyID, spaceID *int) tea.Cmd {
	return func() tea.Msg {
		quest, err := s.GetAPIClient().CreateQuest(title, note, difficulty, journeyID, spaceID)
		return questCreatedMsg{tempID: tempID, quest: quest, err: err}
	}
}

func createJourneyCmd(s *storage.Storage, tempID int, name string, spaceID *int) tea.Cmd {
	return func() tea.Msg {
		journey, err := s.GetAPIClient().CreateJourney(name, spaceID)
		return journeyCreatedMsg{tempID: tempID, journey: journey, err: err}
	}
}

func createHabitCmd(s *storage.Storage, tempID int, name, cycleType string, cycleConfig any) tea.Cmd {
	return func() tea.Msg {
		habit, err := s.GetAPIClient().CreateHabit(name, cycleType, cycleConfig)
		return habitCreatedMsg{tempID: tempID, habit: habit, err: err}
	}
}

func createEventCmd(s *storage.Storage, tempID int, req api.CreateEventRequest) tea.Cmd {
	return func() tea.Msg {
		event, err := s.GetAPIClient().CreateEvent(req)
		return eventCreatedMsg{tempID: tempID, event: event, err: err}
	}
}

func deleteQuestCmd(s *storage.Storage, questID int) tea.Cmd {
	return func() tea.Msg {
		err := s.GetAPIClient().DeleteQuest(questID)
		return questDeletedMsg{questID: questID, err: err}
	}
}

func deleteJourneyCmd(s *storage.Storage, journeyID int) tea.Cmd {
	return func() tea.Msg {
		err := s.GetAPIClient().DeleteJourney(journeyID)
		return journeyDeletedMsg{journeyID: journeyID, err: err}
	}
}

func deleteHabitCmd(s *storage.Storage, habitID int) tea.Cmd {
	return func() tea.Msg {
		err := s.GetAPIClient().DeleteHabit(habitID)
		return habitDeletedMsg{habitID: habitID, err: err}
	}
}

func updateQuestCmd(s *storage.Storage, questID int, updates api.UpdateQuestRequest) tea.Cmd {
	return func() tea.Msg {
		quest, err := s.GetAPIClient().UpdateQuest(questID, updates)
		return questUpdatedMsg{questID: questID, quest: quest, err: err}
	}
}

func updateJourneyCmd(s *storage.Storage, journeyID int, updates api.UpdateJourneyRequest) tea.Cmd {
	return func() tea.Msg {
		journey, err := s.GetAPIClient().UpdateJourney(journeyID, updates)
		return journeyUpdatedMsg{journeyID: journeyID, journey: journey, err: err}
	}
}

func updateHabitCmd(s *storage.Storage, habitID int, updates api.UpdateHabitRequest) tea.Cmd {
	return func() tea.Msg {
		habit, err := s.GetAPIClient().UpdateHabit(habitID, updates)
		return habitUpdatedMsg{habitID: habitID, habit: habit, err: err}
	}
}

func createSpaceCmd(s *storage.Storage, tempID int, name string) tea.Cmd {
	return func() tea.Msg {
		space, err := s.GetAPIClient().CreateSpace(name)
		return spaceCreatedMsg{tempID: tempID, space: space, err: err}
	}
}

func updateSpaceCmd(s *storage.Storage, spaceID int, updates api.UpdateSpaceRequest) tea.Cmd {
	return func() tea.Msg {
		space, err := s.GetAPIClient().UpdateSpace(spaceID, updates)
		return spaceUpdatedMsg{spaceID: spaceID, space: space, err: err}
	}
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
	storage          *storage.Storage
	data             *models.AppData
	mode             ViewMode
	currentSection   string
	spaceSection     string
	questList        list.Model
	habitList        list.Model
	journeyList      list.Model
	spaceList        list.Model
	journeyQuestList list.Model
	spaceQuestList   list.Model
	spaceJourneyList list.Model
	spaceMemberList  list.Model
	calendar         *components.Calendar
	spinner          spinner.Model
	message          string
	errorMessage     string
	width            int
	height           int
	ready            bool
	needsRedraw      bool
	confirmQuest     *models.Quest
	confirmHabit     *models.Habit
	confirmJourney   *models.Journey
	confirmEvent     *models.Event
	confirmSelected  bool
	selectedJourney  *models.Journey
	selectedSpace    *models.Space
	questForm        *huh.Form
	journeyForm      *huh.Form
	habitForm        *huh.Form
	eventForm        *huh.Form
	spaceForm        *huh.Form
	questFormData    *QuestForm
	journeyFormData  *JourneyForm
	habitFormData    *HabitForm
	eventFormData    *EventForm
	spaceFormData    *SpaceForm
	editingQuest     *models.Quest
	editingHabit     *models.Habit
	editingJourney   *models.Journey
	editingEvent     *models.Event
	editingSpace     *models.Space
	syncStatus       SyncStatus
	syncSpinner      spinner.Model
	questCreateCmd   tea.Cmd
	journeyCreateCmd tea.Cmd
	habitCreateCmd   tea.Cmd
	eventCreateCmd   tea.Cmd
	spaceCreateCmd   tea.Cmd
	questDeleteCmd   tea.Cmd
	journeyDeleteCmd tea.Cmd
	habitDeleteCmd   tea.Cmd
	questUpdateCmd   tea.Cmd
	journeyUpdateCmd tea.Cmd
	habitUpdateCmd   tea.Cmd
	spaceUpdateCmd   tea.Cmd
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
		spaceSection:   "quests",
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
