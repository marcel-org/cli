package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"marcel-cli/api"
	"marcel-cli/config"
	"marcel-cli/models"
)

var ErrShowHelp = errors.New("show help")

func Run(args []string, version string) error {
	if len(args) == 0 {
		return ErrShowHelp
	}

	switch args[0] {
	case "help":
		ShowHelp()
		return nil
	case "version":
		fmt.Printf("marcel version %s\n", version)
		return nil
	case "auth":
		return runAuth(args[1:])
	case "quest":
		return runQuest(args[1:])
	case "journey":
		return runJourney(args[1:])
	case "habit":
		return runHabit(args[1:])
	case "event":
		return runEvent(args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func ShowHelp() {
	fmt.Println(`Marcel CLI

USAGE:
    marcel                 Start the interactive TUI
    marcel tui             Start the interactive TUI
    marcel help            Show this help message
    marcel version         Show version information
    marcel auth check      Verify authentication

QUESTS:
    marcel quest list [--json]
    marcel quest add <title> [--note <text>] [--difficulty <level>] [--journey <id>]
    marcel quest done <id>
    marcel quest undo <id>
    marcel quest toggle <id>
    marcel quest update <id> [--title <text>] [--note <text>] [--difficulty <level>]
    marcel quest delete <id>

JOURNEYS:
    marcel journey list [--json]
    marcel journey add <name>
    marcel journey update <id> --name <text>
    marcel journey delete <id>

HABITS:
    marcel habit list [--json]
    marcel habit add <name> [--cycle <type>]
    marcel habit done <id>
    marcel habit undo <id>
    marcel habit update <id> [--name <text>] [--cycle <type>]
    marcel habit delete <id>

EVENTS:
    marcel event list [--json] [--all]
    marcel event add <title> --date <YYYY-MM-DD> [--time <HH:MM>] [--end-date <YYYY-MM-DD>] [--end-time <HH:MM>]
    marcel event update <id> [--title <text>] [--date <YYYY-MM-DD>] [--time <HH:MM>] [--end-date <YYYY-MM-DD>] [--end-time <HH:MM>]
    marcel event delete <id>

AUTHENTICATION:
    Set MARCEL_TOKEN or create ~/.marcel.token`)
}

func runAuth(args []string) error {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("usage: marcel auth check")
		return nil
	}

	if args[0] != "check" {
		return fmt.Errorf("unknown auth command %q", args[0])
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	if err := client.CheckAuth(); err != nil {
		return err
	}

	fmt.Println("Authentication OK")
	return nil
}

func runQuest(args []string) error {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("usage: marcel quest [list|add|done|undo|toggle|update|delete]")
		return nil
	}

	switch args[0] {
	case "list":
		return questList(args[1:])
	case "add":
		return questAdd(args[1:])
	case "done":
		return questDone(args[1:], true)
	case "undo":
		return questDone(args[1:], false)
	case "toggle":
		return questToggle(args[1:])
	case "update":
		return questUpdate(args[1:])
	case "delete":
		return questDelete(args[1:])
	default:
		return fmt.Errorf("unknown quest command %q", args[0])
	}
}

func runJourney(args []string) error {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("usage: marcel journey [list|add|update|delete]")
		return nil
	}

	switch args[0] {
	case "list":
		return journeyList(args[1:])
	case "add":
		return journeyAdd(args[1:])
	case "update":
		return journeyUpdate(args[1:])
	case "delete":
		return journeyDelete(args[1:])
	default:
		return fmt.Errorf("unknown journey command %q", args[0])
	}
}

func runHabit(args []string) error {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("usage: marcel habit [list|add|done|undo|update|delete]")
		return nil
	}

	switch args[0] {
	case "list":
		return habitList(args[1:])
	case "add":
		return habitAdd(args[1:])
	case "done":
		return habitDone(args[1:], true)
	case "undo":
		return habitDone(args[1:], false)
	case "update":
		return habitUpdate(args[1:])
	case "delete":
		return habitDelete(args[1:])
	default:
		return fmt.Errorf("unknown habit command %q", args[0])
	}
}

