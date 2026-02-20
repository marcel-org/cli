# Marcel CLI

A TUI (Terminal User Interface) application for viewing and managing your Marcel quests from the command line.

## Features

- 🔐 Secure authentication with Marcel tokens
- 📋 View your quests organized by journeys with beautiful list component
- ✅ Toggle quest completion with instant sync
- ➕ Create new quests with interactive forms
- 🗑️ Delete quests with confirmation dialog
- 🔍 Filter quests in real-time
- 🔄 Refresh quests from server
- 🎨 Modern UI with Charmbracelet components (Bubbles, Huh, Lipgloss)
- ⚡ Smooth animations and loading states with spinners
- ⚙️ Configuration via environment variables

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/marcel-org/cli/main/install.sh | bash
```

This will clone the repository, build the binary, and install it to `~/.local/bin/marcel`.

### Build from source

```bash
git clone https://github.com/marcel-org/cli.git
cd cli
go build -o marcel
cp marcel ~/.local/bin/marcel
```

## Authentication

### Step 1: Get your Marcel CLI token

1. Go to Marcel web app settings
2. Click "Generate Marcel CLI Token"
3. Copy the token (it starts with `marcel_`)

### Step 2: Configure authentication

Set the `MARCEL_TOKEN` environment variable:

```bash
export MARCEL_TOKEN="marcel_your_token_here"
```

Add this to your shell configuration file (`~/.zshrc`, `~/.bashrc`, etc.) to make it permanent.

Optionally, you can also set the API endpoint (defaults to `https://api.marcel.my`):

```bash
export MARCEL_API_ENDPOINT="http://localhost:3000"  # for local development
```

## Usage

Simply run:

```bash
marcel
```

### Keyboard Controls

**Quest List View:**
- `↑/↓` or `j/k` - Navigate quests
- `gg` - Jump to top
- `G` - Jump to bottom
- `/` - Filter quests
- `Space` or `Enter` - Toggle quest completion
- `n` - Create new quest
- `d` - Delete quest (with confirmation)
- `r` - Refresh quests from server
- `?` - Show/hide help
- `q` or `Ctrl+C` - Quit

**Help View:**
- `?` or `Esc` - Return to quest view
- `q` or `Ctrl+C` - Quit

## Development

### Build

```bash
make build
```

### Build for all platforms

```bash
make build-all
```

### Install locally

```bash
make install
```

### Clean

```bash
make clean
```

## Configuration

### Authentication

Authentication is done **exclusively via the `MARCEL_TOKEN` environment variable**. The token cannot be set in the configuration file for security reasons.

### Configuration File

You can create a `~/.marcel.yml` file to configure the CLI:

```yaml
# API endpoint (can be overridden with MARCEL_API_ENDPOINT env var)
api_endpoint: https://api.marcel.my

# Week start day for calendar (0 = Sunday, 1 = Monday, etc.)
week_start_day: 1
```

**Example configuration file:**

```yaml
# ~/.marcel.yml
api_endpoint: https://api.marcel.my
week_start_day: 1  # Start week on Monday
```

The API endpoint can be configured in two ways:

1. Environment variable: `MARCEL_API_ENDPOINT` (takes priority)
2. Configuration file `~/.marcel.yml` (shown above)

Default: `https://api.marcel.my`

## Tech Stack

- Go 1.25.0
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components (list, spinner, etc.)
- [Huh](https://github.com/charmbracelet/huh) - Interactive forms and prompts
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## License

See the main Marcel project for license information.
