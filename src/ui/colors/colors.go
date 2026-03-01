package colors

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	isDark = detectDarkBackground()

	BrandOrange         = chooseColor("#E85D00", "#FF9600")
	PrimaryText         = chooseColor("#1a1a1a", "#FFFFFF")
	SecondaryText       = chooseColor("#4a4a4a", "#888888")
	MutedText           = chooseColor("#666666", "#aaaaaa")
	BorderColor         = chooseColor("#cccccc", "#333333")
	BackgroundPrimary   = chooseColor("#ffffff", "#1a1a1a")
	BackgroundSecondary = chooseColor("#f5f5f5", "#2a2a2a")
	Green               = chooseColor("#2d9f5d", "#34C759")
	Red                 = chooseColor("#d63030", "#FF3B30")
	AccentPeach         = chooseColor("#e67e22", "#fab387")
	AccentBlue          = chooseColor("#2563eb", "#89b4fa")
	DifficultyEasy      = chooseColor("#0d9668", "#10b981")
	DifficultyMedium    = chooseColor("#2563eb", "#3b82f6")
	DifficultyHard      = chooseColor("#7c3aed", "#a855f7")
	DifficultyEpic      = chooseColor("#dc2626", "#ef4444")
	DifficultyLegendary = chooseColor("#d97706", "#f59e0b")
)

func detectDarkBackground() bool {
	if colorMode := os.Getenv("MARCEL_COLOR_MODE"); colorMode != "" {
		return strings.ToLower(colorMode) == "dark"
	}

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	return termenv.HasDarkBackground()
}

func chooseColor(light, dark string) lipgloss.TerminalColor {
	if isDark {
		return lipgloss.Color(dark)
	}
	return lipgloss.Color(light)
}
