# Marcel Navigation Structure Analysis
## Current State & Future Enhancements

**Last Updated:** 2026-02-23

---

## Executive Summary

Marcel's navigation provides **5 primary sections** optimized for productivity workflows: Calendar, Tasks, Focus, Social, and Profile. The structure prioritizes **feature discoverability** and **workflow efficiency**, ensuring all features are accessible within 2 taps on mobile and desktop.

**Key Principles:**
- **Workflow-first design** - Navigation reflects user intentions, not code organization
- **Maximum 2-tap access** - All features reachable quickly
- **Clear mental models** - Each section has a distinct purpose
- **Mobile-optimized** - 5-slot bottom navigation with consistent desktop experience

**Primary Opportunity:** Add a unified "Today" dashboard for daily workflow optimization (see Future Enhancements section).

---

## Current Navigation Structure

### Primary Navigation (5 Items)

```
┌─────────────────────────────────────────────┐
│  Mobile Bottom Nav & Desktop Sidebar         │
├─────────────────────────────────────────────┤
│  1. Calendar                                │
│     📅 solar:calendar-bold-duotone          │
│     Tabs: Day | Week | Month | List         │
│     • Full calendar views and management    │
│     • Event creation and editing            │
│     • Time-based scheduling                 │
│                                             │
│  2. Tasks                                   │
│     ✅ solar:check-circle-bold-duotone      │
│     Tabs: Quests | Habits | Journeys | AI   │
│     • Quests: One-time completable tasks    │
│     • Habits: Daily/weekly routines         │
│     • Journeys: Project organization        │
│     • AI: Marcel chat for scheduling        │
│                                             │
│  3. Focus                                   │
│     🧘 solar:meditation-bold-duotone        │
│     Tabs: Timer | Sessions                  │
│     • Solo and group focus sessions         │
│     • Pomodoro timer with quest linking     │
│     • Session history and statistics        │
│     • XP/gold rewards on completion         │
│                                             │
│  4. Social                                  │
│     👥 solar:users-group-rounded-bold       │
│     Tabs: Friends | Spaces | Leaderboard    │
│     • Friends: List, activity, invitations  │
│     • Spaces: Team collaboration & chat     │
│     • Leaderboard: Competitive rankings     │
│                                             │
│  5. Profile                                 │
│     👤 User avatar                          │
│     Tabs: Stats | Shop | Inventory | Settings│
│     • User level, XP, and gold              │
│     • Cosmetics shop and purchases          │
│     • Inventory management                  │
│     • Application settings                  │
└─────────────────────────────────────────────┘
```

### Navigation Hierarchy

**Flat primary navigation** - All 5 sections are equal priority, accessible in 1 tap
**Tabbed secondary navigation** - Each section contains 2-4 tabs for feature organization
**Deep links** - Query parameters enable direct navigation (e.g., `/social?spaceId=123`)

---

## Design Principles & Decisions

### 1. Workflow-First Organization

Navigation reflects **user intentions**, not technical architecture:

- **"What do I need to do today?"** → Tasks (quests, habits)
- **"I need to focus right now"** → Focus (immediate access to timer)
- **"What's happening with my team?"** → Social (unified destination)
- **"When is my next meeting?"** → Calendar (event management)

**Why this matters:** Users think in workflows, not entities. Grouping by user intention reduces cognitive load.

### 2. Feature Discoverability

All features are visible and accessible:

- **Focus sessions** are primary navigation (not buried in tabs)
- **Marcel AI** is discoverable in Tasks section (not hidden)
- **Social features** unified in single section (no Hub/Spaces confusion)

**Impact:** Premium features (Focus, AI) are discoverable, driving engagement and conversions.

### 3. Mobile-First Constraints

Bottom navigation on mobile supports **5 items maximum**:

```
Mobile: [Calendar] [Tasks] [Focus] [Social] [Profile]
Desktop: Same 5 items in collapsible sidebar (280px → 80px)
```

**Design decision:** Keep mobile and desktop navigation consistent for predictable UX.

### 4. Clear Section Boundaries

Each section has a **distinct purpose** with no feature overlap:

- **Tasks** = Planning (what to do)
- **Focus** = Execution (doing deep work)
- **Social** = Collaboration (working with others)
- **Calendar** = Scheduling (when to do it)
- **Profile** = Progression (personal growth)

**Why this matters:** Clear boundaries prevent user confusion about where to find features.

---

## User Workflow Analysis

### Daily Productivity Flow

