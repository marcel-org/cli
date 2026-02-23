# Marcel CLI

Terminal interface for managing your Marcel quests and calendar.

## Features

- View and manage quests with an interactive TUI
- Toggle quest completion with instant sync
- Create and delete quests
- Real-time filtering
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
go build -o marcel
cp marcel ~/.local/bin/marcel
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

The API endpoint is hardcoded to `https://api.marcel.my`.

## Tech Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - Components
- [Huh](https://github.com/charmbracelet/huh) - Forms
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
