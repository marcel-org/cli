# Marcel CLI Integration - Complete Implementation

## Overview

The Marcel CLI now has full quest viewing and management capabilities with secure token-based authentication.

## What Was Implemented

### 1. Authentication & Configuration

**Environment Variable Support** (`cli/config/config.go`)
- `MARCEL_TOKEN` - Authentication token (highest priority)
- `MARCEL_API_ENDPOINT` - API endpoint URL
- Environment variables override config file settings

**Config File** (`~/.marcel.yml`)
```yaml
api_endpoint: http://localhost:3000
auth_token: marcel_your_token_here
```

### 2. API Client (`cli/api/`)

**client.go** - Core HTTP client
- Bearer token authentication
- Automatic header management
- Auth verification with `/user/me` endpoint
- Error handling and response parsing

**quests.go** - Quest CRUD operations
- `GetQuests()` - Fetch all user quests
- `CreateQuest()` - Create new quest
- `UpdateQuest()` - Update quest properties
- `ToggleQuest()` - Toggle quest completion status
- `DeleteQuest()` - Delete a quest

### 3. Updated Models (`cli/models/models.go`)

Quest model now matches Marcel backend API:
- `ID` (int) - Quest identifier
- `Title` (string) - Quest title
- `Note` (string) - Quest notes
- `Done` (bool) - Completion status
- `Difficulty` (string) - Quest difficulty
- `XPReward` / `GoldReward` (int) - Rewards
- `JourneyID` (*int) - Optional journey assignment
- Timestamps and metadata

### 4. Storage Layer (`cli/storage/storage.go`)

- Initializes API client with config
- Fetches quests from backend
- Groups quests by journey
- Creates "My Quests" journey for unassigned quests
- Exposes API client for direct operations

### 5. Enhanced UI (`cli/ui/ui.go`)

**Three View Modes:**
1. **QuestView** - Main quest list with navigation
2. **HelpView** - Keyboard shortcuts and configuration help
3. **ErrorView** - Authentication and connection errors