func runEvent(args []string) error {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("usage: marcel event [list|add|update|delete]")
		return nil
	}

	switch args[0] {
	case "list":
		return eventList(args[1:])
	case "add":
		return eventAdd(args[1:])
	case "update":
		return eventUpdate(args[1:])
	case "delete":
		return eventDelete(args[1:])
	default:
		return fmt.Errorf("unknown event command %q", args[0])
	}
}

func questList(args []string) error {
	fs := newFlagSet("quest list")
	jsonOutput := fs.Bool("json", false, "output JSON")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	quests, err := client.GetQuests()
	if err != nil {
		return err
	}

	sort.Slice(quests, func(i, j int) bool {
		if quests[i].Done != quests[j].Done {
			return !quests[i].Done
		}
		return quests[i].ID < quests[j].ID
	})

	if *jsonOutput {
		return printJSON(quests)
	}

	journeyNames := map[int]string{}
	journeys, err := client.GetJourneys()
	if err == nil {
		for _, journey := range journeys {
			journeyNames[journey.ID] = journey.Name
		}
	}

	rows := make([][]string, 0, len(quests))
	for _, quest := range quests {
		journey := "-"
		if quest.JourneyID != nil {
			if name, ok := journeyNames[*quest.JourneyID]; ok {
				journey = name
			} else {
				journey = strconv.Itoa(*quest.JourneyID)
			}
		}

		rows = append(rows, []string{
			strconv.Itoa(quest.ID),
			statusLabel(quest.Done),
			quest.Title,
			emptyFallback(quest.Difficulty, "-"),
			journey,
		})
	}

	printTable([]string{"ID", "STATUS", "TITLE", "DIFFICULTY", "JOURNEY"}, rows)
	return nil
}

func questAdd(args []string) error {
	fs := newFlagSet("quest add")
	note := fs.String("note", "", "quest note")
	difficulty := fs.String("difficulty", "medium", "difficulty level")
	var journeyID intValue
	fs.Var(&journeyID, "journey", "journey ID")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("usage: marcel quest add <title> [--note <text>] [--difficulty <level>] [--journey <id>]")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	var journeyPtr *int
	if journeyID.set {
		journeyPtr = &journeyID.value
	}

	quest, err := client.CreateQuest(fs.Arg(0), *note, *difficulty, journeyPtr, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Created quest %d: %s\n", quest.ID, quest.Title)
	return nil
}

func questDone(args []string, done bool) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: marcel quest %s <id>", map[bool]string{true: "done", false: "undo"}[done])
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(args[0], "quest")
	if err != nil {
		return err
	}

	quest, err := client.ToggleQuest(id, done)
	if err != nil {
		return err
	}

	fmt.Printf("Updated quest %d: %s\n", quest.ID, statusLabel(quest.Done))
	return nil
}

func questToggle(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: marcel quest toggle <id>")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(args[0], "quest")
	if err != nil {
		return err
	}

	quests, err := client.GetQuests()
	if err != nil {
		return err
	}

	for _, quest := range quests {
		if quest.ID == id {
			updated, err := client.ToggleQuest(id, !quest.Done)
			if err != nil {
				return err
			}
			fmt.Printf("Updated quest %d: %s\n", updated.ID, statusLabel(updated.Done))
			return nil
		}
	}

	return fmt.Errorf("quest %d not found", id)
}

func questUpdate(args []string) error {
	fs := newFlagSet("quest update")
	var title trackedString
	var note trackedString
	var difficulty trackedString
	fs.Var(&title, "title", "new title")
	fs.Var(&note, "note", "new note")
	fs.Var(&difficulty, "difficulty", "new difficulty")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("usage: marcel quest update <id> [--title <text>] [--note <text>] [--difficulty <level>]")
	}

	updates := api.UpdateQuestRequest{}
	if title.set {
		updates.Title = &title.value
	}
	if note.set {
		updates.Note = &note.value
	}
	if difficulty.set {
		updates.Difficulty = &difficulty.value
	}
	if updates.Title == nil && updates.Note == nil && updates.Difficulty == nil {
		return fmt.Errorf("no updates provided")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(fs.Arg(0), "quest")
	if err != nil {
		return err
	}

	quest, err := client.UpdateQuest(id, updates)
	if err != nil {
		return err
	}

	fmt.Printf("Updated quest %d: %s\n", quest.ID, quest.Title)
	return nil
}

