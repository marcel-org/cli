# Marcel CLI

A TUI (Terminal User Interface) application for viewing and managing your Marcel quests from the command line.

## Features

- 🔐 Secure authentication with Marcel tokens
- 📋 View your quests organized by journeys
- ✅ Toggle quest completion with instant sync
- 🔄 Refresh quests in real-time
- 🎨 Clean and intuitive interface using Bubble Tea
- ⚙️ Configuration via environment variables

## Installation

### Build from source

```bash
cd cli
make build
make install
```

This will install the `marcel` binary to `~/.local/bin/`.

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

**Quest View:**
- `↑/↓` or `j/k` - Navigate quests
- `Space` - Toggle quest completion
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

### API Endpoint

The API endpoint can be configured in two ways:

1. Environment variable: `MARCEL_API_ENDPOINT` (takes priority)
2. Configuration file `~/.marcel.yml`:
   ```yaml
   api_endpoint: https://api.marcel.my
   ```

Default: `https://api.marcel.my`

## Tech Stack

- Go 1.25.0
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## License

See the main Marcel project for license information.
