package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	brandOrange = lipgloss.Color("#FF9600")
	white       = lipgloss.Color("#FFFFFF")
	gray        = lipgloss.Color("#666666")
	lightGray   = lipgloss.Color("#888888")
	darkGray    = lipgloss.Color("#333333")
	green       = lipgloss.Color("#34C759")
	red         = lipgloss.Color("#FF3B30")
	background  = lipgloss.Color("#1a1a1a")
)

var (
	BaseStyle = lipgloss.NewStyle().
			Foreground(white).
			Background(background)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(brandOrange).
			MarginBottom(1).
			Padding(0, 2)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(brandOrange).
			Background(darkGray).
			Padding(0, 1).
			MarginBottom(1)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(background).
				Background(brandOrange).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(white)

	CompletedItemStyle = lipgloss.NewStyle().
				Foreground(gray).
				Strikethrough(true).
				Padding(0, 1)

	MutedStyle = lipgloss.NewStyle().
			Foreground(lightGray)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(brandOrange).
			Padding(1, 2)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lightGray).
			MarginTop(1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(white).
			Background(darkGray).
			Padding(0, 1)

	BadgeStyle = lipgloss.NewStyle().
			Foreground(background).
			Background(brandOrange).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(brandOrange)

	DividerStyle = lipgloss.NewStyle().
			Foreground(darkGray).
			Bold(true)
)

func RenderDivider(width int) string {
	if width <= 0 {
		width = 80
	}
	divider := ""
	for i := 0; i < width; i++ {
		divider += "─"
	}
	return DividerStyle.Render(divider)
}

func RenderBox(title, content string) string {
	titleWithStyle := TitleStyle.Render(title)
	boxContent := titleWithStyle + "\n\n" + content
	return BoxStyle.Render(boxContent)
}

func RenderBadge(text string) string {
	return BadgeStyle.Render(text)
}

func RenderQuestCheckbox(done bool) string {
	if done {
		return CompletedItemStyle.Render("[✓]")
	}
	return NormalItemStyle.Render("[ ]")
}

func RenderReward(xp, gold int) string {
	return MutedStyle.Render(
		fmt.Sprintf(" ⚡ +%d  🪙 +%d", xp, gold),
	)
}