func questDelete(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: marcel quest delete <id>")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(args[0], "quest")
	if err != nil {
		return err
	}

	if err := client.DeleteQuest(id); err != nil {
		return err
	}

	fmt.Printf("Deleted quest %d\n", id)
	return nil
}

func journeyList(args []string) error {
	fs := newFlagSet("journey list")
	jsonOutput := fs.Bool("json", false, "output JSON")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	journeys, err := client.GetJourneys()
	if err != nil {
		return err
	}

	sort.Slice(journeys, func(i, j int) bool {
		return journeys[i].ID < journeys[j].ID
	})

	if *jsonOutput {
		return printJSON(journeys)
	}

	rows := make([][]string, 0, len(journeys))
	for _, journey := range journeys {
		rows = append(rows, []string{
			strconv.Itoa(journey.ID),
			journey.Name,
		})
	}

	printTable([]string{"ID", "NAME"}, rows)
	return nil
}

func journeyAdd(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: marcel journey add <name>")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	journey, err := client.CreateJourney(args[0], nil)
	if err != nil {
		return err
	}

	fmt.Printf("Created journey %d: %s\n", journey.ID, journey.Name)
	return nil
}

func journeyUpdate(args []string) error {
	fs := newFlagSet("journey update")
	var name trackedString
	fs.Var(&name, "name", "new name")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("usage: marcel journey update <id> --name <text>")
	}
	if !name.set {
		return fmt.Errorf("missing required flag --name")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(fs.Arg(0), "journey")
	if err != nil {
		return err
	}

	journey, err := client.UpdateJourney(id, api.UpdateJourneyRequest{Name: &name.value})
	if err != nil {
		return err
	}

	fmt.Printf("Updated journey %d: %s\n", journey.ID, journey.Name)
	return nil
}

func journeyDelete(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: marcel journey delete <id>")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(args[0], "journey")
	if err != nil {
		return err
	}

	if err := client.DeleteJourney(id); err != nil {
		return err
	}

	fmt.Printf("Deleted journey %d\n", id)
	return nil
}

func habitList(args []string) error {
	fs := newFlagSet("habit list")
	jsonOutput := fs.Bool("json", false, "output JSON")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	habits, err := client.GetHabits()
	if err != nil {
		return err
	}

	sort.Slice(habits, func(i, j int) bool {
		if habits[i].IsDueToday != habits[j].IsDueToday {
			return habits[i].IsDueToday
		}
		return habits[i].ID < habits[j].ID
	})

	if *jsonOutput {
		return printJSON(habits)
	}

	rows := make([][]string, 0, len(habits))
	for _, habit := range habits {
		rows = append(rows, []string{
			strconv.Itoa(habit.ID),
			habit.Name,
			habit.CycleType,
			boolLabel(habit.IsDueToday),
			strconv.Itoa(habit.CurrentStreak),
		})
	}

	printTable([]string{"ID", "NAME", "CYCLE", "DUE_TODAY", "STREAK"}, rows)
	return nil
}

func habitAdd(args []string) error {
	fs := newFlagSet("habit add")
	cycle := fs.String("cycle", "daily", "habit cycle type")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("usage: marcel habit add <name> [--cycle <type>]")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	habit, err := client.CreateHabit(fs.Arg(0), *cycle, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Created habit %d: %s\n", habit.ID, habit.Name)
	return nil
}

func habitDone(args []string, done bool) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: marcel habit %s <id>", map[bool]string{true: "done", false: "undo"}[done])
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(args[0], "habit")
	if err != nil {
		return err
	}

	habit, err := client.ToggleHabit(id, done)
	if err != nil {
		return err
	}

	fmt.Printf("Updated habit %d: %s\n", habit.ID, habit.Name)
	return nil
}

