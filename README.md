# Marcel CLI

Terminal client for managing Marcel from either an interactive TUI or direct shell commands.

## Features

- Interactive TUI for quests and calendar
- Scriptable CLI subcommands for quests, journeys, habits, and events
- Plain table output for humans and `--json` output for scripts
- Instant sync for create, update, toggle, and delete actions
- Calendar view with week/month modes
- Modern UI with Charmbracelet components

## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/marcel-org/cli/main/install.sh | bash
```

### Build from source

```bash
git clone https://github.com/marcel-org/cli.git
cd cli
mise run install
```

Common developer tasks:

```bash
mise run build
mise run test
mise run install
mise run clean
```

## Authentication

1. Go to Marcel web app Settings → Marcel CLI
2. Generate a CLI token
3. Set the token in your environment:

```bash
export MARCEL_TOKEN="marcel_your_token_here"
```

Add this to your shell config (`~/.zshrc`, `~/.bashrc`) to make it permanent.

## Usage

```bash
marcel
```

This still starts the TUI by default. You can also start it explicitly:

```bash
marcel tui
```

### Direct CLI Commands

```bash
marcel help
marcel auth check

marcel quest list
marcel quest list --json
marcel quest add "Write proposal" --note "Draft outline" --difficulty medium
marcel quest done 42
marcel quest undo 42
marcel quest toggle 42
marcel quest update 42 --title "Write final proposal"
marcel quest delete 42

marcel journey list
marcel journey add "Health"
marcel journey update 3 --name "Fitness"
marcel journey delete 3

marcel habit list
marcel habit add "Meditate" --cycle daily
marcel habit done 8
marcel habit undo 8
marcel habit update 8 --name "Walk" --cycle weekly
marcel habit delete 8

marcel event list
marcel event list --all
marcel event add "Doctor" --date 2026-03-28 --time 09:00 --location "Berlin"
marcel event update 15 --title "Dentist" --time 10:30
marcel event delete 15
```

### Keyboard Controls

**Quest List:**
- `↑/↓` or `j/k` - Navigate
- `Space/Enter` - Toggle completion
- `n` - New quest
- `d` - Delete quest
- `/` - Filter
- `r` - Refresh
- `?` - Help
- `q` - Quit

**Calendar:**
- `w` - Week view
- `m` - Month view
- `tab` - Switch views

## Configuration

Optional `~/.marcel.yml` file:

```yaml
week_start_day: sunday  # Options: sunday, monday, tuesday, etc.
```

`list` commands support `--json` for machine-readable output. `marcel event list` shows only today's and future events by default; use `--all` to include past events.

## Tech Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - Components
- [Huh](https://github.com/charmbracelet/huh) - Forms
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