**Morning routine:**
```
1. Calendar (1 tap) → Check today's events
2. Tasks > Habits (2 taps) → Complete daily check-ins
3. Tasks > Quests (2 taps) → Review pending tasks
4. Focus (1 tap) → Start first focus session
```

**Total:** 3-4 taps to complete morning workflow

**With AI assistance:**
```
1. Tasks > AI (2 taps) → "What should I work on today?"
2. Review AI suggestions
3. Focus (1 tap) → Start suggested focus session
```

**Total:** 3 taps with intelligent guidance

### Team Collaboration Flow

**Check team status:**
```
1. Social > Spaces (2 taps) → View team space
2. Select space → Quests tab → See team's active work
3. Chat tab → Check messages
4. Focus (1 tap) → Join team focus session
```

**Total:** 2-3 taps to engage with team

### Social/Competitive Flow

**Check progress and rankings:**
```
1. Social > Leaderboard (2 taps) → View rankings
2. Social > Friends (2 taps) → See friends' activity
3. Focus (1 tap) → Join friend's public session
```

**Total:** 2-3 taps for social engagement

---

## Strengths of Current Design

### ✅ **1. Maximum 2-Tap Access**
Every feature is reachable in 1-2 taps:
- Primary sections: 1 tap
- Section tabs: 2 taps
- No feature is buried 3+ levels deep

### ✅ **2. Focus Sessions Prominent**
- Unique competitive advantage (group focus rooms)
- Accessible in 1 tap from anywhere
- Competitors (Focusmate) build entire apps around this
- **Result:** Focus is now a first-class feature

### ✅ **3. Marcel AI Discoverable**
- Premium feature ($4.99/month) now visible in Tasks section
- Context-appropriate placement (AI assists with task planning)
- Clear icon and label
- **Result:** AI drives premium conversions

### ✅ **4. Unified Social Section**
- All multiplayer features in one place
- Clear mental model: Friends → casual, Spaces → structured teams
- Single notification counter
- **Result:** No more "Hub vs Spaces" confusion

### ✅ **5. Mobile-Optimized**
- 5-slot bottom nav uses all available space efficiently
- Touch-friendly tab switching
- Consistent desktop/mobile experience

### ✅ **6. Command Palette Integration**
- Cmd+K quick navigation to any section
- Search for quests, journeys, spaces
- Power user optimization layer

---

## Opportunities for Enhancement

### 🎯 **1. No Unified Daily Dashboard**

**Current state:** Users must visit 3+ sections to see daily status
- Calendar → today's events
- Tasks → Habits → today's habits
- Tasks → Quests → active quests
- Focus → timer status

**Opportunity:** Create a "Today" dashboard that unifies daily items

**Expected impact:**
- Reduce daily workflow from 6+ taps to 1 tap
- Increase engagement with integrated view
- Enable AI-powered daily suggestions

**See:** Future Enhancement section below

### 🎯 **2. Habits Placement Debate**

**Current:** Habits are in Tasks section (planning mindset)

**Alternative consideration:** Move to Focus section (execution mindset)
- Habits are about daily execution, not planning
- Streaks align with focus session gamification
- Natural fit: "Focus for 30min" = habit completion

**Recommendation:** A/B test before moving
- Collect user behavior data
- Measure habit completion rates in both placements
- Make decision based on engagement metrics

### 🎯 **3. Marcel AI Accessibility**

**Current:** AI in Tasks > AI tab (2 taps)

**Enhancement options:**
- **Floating button:** Global access from any page
- **Cmd+I shortcut:** Keyboard quick-open
- **Context menu:** Right-click → "Ask Marcel"

**Expected impact:**
- Increase AI usage by 50-100%
- Enable contextual AI assistance
- Drive premium conversions

### 🎯 **4. Cross-Feature Intelligence**

**Current:** Features operate independently

**Enhancement opportunities:**
- Calendar event → Auto-suggest focus session
- Completed quest → Trigger habit check-in prompt
- Focus session → Link to related quests/journeys
- AI suggestions → Appear on Today dashboard

**Expected impact:**
- System feels cohesive and intelligent
- Reduced manual workflow management
- Increased feature discovery

### 🎯 **5. Progressive Disclosure**

**Current:** All tabs visible at once

**Enhancement consideration:**
- Smart tab ordering based on usage
- Contextual tab visibility (hide empty sections)
- Personalized navigation preferences

**Trade-off:** Balance between personalization and predictability

