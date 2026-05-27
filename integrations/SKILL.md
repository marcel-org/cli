---
name: marcel
description: >
  Personal task and calendar CLI. Use when the user asks to manage tasks (quests),
  habits, calendar events, focus sessions, or journeys.
---

# marcel — Personal tasks & calendar

Binary: `marcel`
Config: `~/.marcel.yml`, token in `~/.marcel.token` or `MARCEL_TOKEN` env var

## When to apply

Use when the user mentions personal tasks, quests, habits, calendar events, focus sessions, journeys, reminders, todos, deadlines, or appointments.
Triggers: "add a task", "create an event", "remind me", "calendar", "focus session", "habit", "quest", "journey"

## Commands

### Quests (tasks)
```
marcel quest list [--json]
marcel quest add <title> [--note <text>] [--difficulty easy|medium|hard] [--journey <id>]
marcel quest done <id>
marcel quest undo <id>
marcel quest toggle <id>
marcel quest update <id> [--title <text>] [--note <text>] [--difficulty <level>]
marcel quest delete <id>
```

### Events (calendar)
```
marcel event list [--json] [--all]
marcel event add <title> --date <YYYY-MM-DD> [--time <HH:MM>] [--end-date <YYYY-MM-DD>] [--end-time <HH:MM>] [--location <text>] [--description <text>]
marcel event update <id> [--title <text>] [--date ...] [--time ...]
marcel event delete <id>
```

### Habits
```
marcel habit list [--json]
marcel habit add <name> [--cycle daily|weekly]
marcel habit done <id>
marcel habit undo <id>
marcel habit update <id> [--name <text>] [--cycle <type>]
marcel habit delete <id>
```

### Focus sessions
```
marcel focus start [--name <text>] [--duration <minutes>]    Default 25 min
marcel focus pause <id>
marcel focus complete <id>
marcel focus status <id>
marcel focus history [--json]
```

### Journeys (quest groups)
```
marcel journey list [--json]
marcel journey add <name>
marcel journey update <id> --name <text>
marcel journey delete <id>
```

### Other
```
marcel auth check
marcel tui                     Launch interactive TUI
marcel version
marcel help
```

## Rules
- `--json` available on list commands and focus history
- Event dates: `YYYY-MM-DD`, times: `HH:MM`
- `event list` excludes past events by default; use `--all` for past
- Confirm before deleting
- After using marcel, tell the user what was added, changed, or found
