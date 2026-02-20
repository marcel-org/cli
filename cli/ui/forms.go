package ui

import (
	"fmt"
	"marcel-cli/models"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type QuestForm struct {
	Title      string
	Note       string
	Difficulty string
	JourneyID  *int
}

func NewQuestForm(journeys []models.Journey) (*QuestForm, error) {
	form := &QuestForm{}

	difficultyOptions := []huh.Option[string]{
		{Key: "Easy (1-2 XP)", Value: "easy"},
		{Key: "Medium (3-5 XP)", Value: "medium"},
		{Key: "Hard (6-10 XP)", Value: "hard"},
	}

	journeyOptions := []huh.Option[int]{
		{Key: "No Journey", Value: 0},
	}
	for _, j := range journeys {
		journeyOptions = append(journeyOptions, huh.Option[int]{
			Key:   j.Name,
			Value: j.ID,
		})
	}

	theme := huh.ThemeCharm()
	theme.Focused.Base = lipgloss.NewStyle().BorderForeground(brandOrange)
	theme.Focused.Title = lipgloss.NewStyle().Foreground(brandOrange).Bold(true)
	theme.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(brandOrange)
	theme.Focused.SelectSelector = lipgloss.NewStyle().Foreground(brandOrange).SetString(">")
	theme.Focused.SelectedOption = lipgloss.NewStyle().Foreground(brandOrange)

	huhForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Quest Title").
				Placeholder("What do you want to accomplish?").
				Value(&form.Title).
				Validate(func(s string) error {
					if len(s) == 0 {
						return fmt.Errorf("title cannot be empty")
					}
					return nil
				}),

			huh.NewText().
				Title("Quest Note (optional)").
				Placeholder("Add any additional details...").
				Value(&form.Note),

			huh.NewSelect[string]().
				Title("Difficulty").
				Options(difficultyOptions...).
				Value(&form.Difficulty),
		),
	).WithTheme(theme)

	err := huhForm.Run()
	if err != nil {
		return nil, err
	}

	return form, nil
}

type ConfirmDialog struct {
	Confirmed bool
}

func NewConfirmDialog(title, description string) (*ConfirmDialog, error) {
	dialog := &ConfirmDialog{}

	theme := huh.ThemeCharm()
	theme.Focused.Base = lipgloss.NewStyle().BorderForeground(red)
	theme.Focused.Title = lipgloss.NewStyle().Foreground(red).Bold(true)
	theme.Focused.SelectSelector = lipgloss.NewStyle().Foreground(red).SetString(">")
	theme.Focused.SelectedOption = lipgloss.NewStyle().Foreground(red)

	huhForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Description(description).
				Affirmative("Yes").
				Negative("No").
				Value(&dialog.Confirmed),
		),
	).WithTheme(theme)

	err := huhForm.Run()
	if err != nil {
		return nil, err
	}

	return dialog, nil
}
