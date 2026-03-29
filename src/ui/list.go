package ui

import (
	"fmt"
	"io"
	"marcel-cli/models"
	"marcel-cli/ui/colors"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type questItem struct {
	quest   models.Quest
	journey string
}

func (i questItem) Title() string {
	checkbox := "[ ]"
	if i.quest.Done {
		checkbox = "[✓]"
	}
	return fmt.Sprintf("%s %s", checkbox, i.quest.Title)
}

func (i questItem) Description() string {
	reward := fmt.Sprintf("⚡ +%d  💰 +%d", i.quest.XPReward, i.quest.GoldReward)
	if i.journey != "" {
		return fmt.Sprintf("%s • %s", i.journey, reward)
	}
	return reward
}

func (i questItem) FilterValue() string {
	return i.quest.Title
}

type questDelegate struct{}

func (d questDelegate) Height() int { return 1 }

func (d questDelegate) Spacing() int { return 0 }

func (d questDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d questDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(questItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	checkbox := "☐"
	if i.quest.Done {
		checkbox = "☑"
	}

	reward := fmt.Sprintf("⚡ +%d  💰 +%d", i.quest.XPReward, i.quest.GoldReward)
	journey := ""
	if i.journey != "" && i.journey != "My Quests" {
		journey = fmt.Sprintf("[%s] ", i.journey)
	}

	var str string
	if i.quest.Done {
		title := lipgloss.NewStyle().
			Foreground(colors.MutedText).
			Strikethrough(true).
			Render(fmt.Sprintf("%s %s", checkbox, i.quest.Title))

		meta := lipgloss.NewStyle().
			Foreground(colors.MutedText).
			Render(fmt.Sprintf(" %s%s", journey, reward))

		if isSelected {
			str = lipgloss.NewStyle().
				Background(colors.BrandOrange).
				Foreground(colors.BackgroundPrimary).
				Render(title + meta)
		} else {
			str = title + meta
		}
	} else {
		journeyStyled := ""
		if journey != "" && journey != "My Quests" {
			journeyStyled = lipgloss.NewStyle().
				Foreground(colors.SecondaryText).
				Render(fmt.Sprintf(" %s", journey))
		}

		titleText := fmt.Sprintf("%s %s", checkbox, i.quest.Title)

		if isSelected {
			str = SelectedItemStyle.Render(titleText) + journeyStyled
		} else {
			difficultyColor := colors.PrimaryText
			switch i.quest.Difficulty {
			case "easy":
				difficultyColor = colors.DifficultyEasy
			case "medium":
				difficultyColor = colors.DifficultyMedium
			case "hard":
				difficultyColor = colors.DifficultyHard
			case "epic":
				difficultyColor = colors.DifficultyEpic
			case "legendary":
				difficultyColor = colors.DifficultyLegendary
			}

			incompleteStyle := lipgloss.NewStyle().
				Foreground(difficultyColor).
				Bold(true)
			str = incompleteStyle.Render(titleText) + journeyStyled
		}
	}

	fmt.Fprint(w, str)
}

func newQuestList(data *models.AppData, width, height int) list.Model {
	return newQuestListFromJourneys(data.Journeys, width, height)
}

func newQuestListFromJourneys(journeys []models.Journey, width, height int) list.Model {
	incompleteItems := []list.Item{}
	completedItems := []list.Item{}

	for _, journey := range journeys {
		for _, quest := range journey.Quests {
			item := questItem{
				quest:   quest,
				journey: journey.Name,
			}
			if quest.Done {
				completedItems = append(completedItems, item)
			} else {
				incompleteItems = append(incompleteItems, item)
			}
		}
	}

	items := append(incompleteItems, completedItems...)

	delegate := questDelegate{}

	l := list.New(items, delegate, width, height)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(colors.BrandOrange)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(colors.BrandOrange)

	return l
}

type habitItem struct {
	habit models.Habit
}

func (i habitItem) Title() string {
	checkbox := "☐"
	if i.habit.IsDueToday {
		today := ""
		for _, d := range i.habit.Completed {
			if d[:10] == fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day()) {
				today = d
				break
			}
		}
		if today != "" {
			checkbox = "☑"
		}
	}
	return fmt.Sprintf("%s %s", checkbox, i.habit.Name)
}

func (i habitItem) Description() string {
	return fmt.Sprintf("🔥 %d streak • %s", i.habit.CurrentStreak, i.habit.CycleDescription)
}

func (i habitItem) FilterValue() string {
	return i.habit.Name
}

type habitDelegate struct{}

func (d habitDelegate) Height() int { return 1 }

func (d habitDelegate) Spacing() int { return 0 }

func (d habitDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d habitDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(habitItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	checkbox := "☐"
	completedToday := false
	today := fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	for _, d := range i.habit.Completed {
		if len(d) >= 10 && d[:10] == today {
			completedToday = true
			checkbox = "☑"
			break
		}
	}

	streak := lipgloss.NewStyle().
		Foreground(colors.SecondaryText).
		Render(fmt.Sprintf(" 🔥 %d • %s", i.habit.CurrentStreak, i.habit.CycleDescription))

	var str string
	if completedToday {
		title := lipgloss.NewStyle().
			Foreground(colors.MutedText).
			Strikethrough(true).
			Render(fmt.Sprintf("%s %s", checkbox, i.habit.Name))

		if isSelected {
			str = lipgloss.NewStyle().
				Background(colors.BrandOrange).
				Foreground(colors.BackgroundPrimary).
				Render(title + streak)
		} else {
			str = title + streak
		}
	} else {
		if isSelected {
			str = SelectedItemStyle.Render(fmt.Sprintf("%s %s", checkbox, i.habit.Name)) + streak
		} else {
			str = NormalItemStyle.Render(fmt.Sprintf("%s %s", checkbox, i.habit.Name)) + streak
		}
	}

	fmt.Fprint(w, str)
}

func newHabitList(data *models.AppData, width, height int) list.Model {
	items := []list.Item{}

	for _, habit := range data.Habits {
		items = append(items, habitItem{
			habit: habit,
		})
	}

	delegate := habitDelegate{}

	l := list.New(items, delegate, width, height)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(colors.BrandOrange)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(colors.BrandOrange)

	return l
}

type journeyItem struct {
	journey models.Journey
}

func (i journeyItem) Title() string {
	return i.journey.Name
}

func (i journeyItem) Description() string {
	return fmt.Sprintf("%d quests", len(i.journey.Quests))
}

func (i journeyItem) FilterValue() string {
	return i.journey.Name
}

type journeyDelegate struct{}

func (d journeyDelegate) Height() int { return 1 }

func (d journeyDelegate) Spacing() int { return 0 }

func (d journeyDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d journeyDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(journeyItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	questCount := lipgloss.NewStyle().
		Foreground(colors.SecondaryText).
		Render(fmt.Sprintf(" %d quests", len(i.journey.Quests)))

	var str string
	if isSelected {
		str = SelectedItemStyle.Render(i.journey.Name) + questCount
	} else {
		str = NormalItemStyle.Render(i.journey.Name) + questCount
	}

	fmt.Fprint(w, str)
}

func newJourneyList(data *models.AppData, width, height int) list.Model {
	return newJourneyListFromJourneys(data.Journeys, width, height)
}

func newJourneyListFromJourneys(journeys []models.Journey, width, height int) list.Model {
	items := []list.Item{}

	for _, journey := range journeys {
		if journey.ID != 0 {
			items = append(items, journeyItem{
				journey: journey,
			})
		}
	}

	delegate := journeyDelegate{}

	l := list.New(items, delegate, width, height)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(colors.BrandOrange)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(colors.BrandOrange)

	return l
}

func newJourneyQuestList(journey *models.Journey, width, height int) list.Model {
	return newQuestListFromJourneys([]models.Journey{*journey}, width, height)
}

type spaceItem struct {
	space models.Space
}

func (i spaceItem) Title() string {
	return i.space.Name
}

func (i spaceItem) Description() string {
	return fmt.Sprintf("%d members • owner %s", len(i.space.Members), i.space.Owner.Username)
}

func (i spaceItem) FilterValue() string {
	return i.space.Name
}

type spaceDelegate struct{}

func (d spaceDelegate) Height() int { return 1 }

func (d spaceDelegate) Spacing() int { return 0 }

func (d spaceDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d spaceDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(spaceItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()
	meta := lipgloss.NewStyle().
		Foreground(colors.SecondaryText).
		Render(fmt.Sprintf(" %d members • owner %s", len(i.space.Members), i.space.Owner.Username))

	var str string
	if isSelected {
		str = SelectedItemStyle.Render(i.space.Name) + meta
	} else {
		str = NormalItemStyle.Render(i.space.Name) + meta
	}

	fmt.Fprint(w, str)
}

func newSpaceList(data *models.AppData, width, height int) list.Model {
	items := []list.Item{}

	for _, space := range data.Spaces {
		items = append(items, spaceItem{space: space})
	}

	l := list.New(items, spaceDelegate{}, width, height)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(colors.BrandOrange)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(colors.BrandOrange)

	return l
}

type memberItem struct {
	member models.SpaceMember
}

func (i memberItem) Title() string {
	return i.member.User.Username
}

func (i memberItem) Description() string {
	return fmt.Sprintf("%s permission • level %d", i.member.Permission, i.member.User.Level)
}

func (i memberItem) FilterValue() string {
	return i.member.User.Username
}

type memberDelegate struct{}

func (d memberDelegate) Height() int { return 1 }

func (d memberDelegate) Spacing() int { return 0 }

func (d memberDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d memberDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(memberItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()
	meta := lipgloss.NewStyle().
		Foreground(colors.SecondaryText).
		Render(fmt.Sprintf(" %s • level %d", i.member.Permission, i.member.User.Level))

	var str string
	if isSelected {
		str = SelectedItemStyle.Render(i.member.User.Username) + meta
	} else {
		str = NormalItemStyle.Render(i.member.User.Username) + meta
	}

	fmt.Fprint(w, str)
}

func newMemberList(members []models.SpaceMember, width, height int) list.Model {
	items := make([]list.Item, 0, len(members))
	for _, member := range members {
		items = append(items, memberItem{member: member})
	}

	l := list.New(items, memberDelegate{}, width, height)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(colors.BrandOrange)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(colors.BrandOrange)

	return l
}
