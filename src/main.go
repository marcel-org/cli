package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"marcel-cli/cli"
	"marcel-cli/ui"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = cli.ShowHelp

	var showVersion = fs.Bool("version", false, "Show version information")
	var showHelp = fs.Bool("help", false, "Show help information")
	var showHelpShort = fs.Bool("h", false, "Show help information")

	if err := fs.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		os.Exit(2)
	}

	if *showVersion {
		fmt.Printf("marcel version %s\n", version)
		return
	}

	if *showHelp || *showHelpShort {
		cli.ShowHelp()
		return
	}

	if fs.NArg() > 0 {
		switch fs.Arg(0) {
		case "tui":
			runTUI()
			return
		default:
			if err := cli.Run(fs.Args(), version); err != nil {
				if errors.Is(err, cli.ErrShowHelp) {
					cli.ShowHelp()
					return
				}
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
	}

	runTUI()
}

func runTUI() {
	model, err := ui.NewModel()
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