func habitUpdate(args []string) error {
	fs := newFlagSet("habit update")
	var name trackedString
	var cycle trackedString
	fs.Var(&name, "name", "new habit name")
	fs.Var(&cycle, "cycle", "new cycle type")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("usage: marcel habit update <id> [--name <text>] [--cycle <type>]")
	}

	updates := api.UpdateHabitRequest{}
	if name.set {
		updates.Name = &name.value
	}
	if cycle.set {
		updates.CycleType = &cycle.value
	}
	if updates.Name == nil && updates.CycleType == nil {
		return fmt.Errorf("no updates provided")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(fs.Arg(0), "habit")
	if err != nil {
		return err
	}

	habit, err := client.UpdateHabit(id, updates)
	if err != nil {
		return err
	}

	fmt.Printf("Updated habit %d: %s\n", habit.ID, habit.Name)
	return nil
}

func habitDelete(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: marcel habit delete <id>")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(args[0], "habit")
	if err != nil {
		return err
	}

	if err := client.DeleteHabit(id); err != nil {
		return err
	}

	fmt.Printf("Deleted habit %d\n", id)
	return nil
}

func eventList(args []string) error {
	fs := newFlagSet("event list")
	jsonOutput := fs.Bool("json", false, "output JSON")
	allEvents := fs.Bool("all", false, "include past events")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	events, err := client.GetEvents()
	if err != nil {
		return err
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Before(events[j].Date)
	})

	if !*allEvents {
		today := time.Now().In(time.Local)
		startOfToday := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

		filtered := make([]models.Event, 0, len(events))
		for _, event := range events {
			if !event.Date.Before(startOfToday) {
				filtered = append(filtered, event)
			}
		}
		events = filtered
	}

	if *jsonOutput {
		return printJSON(events)
	}

	rows := make([][]string, 0, len(events))
	for _, event := range events {
		rows = append(rows, []string{
			strconv.Itoa(event.ID),
			event.Title,
			formatDate(event.Date),
			ptrString(event.Time),
			ptrString(event.Location),
		})
	}

	printTable([]string{"ID", "TITLE", "DATE", "TIME", "LOCATION"}, rows)
	return nil
}

func eventAdd(args []string) error {
	fs := newFlagSet("event add")
	date := fs.String("date", "", "event date (YYYY-MM-DD)")
	timeValue := fs.String("time", "", "event time (HH:MM)")
	endDate := fs.String("end-date", "", "end date (YYYY-MM-DD)")
	endTime := fs.String("end-time", "", "end time (HH:MM)")
	location := fs.String("location", "", "event location")
	description := fs.String("description", "", "event description")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("usage: marcel event add <title> --date <YYYY-MM-DD> [--time <HH:MM>] [--end-date <YYYY-MM-DD>] [--end-time <HH:MM>]")
	}
	if *date == "" {
		return fmt.Errorf("missing required flag --date")
	}

	req := api.CreateEventRequest{
		Title: fs.Arg(0),
		Date:  *date,
	}
	if *timeValue != "" {
		req.Time = timeValue
	}
	if *endDate != "" {
		req.EndDate = endDate
	}
	if *endTime != "" {
		req.EndTime = endTime
	}
	if *location != "" {
		req.Location = location
	}
	if *description != "" {
		req.Description = description
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	event, err := client.CreateEvent(req)
	if err != nil {
		return err
	}

	fmt.Printf("Created event %d: %s\n", event.ID, event.Title)
	return nil
}

