package ui

import (
	"fmt"
	"io"
	"marcel-cli/models"
	"strings"

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
	reward := fmt.Sprintf("⚡ +%d XP  💰 +%d gold", i.quest.XPReward, i.quest.GoldReward)
	if i.journey != "" {
		return fmt.Sprintf("%s • %s", i.journey, reward)
	}
	return reward
}

func (i questItem) FilterValue() string {
	return i.quest.Title
}

type questDelegate struct{}

func (d questDelegate) Height() int { return 2 }

func (d questDelegate) Spacing() int { return 1 }

func (d questDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d questDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(questItem)
	if !ok {
		return
	}

	var str string
	isSelected := index == m.Index()

	if i.quest.Done {
		titleStyle := CompletedItemStyle
		if isSelected {
			titleStyle = lipgloss.NewStyle().
				Foreground(gray).
				Background(brandOrange).
				Strikethrough(true).
				Padding(0, 1)
		}
		str = fmt.Sprintf("%s\n%s",
			titleStyle.Render(i.Title()),
			MutedStyle.Render("  "+i.Description()),
		)
	} else {
		titleStyle := NormalItemStyle
		if isSelected {
			titleStyle = SelectedItemStyle
		}
		str = fmt.Sprintf("%s\n%s",
			titleStyle.Render(i.Title()),
			MutedStyle.Render("  "+i.Description()),
		)
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
	l.Title = "Your Quests"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = TitleStyle
	l.Styles.TitleBar = lipgloss.NewStyle().
		Background(darkGray).
		Foreground(brandOrange).
		Bold(true).
		Padding(0, 1)
	l.Styles.StatusBar = StatusBarStyle
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(brandOrange)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(brandOrange)

	return l
}

func questListHelp() string {
	helpItems := []string{
		"↑/k up",
		"↓/j down",
		"/ filter",
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