**Features:**
- ✅ Real-time quest toggling with API sync
- 🔄 Manual refresh with 'r' key
- 🎨 Color-coded UI (Marcel purple #5856D6)
- ✨ Success/error feedback messages
- 📊 Shows XP and gold rewards per quest
- 🎯 Journey-based quest grouping
- ⌨️ Vim-style navigation (j/k)

**Keyboard Controls:**
- `↑/↓` or `j/k` - Navigate quests
- `Space` - Toggle quest completion
- `r` - Refresh from server
- `?` - Toggle help
- `q` or `Ctrl+C` - Quit

## File Structure

```
cli/
├── api/
│   ├── client.go       # HTTP client with auth
│   └── quests.go       # Quest CRUD operations
├── config/
│   └── config.go       # Config loader with env support
├── models/
│   └── models.go       # Updated Quest & Journey models
├── storage/
│   └── storage.go      # API integration layer
├── ui/
│   └── ui.go           # Complete TUI implementation
├── main.go             # Entry point
├── go.mod              # Dependencies
├── Makefile            # Build commands
└── README.md           # Usage documentation
```

## Usage Flow

### 1. Authentication Setup

**Web App:**
1. Navigate to Marcel settings
2. Click "Generate Marcel CLI Token"
3. Copy the `marcel_xxx` token

**CLI Configuration:**
```bash
# Option A: Environment variable (recommended)
export MARCEL_TOKEN="marcel_abc123..."

# Option B: Config file
echo "auth_token: marcel_abc123..." >> ~/.marcel-cli.yml
```

### 2. Running the CLI

```bash
# Build and install
cd cli
make build
make install

# Run
marcel
```

### 3. Using the Interface

**On Launch:**
- Authenticates with backend
- Fetches all user quests
- Groups by journeys
- Displays in TUI

**Quest Management:**
- Navigate with arrow keys or j/k
- Press `Space` to mark quest as complete
- Quest updates sync immediately to backend
- XP and gold rewards shown for each quest

**Error Handling:**
- If authentication fails: Shows error screen with retry option
- If API unreachable: Shows connection error
- If token missing: Displays configuration instructions

## API Integration Details

### Authentication Flow

1. CLI loads config (file + env variables)
2. API client initialized with token
3. `CheckAuth()` called to verify token
4. If valid: Fetch quests
5. If invalid: Show error screen

### Quest Toggle Flow

1. User presses `Space` on selected quest
2. `ToggleQuest(questID, !currentStatus)` called
3. PUT request to `/quest/:id` with `done` field
4. On success: Update local state
5. Show success message
6. If error: Show error message

### Data Sync

- **On Launch:** Full quest list fetched
- **On Toggle:** Single quest updated via API
- **On Refresh (r):** Full quest list re-fetched
- **No local storage:** All data from API

## Backend API Endpoints Used

| Endpoint | Method | Purpose | Auth |
|----------|--------|---------|------|
| `/user/me` | GET | Verify token | Token |
| `/quest` | GET | Fetch all quests | Token |
| `/quest/:id` | PUT | Update quest | Token |
| `/quest` | POST | Create quest* | Token |
| `/quest/:id` | DELETE | Delete quest* | Token |

\* Create/Delete implemented in API client but not yet in UI

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `MARCEL_TOKEN` | Marcel CLI token | `marcel_abc123...` |
| `MARCEL_API_ENDPOINT` | API base URL | `http://localhost:3000` |

## Build & Installation

```bash
# Development build
make build

# Install to ~/.local/bin
make install

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

## Testing the System

### Prerequisites
1. Ensure backend is running: `cd backend && bun run dev`
2. Ensure frontend is running: `cd frontend && bun run dev`
3. Create a user account in the web app

### Step-by-Step Testing

#### 1. Generate Marcel Token (Frontend)
- [ ] Navigate to Settings page in web app
- [ ] Scroll to "Marcel CLI Token" section
- [ ] Click "Generate Token" button
- [ ] Verify token appears (starts with `marcel_`)
- [ ] Click "Copy Token" button
- [ ] Verify "Copied!" feedback appears

#### 2. Configure CLI
**Option A: Environment Variable**
```bash
export MARCEL_TOKEN="paste_your_token_here"
export MARCEL_API_ENDPOINT="http://localhost:3000"
```

**Option B: Config File**
```bash
cat > ~/.marcel.yml << EOF
api_endpoint: http://localhost:3000
auth_token: paste_your_token_here
EOF
```

#### 3. Build and Run CLI
```bash
cd cli
make build
./build/marcel
```

#### 4. Verify Functionality
- [ ] CLI launches without errors
- [ ] Authentication succeeds
- [ ] Quests appear organized by journeys
- [ ] Navigate with ↑/↓ or j/k
- [ ] Select a quest and press Space
- [ ] Quest status toggles (☐ ↔ ☑)
- [ ] Success message appears
- [ ] Refresh web app - verify quest status changed
- [ ] Press 'r' in CLI to refresh
- [ ] Quest list updates from server
- [ ] Press '?' to view help
- [ ] Help screen appears with instructions

#### 5. Error Handling Tests
- [ ] Exit CLI and unset MARCEL_TOKEN
- [ ] Run CLI - should show error screen with instructions
- [ ] Set invalid token - should show auth error
- [ ] With valid token but backend down - should show connection error
- [ ] Press 'r' on error screen - should retry connection

### Expected Behavior

**Authentication Flow:**
1. CLI reads token from env or config
2. Sends GET /user/me with Authorization header
3. Backend validates token via hybridAuthMiddleware
4. Returns user data
5. CLI fetches quests via GET /quest
6. Displays quests in TUI

**Quest Toggle Flow:**
1. User selects quest and presses Space
2. CLI sends PUT /quest/:id with done=true/false
3. Backend updates quest in database
4. Returns updated quest
5. CLI updates local state
6. Displays success message

## Future Enhancements (Not Implemented)

These are in the API client but not yet in the UI:
- Quest creation from CLI
- Quest deletion from CLI
- Quest editing (title, notes, difficulty)
- Journey management
- Filtering quests by status
- Search functionality

## Technical Notes

- **No Local Storage:** All data fetched from API in real-time
- **Stateless:** Each run authenticates and fetches fresh data
- **Token Format:** Must start with `marcel_` (validated by backend)
- **API Client:** 30-second timeout for HTTP requests
- **Error Recovery:** User can retry on any error screen

## Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `gopkg.in/yaml.v3` - YAML config parsing
- Go 1.25.0

## Color Scheme

Matches Marcel branding:
- Primary: `#5856D6` (Marcel purple)
- Success: `#34C759` (green)
- Error: `#FF3B30` (red)
- Muted: `#666666` / `#888888` (grays)
- Completed: Strikethrough gray text

## Summary

The Marcel CLI is now fully functional with:
- ✅ Secure token-based authentication
- ✅ Real-time quest viewing
- ✅ Quest completion toggling
- ✅ Journey-based organization
- ✅ Environment variable support
- ✅ Error handling and recovery
- ✅ Clean, intuitive UI

Users can now manage their Marcel quests from the terminal with the same experience as the web app!