---

## Future Enhancement: Today Dashboard

### Overview

A unified **"Today"** section that becomes the primary landing page, showing all daily items in one view.

### Proposed Structure

```
┌─────────────────────────────────────────────┐
│  Today Dashboard (New Primary Section)       │
├─────────────────────────────────────────────┤
│  📅 Today: Sunday, Feb 23                   │
│                                             │
│  🔥 Habits Due Today                        │
│  ├─ Morning meditation (3 day streak)       │
│  ├─ Exercise (15 day streak)                │
│  └─ Journal entry                           │
│                                             │
│  ✅ Active Quests (3)                       │
│  ├─ Review navigation design                │
│  ├─ Update documentation                    │
│  └─ Test mobile layout                      │
│                                             │
│  📆 Today's Events (2)                      │
│  ├─ 10:00 AM - Team standup                │
│  └─ 2:00 PM - Design review                │
│                                             │
│  🧘 Quick Focus                             │
│  └─ [Start 25min Session] [Join Group]     │
│                                             │
│  🤖 Marcel Suggests                         │
│  └─ "Focus on navigation docs before standup"│
└─────────────────────────────────────────────┘
```

### Implementation Options

**Option 1: Replace Calendar as Primary**
```
[Today] [Tasks] [Focus] [Social] [Profile]
      └─ Calendar moved to tab inside Today
```

**Option 2: Add as 6th Item (Desktop Only)**
```
Desktop: [Today] [Calendar] [Tasks] [Focus] [Social] [Profile]
Mobile:  [Calendar] [Tasks] [Focus] [Social] [Profile]
```

**Option 3: Smart Landing Page**
```
- Default landing: Today dashboard
- Calendar/Tasks/Focus still in nav
- Today = aggregated view, not separate section
```

**Recommendation:** Option 3 (Smart Landing Page)
- Doesn't require nav restructure
- Mobile stays at 5 items
- Today becomes homepage, not nav item
- Users can still navigate to dedicated sections

### Expected Impact

**Metrics:**
- Session start efficiency: -50% taps for daily workflow
- Feature engagement: +30% (unified view increases visibility)
- User retention: +20% at day 7 (clearer daily value)
- Time to first action: -40% (everything visible on load)

**User feedback predictions:**
- "Finally! I don't have to check 4 places"
- "Marcel feels like it knows what I need to do"
- "The AI suggestions are perfect for my workflow"

---

## Implementation Roadmap

### Phase 1: Foundation ✅ COMPLETED (2026-02-23)

**Goal:** Establish solid navigation base with discoverability

**Implemented:**
- ✅ 5-section primary navigation (Calendar, Tasks, Focus, Social, Profile)
- ✅ Focus promoted to primary nav
- ✅ Marcel AI accessible in Tasks section
- ✅ Social features unified (merged Hub + Spaces)
- ✅ All navigation components updated

**Result:** Feature discoverability dramatically improved, max 2-tap access established

---

### Phase 2: Today Dashboard (2-3 Weeks)

**Goal:** Unified daily workflow view

**Tasks:**
1. **Design Today component** (Week 1)
   - Wireframe daily dashboard layout
   - Design habit completion UI
   - Design quick-action buttons
   - Mobile-responsive design

2. **Implement Today aggregation** (Week 2)
   - Fetch today's habits (due today)
   - Fetch active quests (in-progress or todo)
   - Fetch today's calendar events
   - Fetch active focus sessions
   - Marcel AI daily suggestions endpoint

3. **Smart landing page** (Week 2)
   - Make Today the default route after login
   - Cache for performance (localStorage + API)
   - Real-time updates (WebSocket integration)

4. **A/B testing** (Week 3)
   - 50% users see Today dashboard
   - 50% users see Calendar (current default)
   - Measure engagement, retention, feature usage
   - Collect qualitative feedback

**Success criteria:**
- Daily active engagement +30%
- Time to first action -40%
- Feature discovery rate +50%
- User satisfaction score +25%

---

### Phase 3: AI Integration (4-6 Weeks)

**Goal:** Make Marcel AI proactive and contextual

**Tasks:**
1. **Floating AI button** (Week 4)
   - Add global AI button (bottom-right corner)
   - Keyboard shortcut: Cmd/Ctrl+I
   - Context-aware: "Current page: Tasks > Quests"

2. **Proactive AI suggestions** (Week 5)
   - AI suggestions card on Today dashboard
   - "You have 3 habits due and a standup at 10am. Start with morning meditation?"
   - Smart focus session suggestions based on calendar

