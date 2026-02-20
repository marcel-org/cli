package ui

import (
	"fmt"
	"io"
	"marcel-cli/models"
	"strings"
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
	reward := fmt.Sprintf("⚡ +%d  🪙 +%d", i.quest.XPReward, i.quest.GoldReward)
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

	reward := fmt.Sprintf("⚡ +%d  🪙 +%d", i.quest.XPReward, i.quest.GoldReward)
	journey := ""
	if i.journey != "" && i.journey != "My Quests" {
		journey = fmt.Sprintf("[%s] ", i.journey)
	}

	var str string
	if i.quest.Done {
		title := lipgloss.NewStyle().
			Foreground(gray).
			Strikethrough(true).
			Render(fmt.Sprintf("%s %s", checkbox, i.quest.Title))

		meta := lipgloss.NewStyle().
			Foreground(gray).
			Render(fmt.Sprintf(" %s%s", journey, reward))

		if isSelected {
			str = lipgloss.NewStyle().
				Background(brandOrange).
				Render(title + meta)
		} else {
			str = title + meta
		}
	} else {
		rewardStyled := lipgloss.NewStyle().
			Foreground(gray).
			Render(fmt.Sprintf(" %s%s", journey, reward))

		if isSelected {
			str = SelectedItemStyle.Render(fmt.Sprintf("%s %s", checkbox, i.quest.Title)) + rewardStyled
		} else {
			str = NormalItemStyle.Render(fmt.Sprintf("%s %s", checkbox, i.quest.Title)) + rewardStyled
		}
	}

	fmt.Fprint(w, str)
}

func newQuestList(data *models.AppData, width, height int) list.Model {
	items := []list.Item{}

	for _, journey := range data.Journeys {
		for _, quest := range journey.Quests {
			items = append(items, questItem{
				quest:   quest,
				journey: journey.Name,
			})
		}
	}

	delegate := questDelegate{}

	l := list.New(items, delegate, width, height)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(brandOrange)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(brandOrange)

	return l
}

func questListHelp() string {
	helpItems := []string{
		"↑/k up",
		"↓/j down",
		"tab next",
		"shift+tab prev",
		"space toggle",
		"d delete",
		"r refresh",
		"n new",
		"? help",
		"q quit",
	}

	var styledItems []string
	for _, item := range helpItems {
		parts := strings.SplitN(item, " ", 2)
		if len(parts) == 2 {
			keyStyle := lipgloss.NewStyle().Foreground(brandOrange).Bold(true)
			descStyle := lipgloss.NewStyle().Foreground(lightGray)
			styledItems = append(styledItems, keyStyle.Render(parts[0])+" "+descStyle.Render(parts[1]))
		}
	}

	return HelpStyle.Render(strings.Join(styledItems, "  •  "))
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
		Foreground(gray).
		Render(fmt.Sprintf(" 🔥 %d • %s", i.habit.CurrentStreak, i.habit.CycleDescription))

	var str string
	if completedToday {
		title := lipgloss.NewStyle().
			Foreground(gray).
			Strikethrough(true).
			Render(fmt.Sprintf("%s %s", checkbox, i.habit.Name))

		if isSelected {
			str = lipgloss.NewStyle().
				Background(brandOrange).
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
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(brandOrange)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(brandOrange)

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
		Foreground(gray).
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
	items := []list.Item{}

	for _, journey := range data.Journeys {
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
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(brandOrange)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(brandOrange)

	return l
}

func newJourneyQuestList(journey *models.Journey, width, height int) list.Model {
	items := []list.Item{}

	for _, quest := range journey.Quests {
		items = append(items, questItem{
			quest:   quest,
			journey: journey.Name,
		})
	}

	delegate := questDelegate{}

	l := list.New(items, delegate, width, height)
	l.Title = ""
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(brandOrange)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(brandOrange)

	return l
}
