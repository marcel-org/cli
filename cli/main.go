package main

import (
	"flag"
	"fmt"
	"log"

	"marcel-cli/ui"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	var showVersion = flag.Bool("version", false, "Show version information")
	var showHelp = flag.Bool("help", false, "Show help information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("marcel version %s\n", version)
		return
	}

	if *showHelp {
		showHelpText()
		return
	}

	model, err := ui.NewModel()
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func showHelpText() {
	fmt.Println(`Marcel CLI - Gamified productivity TUI application

USAGE:
    marcel [OPTIONS]

OPTIONS:
    --version    Show version information
    --help       Show this help message

KEYBOARD CONTROLS:

Quest View:
    j/k, arrows  Navigate quests
    gg           Jump to top
    G            Jump to bottom
    Space        Toggle quest completion
    ?            Show/hide help
    q, Ctrl+C    Quit

Input Mode:
    Type         Enter text
    Enter        Confirm
    Esc          Cancel
    Backspace    Delete character

CONFIGURATION:
    API endpoint: Set in ~/.marcel.yml
    Auth token: Set MARCEL_TOKEN environment variable

For more information, visit: https://github.com/marcel-org/cli`)
}
