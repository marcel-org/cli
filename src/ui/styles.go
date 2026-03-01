package ui

import (
	"fmt"
	"marcel-cli/ui/colors"

	"github.com/charmbracelet/lipgloss"
)

var (
	BaseStyle = lipgloss.NewStyle().
			Foreground(colors.PrimaryText).
			Background(colors.BackgroundPrimary)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.BrandOrange).
			MarginBottom(1).
			Padding(0, 2)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.BrandOrange).
			Background(colors.BackgroundSecondary).
			Padding(0, 1).
			MarginBottom(1)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(colors.BackgroundPrimary).
				Background(colors.BrandOrange).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(colors.PrimaryText)

	CompletedItemStyle = lipgloss.NewStyle().
				Foreground(colors.MutedText).
				Strikethrough(true).
				Padding(0, 1)

	MutedStyle = lipgloss.NewStyle().
			Foreground(colors.SecondaryText)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(colors.Red).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(colors.Green).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Padding(1, 2)

	HelpStyle = lipgloss.NewStyle().
			Foreground(colors.SecondaryText).
			MarginTop(1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(colors.PrimaryText).
			Background(colors.BackgroundSecondary).
			Padding(0, 1)

	BadgeStyle = lipgloss.NewStyle().
			Foreground(colors.BackgroundPrimary).
			Background(colors.BrandOrange).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(colors.BrandOrange)

	DividerStyle = lipgloss.NewStyle().
			Foreground(colors.BorderColor).
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
		fmt.Sprintf(" ⚡ +%d  💰 +%d", xp, gold),
	)
}