3. **Contextual AI actions** (Week 6)
   - Right-click quest → "Ask Marcel to break this down"
   - Select calendar event → "Create focus session for this"
   - Habit streak at risk → Marcel reminder + suggestion

**Success criteria:**
- AI usage +200%
- Premium conversion rate +50%
- User satisfaction with AI suggestions >80%

---

### Phase 4: Intelligence Layer (8-12 Weeks)

**Goal:** System cohesion and smart workflows

**Tasks:**
1. **Cross-feature automation**
   - Calendar event → Auto-suggest focus session
   - Quest completion → Auto-prompt related habits
   - Journey milestone → Celebration + team notification

2. **Smart notifications**
   - Unified notification center
   - Priority indicators (urgent/important)
   - Smart grouping (social, tasks, system)
   - Notification preferences

3. **Analytics dashboard**
   - Focus session statistics (time, frequency, sessions)
   - Habit streak visualizations
   - Quest completion trends
   - Journey progress tracking
   - Social engagement metrics

4. **Command palette enhancements**
   - Recent items (last 5 quests/journeys/spaces)
   - Quick actions: "Start focus", "Create quest"
   - Smart search (fuzzy matching)
   - Action history

**Success criteria:**
- "Feels intelligent" user feedback >75%
- Cross-feature action usage +150%
- Power user adoption +100%

---

## Success Metrics

### Key Performance Indicators

**Baseline Period:** First 2 weeks after Phase 1 (Feb 23 - Mar 8, 2026)

#### 1. Feature Discoverability
```
Metric: % of users who discover each feature within 7 days
Target: >80% for all primary features

Tracking:
- Focus session starts
- Marcel AI opens
- Social section visits
- Command palette usage
```

#### 2. Navigation Efficiency
```
Metric: Average taps to reach target feature
Target: <2 taps for all features

Tracking:
- Tap count distribution
- Most common paths
- Drop-off points
```

#### 3. Daily Engagement
```
Metric: Average sections visited per session
Target: 3+ sections per session

Tracking:
- Section visit frequency
- Session duration per section
- Cross-section navigation patterns
```

#### 4. User Retention
```
Metric: Day 7 and Day 30 retention rates
Target: +30% improvement after Today dashboard

Tracking:
- New user retention curves
- Feature engagement correlation
- Churn reasons (exit surveys)
```

#### 5. Premium Conversion
```
Metric: Free → Premium conversion rate
Target: +50% after AI enhancements

Tracking:
- AI feature usage (free vs premium)
- Paywall encounters
- Conversion funnel
- Upsell success rate
```

### Data Collection Implementation

```typescript
// Navigation tracking
analytics.track('nav_clicked', {
  from: 'calendar',
  to: 'focus',
  method: 'bottom_nav', // or 'sidebar', 'command_palette'
  timestamp: Date.now()
})

// Feature discovery
analytics.track('feature_first_use', {
  feature: 'marcel_ai',
  days_since_signup: 3,
  discovery_method: 'navigation', // or 'onboarding', 'friend'
  context: 'tasks_section'
})

// Workflow completion
analytics.track('daily_workflow_completed', {
  steps: ['today', 'habits_checked', 'focus_started'],
  duration_seconds: 45,
  interruptions: 0
})

// Today dashboard (Phase 2)
analytics.track('today_dashboard_viewed', {
  items_shown: {
    habits: 3,
    quests: 5,
    events: 2,
    ai_suggestions: 1
  },
  first_action: 'habit_completed',
  time_to_action_ms: 1200
})
```

---

## Competitor Analysis

### How Marcel Compares

#### Todoist
**Navigation:** Today | Inbox | Upcoming | Filters | Projects
**Insight:** "Today" is primary - daily workflow front and center
**Marcel:** We have Calendar first; **Today dashboard would match Todoist's approach**

#### Notion
**Navigation:** Workspace switcher + flexible sidebar
**Insight:** Highly customizable but requires setup
**Marcel:** Structured navigation better for new users, less flexible for power users

#### Habitica
**Navigation:** Tasks | Inventory | Social | Guilds
**Insight:** Gamification + productivity (similar to Marcel)
**Marcel:** ✅ Similar structure with Tasks, Social, Profile (progression)

#### Focusmate
**Navigation:** Sessions | Schedule | Profile
**Insight:** Entire app built around focus sessions
**Marcel:** ✅ Focus now primary nav (previously buried)