func eventUpdate(args []string) error {
	fs := newFlagSet("event update")
	var title trackedString
	var date trackedString
	var timeValue trackedString
	var endDate trackedString
	var endTime trackedString
	var location trackedString
	var description trackedString
	fs.Var(&title, "title", "new title")
	fs.Var(&date, "date", "new date")
	fs.Var(&timeValue, "time", "new time")
	fs.Var(&endDate, "end-date", "new end date")
	fs.Var(&endTime, "end-time", "new end time")
	fs.Var(&location, "location", "new location")
	fs.Var(&description, "description", "new description")
	if err := parseArgs(fs, args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("usage: marcel event update <id> [--title <text>] [--date <YYYY-MM-DD>] [--time <HH:MM>] [--end-date <YYYY-MM-DD>] [--end-time <HH:MM>]")
	}

	updates := api.UpdateEventRequest{}
	if title.set {
		updates.Title = &title.value
	}
	if date.set {
		updates.Date = &date.value
	}
	if timeValue.set {
		updates.Time = &timeValue.value
	}
	if endDate.set {
		updates.EndDate = &endDate.value
	}
	if endTime.set {
		updates.EndTime = &endTime.value
	}
	if location.set {
		updates.Location = &location.value
	}
	if description.set {
		updates.Description = &description.value
	}
	if updates.Title == nil && updates.Date == nil && updates.Time == nil && updates.EndDate == nil && updates.EndTime == nil && updates.Location == nil && updates.Description == nil {
		return fmt.Errorf("no updates provided")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(fs.Arg(0), "event")
	if err != nil {
		return err
	}

	event, err := client.UpdateEvent(id, updates)
	if err != nil {
		return err
	}

	fmt.Printf("Updated event %d: %s\n", event.ID, event.Title)
	return nil
}

func eventDelete(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: marcel event delete <id>")
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	id, err := parseID(args[0], "event")
	if err != nil {
		return err
	}

	if err := client.DeleteEvent(id); err != nil {
		return err
	}

	fmt.Printf("Deleted event %d\n", id)
	return nil
}

func newClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return api.NewClient(cfg), nil
}

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {}
	return fs
}

func parseArgs(fs *flag.FlagSet, args []string) error {
	normalized := normalizeArgs(fs, args)
	if err := fs.Parse(normalized); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return ErrShowHelp
		}
		return err
	}
	return nil
}

func normalizeArgs(fs *flag.FlagSet, args []string) []string {
	flags := make([]string, 0, len(args))
	positionals := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			positionals = append(positionals, args[i+1:]...)
			break
		}

		if !looksLikeFlag(arg) {
			positionals = append(positionals, arg)
			continue
		}

		flags = append(flags, arg)
		if flagConsumesValue(fs, arg) && i+1 < len(args) {
			i++
			flags = append(flags, args[i])
		}
	}

	return append(flags, positionals...)
}

func looksLikeFlag(arg string) bool {
	return len(arg) > 1 && strings.HasPrefix(arg, "-")
}

func flagConsumesValue(fs *flag.FlagSet, arg string) bool {
	if strings.Contains(arg, "=") {
		return false
	}

	name := strings.TrimLeft(arg, "-")
	if name == "" {
		return false
	}

	flagDef := fs.Lookup(name)
	if flagDef == nil {
		return false
	}

	if boolFlag, ok := flagDef.Value.(interface{ IsBoolFlag() bool }); ok && boolFlag.IsBoolFlag() {
		return false
	}

	return true
}

func parseID(raw, resource string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid %s ID %q", resource, raw)
	}
	return id, nil
}

func printJSON(v any) error {
	encoded, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(encoded))
	return nil
}

func printTable(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}

func statusLabel(done bool) string {
	if done {
		return "done"
	}
	return "open"
}

func boolLabel(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func ptrString(value *string) string {
	if value == nil || *value == "" {
		return "-"
	}
	return *value
}

func emptyFallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func formatDate(date time.Time) string {
	if date.IsZero() {
		return "-"
	}
	return date.Format("2006-01-02")
}

type trackedString struct {
	value string
	set   bool
}

func (t *trackedString) Set(value string) error {
	t.value = value
	t.set = true
	return nil
}

func (t *trackedString) String() string {
	return t.value
}

type intValue struct {
	value int
	set   bool
}

func (i *intValue) Set(value string) error {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	i.value = parsed
	i.set = true
	return nil
}

func (i *intValue) String() string {
	if !i.set {
		return ""
	}
	return strconv.Itoa(i.value)
}

var (
	_ flag.Value = (*trackedString)(nil)
	_ flag.Value = (*intValue)(nil)
)
