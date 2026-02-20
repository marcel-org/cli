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
		{Key: "Easy", Value: "easy"},
		{Key: "Medium", Value: "medium"},
		{Key: "Hard", Value: "hard"},
		{Key: "Epic", Value: "epic"},
		{Key: "Legendary", Value: "legendary"},
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
	).
		WithTheme(theme).
		WithWidth(60).
		WithHeight(10)

	err := huhForm.Run()
	if err != nil {
		return nil, err
	}

	return dialog, nil
}

type HabitForm struct {
	Name        string
	CycleType   string
	CycleConfig any
}

func NewHabitForm() (*HabitForm, error) {
	form := &HabitForm{}

	cycleOptions := []huh.Option[string]{
		{Key: "Daily", Value: "daily"},
		{Key: "Weekly", Value: "weekly"},
		{Key: "Interval", Value: "interval"},
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
				Title("Habit Name").
				Placeholder("What habit do you want to track?").
				Value(&form.Name).
				Validate(func(s string) error {
					if len(s) == 0 {
						return fmt.Errorf("name cannot be empty")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Cycle Type").
				Options(cycleOptions...).
				Value(&form.CycleType),
		),
	).WithTheme(theme)

	err := huhForm.Run()
	if err != nil {
		return nil, err
	}

	return form, nil
}

type JourneyForm struct {
	Name string
}

func NewJourneyForm() (*JourneyForm, error) {
	form := &JourneyForm{}

	theme := huh.ThemeCharm()
	theme.Focused.Base = lipgloss.NewStyle().BorderForeground(brandOrange)
	theme.Focused.Title = lipgloss.NewStyle().Foreground(brandOrange).Bold(true)
	theme.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(brandOrange)

	huhForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Journey Name").
				Placeholder("Name your journey").
				Value(&form.Name).
				Validate(func(s string) error {
					if len(s) == 0 {
						return fmt.Errorf("name cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(theme)

	err := huhForm.Run()
	if err != nil {
		return nil, err
	}

	return form, nil
}

func NewQuestFormSimple() (*QuestForm, error) {
	form := &QuestForm{}

	difficultyOptions := []huh.Option[string]{
		{Key: "Easy", Value: "easy"},
		{Key: "Medium", Value: "medium"},
		{Key: "Hard", Value: "hard"},
		{Key: "Epic", Value: "epic"},
		{Key: "Legendary", Value: "legendary"},
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