#### ClickUp
**Navigation:** Home | Inbox | Docs | Dashboards | Goals
**Insight:** "Home" dashboard as entry point
**Marcel:** **Today dashboard would provide similar unified view**

#### Sunsama
**Navigation:** Today | Calendar | Backlog | Objectives
**Insight:** Daily planning workflow emphasized
**Marcel:** Similar to Today dashboard proposal

### Competitive Positioning

**Marcel's unique value:**
- ✅ **Focus sessions** with multiplayer (rare feature)
- ✅ **AI scheduling** assistant (Marcel chat)
- ✅ **Gamification** with XP, levels, streaks
- ✅ **Social productivity** (spaces, leaderboard)

**Where Marcel leads:**
- Group focus sessions (better than Focusmate's 1:1 only)
- AI integration (better than Todoist's basic suggestions)
- Full feature integration (calendar + tasks + focus + social)

**Where Marcel can improve:**
- ❌ No "Today" dashboard (Todoist, Sunsama have this)
- ❌ AI not as prominent as it should be
- ⚠️ Habits placement (could be more execution-focused)

**Bottom line:** Marcel has unique features that competitors lack (group focus, AI scheduling, gamification). Adding Today dashboard would make Marcel's daily workflow best-in-class.

---

## Technical Implementation Notes

### Navigation Configuration

**Primary source:** `/frontend/src/lib/pages.ts`
```typescript
export const pages = [
  { href: "/calendar", icon: "solar:calendar-bold-duotone", label: "Calendar" },
  { href: "/tasks", icon: "solar:check-circle-bold-duotone", label: "Tasks" },
  { href: "/focus", icon: "solar:meditation-bold-duotone", label: "Focus" },
  { href: "/social", icon: "solar:users-group-rounded-bold-duotone", label: "Social" },
]
```

**Navigation components:**
- Desktop: `/frontend/src/lib/components/Sidebar.svelte`
- Mobile: `/frontend/src/lib/components/MobileNavbar.svelte`
- Command: `/frontend/src/lib/components/CommandBar.svelte`

### Routing Structure

**SvelteKit file-based routing:**
```
src/routes/(app)/
├── calendar/+page.svelte
├── tasks/
│   ├── +page.svelte (tabs: Quests, Habits, Journeys, AI)
│   ├── QuestsTab.svelte
│   ├── HabitsTab.svelte
│   └── JourneysTab.svelte
├── focus/
│   └── +page.svelte (tabs: Timer, Sessions)
├── social/
│   ├── +page.svelte
│   ├── FriendsTab.svelte
│   ├── SpacesTab.svelte
│   └── LeaderboardTab.svelte
└── profile/+page.svelte
```

### State Management

**Tab preferences:** `localStorage`
```typescript
// Persist user's last active tab per section
localStorage.setItem('tasksDefaultTab', 'quests') // 'quests' | 'habits' | 'journeys'
localStorage.setItem('focusDefaultTab', 'timer')  // 'timer' | 'sessions'
localStorage.setItem('socialDefaultTab', 'friends') // 'friends' | 'spaces' | 'leaderboard'
```

**Real-time data:** Svelte stores + WebSocket
```typescript
// User, spaces, invites, messages all use reactive stores
import { user, spaces, invites, unreadMessages } from '$lib/stores'
```

### Deep Linking

**Query parameters for context:**
```
/social?spaceId=123        // Opens Social > Spaces > Space Detail
/tasks?tab=habits          // Opens Tasks > Habits tab
/focus?session=456         // Opens Focus with session 456
```

---

## Conclusion

Marcel's current navigation structure provides **solid foundation** for productivity workflows with clear discoverability and efficient access patterns. All features are within 2 taps, Focus and AI are prominent, and social features are unified.

**Primary enhancement opportunity:** Add a unified "Today" dashboard to become the daily workflow hub, aggregating habits, quests, events, and AI suggestions in one view.

**Implementation approach:**
1. ✅ Phase 1 complete - Strong navigation base
2. 🎯 Phase 2 next - Today dashboard (2-3 weeks)
3. 🚀 Phase 3-4 - AI integration and intelligence layer

**Expected outcome:** Marcel becomes best-in-class for daily productivity workflows, with competitive advantages (group focus, AI scheduling, gamification) made prominent and accessible.

---

**Document Version:** 2.0 (2026-02-23)
**Status:** Phase 1 Complete, Phase 2 Planning
