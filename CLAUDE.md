# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Marcel is a gamified productivity application with multiple client applications (web, mobile, backoffice) and a unified backend API. The platform includes features like focus sessions, quests, habits, spaces (collaborative areas), an AI assistant (Marcel), calendar integration, and a premium subscription system.

## Monorepo Structure

```
/backend       - Bun + Hono API server with Prisma ORM
/frontend      - SvelteKit web application
/mobile        - React Native (Expo) mobile app
/backoffice    - Next.js admin dashboard
```

## Tech Stack

- **Runtime**: Bun (use `bun` commands, not `npm` or `node`)
- **Backend**: Hono framework, Prisma ORM, PostgreSQL, WebSocket support
- **Frontend**: SvelteKit with Vite, TailwindCSS, bits-ui components
- **Mobile**: Expo Router (file-based routing), React Native, Zustand for state
- **Backoffice**: Next.js 15, React 19, TailwindCSS
- **UI Icons**: Solar icon set (always use these, no other icon packs)
- **AI**: OpenAI API, Qdrant vector database for RAG

## Common Commands

### Backend (backend/)
```bash
bun install                    # Install dependencies
bun run dev                    # Run dev server with hot reload
bun run start                  # Run production server
bun run migrate                # Deploy Prisma migrations
bun run migrate-dev            # Create and apply dev migration
bun run generate               # Generate Prisma client
bun run studio                 # Open Prisma Studio
bun run build                  # Compile to standalone binary
```

### Frontend (frontend/)
```bash
bun install                    # Install dependencies
bun run dev                    # Run dev server (port 5173)
bun run build                  # Build for production
bun run preview                # Preview production build
bun run check                  # Type-check with svelte-check
```

### Mobile (mobile/)
```bash
bun install                    # Install dependencies
bun run start                  # Start Expo dev server
bun run android                # Run on Android
bun run ios                    # Run on iOS
bun run web                    # Run web version
bun run test                   # Run Jest tests
```

### Backoffice (backoffice/)
```bash
bun install                    # Install dependencies
bun run dev                    # Run dev server (port 3003)
bun run build                  # Build for production
bun run start                  # Start production server
```

## Architecture Patterns

### Backend Architecture

**Route Structure** (`backend/src/routes/`):
- Each domain has its own route file (auth.ts, user.ts, quest.ts, space.ts, etc.)
- Routes are organized as Hono sub-applications
- All routes are imported and mounted in `webserver.ts`

**Services** (`backend/src/services/`):
- AI-related services (marcel-chat, marcel-tools, embedding, ai-context, ai-schedule)
- Services contain business logic separate from route handlers
- Marcel AI uses OpenAI function calling with custom tools

**WebSocket** (`backend/src/websocket.ts`):
- Real-time updates for focus sessions, spaces, messages
- Uses Bun's native WebSocket support
- Connected clients tracked in memory maps

**Key Backend Files**:
- `src/index.ts` - Entry point, binds Hono app with WebSocket
- `src/webserver.ts` - Main Hono app, CORS config, route mounting
- `src/lib/jwt.ts` - JWT authentication utilities
- `src/lib/timer-manager.ts` - Focus session timer management
- `src/lib/stripe.ts` - Stripe payment integration
- `src/lib/qdrant.ts` - Vector database for AI context

**Prisma**:
- Schema at `backend/prisma/schema.prisma`
- Client generated to `backend/src/lib/prisma/client`
- Run migrations with `bun run migrate-dev` before making schema changes

### Frontend Architecture (SvelteKit)

**Routing** (`frontend/src/routes/`):
- File-based routing with SvelteKit conventions
- `(app)/` group for authenticated routes
- `+page.svelte` files for routes, `+layout.svelte` for layouts

**State Management** (`frontend/src/lib/stores/`):
- Custom `storage()` wrapper for persisted stores
- Domain stores: user.ts, quests.ts, spaces.ts, habits.ts, events.ts, etc.
- Stores expose fetch functions and reactive Svelte stores

**API Communication** (`frontend/src/lib/api.ts`):
- `requestApi()` centralized API client with JWT auth
- Automatically adds auth headers from localStorage

**WebSocket** (`frontend/src/lib/websockets.ts`):
- Real-time connection to backend WebSocket
- Message handlers in `websocketMessageHandler.ts`

**UI Components** (`frontend/src/lib/components/ui/`):
- Built with bits-ui (Svelte headless components)
- TailwindCSS for styling
- Use Solar icons only

### Mobile Architecture (React Native/Expo)

**Routing** (`mobile/app/`):
- Expo Router with file-based routing
- `(auth)/` group for auth screens
- `app/` group for main app screens
- `_layout.tsx` files define navigation structure

**State Management** (`mobile/stores/`):
- Zustand stores for each domain (user.ts, quests.ts, spaces.ts, etc.)
- Stores imported from Expo's AsyncStorage for persistence

**Navigation**:
- Tab navigation for main app screens
- Stack navigation for auth flow
- Expo Router handles deep linking

### Backoffice Architecture (Next.js)

**Structure** (`backoffice/app/`):
- App Router (Next.js 13+)
- Server and client components
- Simple admin interface for user management

## Database Schema (Key Models)

**User**: Core user model with gamification (level, xp, gold, gems), customization (hat, glasses, shape, color), premium subscription
**Space**: Collaborative workspaces with members, messages, invite links
**Quest**: Tasks with XP/rewards, can be assigned, part of journeys
**Journey**: Collections of quests (like projects)
**Habit**: Recurring tasks with streak tracking
**Event**: Calendar events with Google Calendar sync support
**FocusSession**: Pomodoro-style focus sessions, can be collaborative
**Item**: Shop items for avatar customization
**Message**: Chat messages in spaces
**Conversation**: AI chat conversations with Marcel
**PendingAction**: AI-suggested actions requiring user approval
**AIScheduleSuggestion**: AI-generated schedule suggestions
**Subscription**: Stripe subscription tracking
**GoogleCalendarToken**: OAuth tokens for calendar sync

## Key Features

### Marcel AI Assistant
- Context-aware AI using RAG (Qdrant vector DB stores user context)
- OpenAI function calling for tool use (create quests, schedule events, etc.)
- Conversations stored in database
- Pending actions require user approval before execution
- Context includes user's quests, habits, events, and past conversations

### Focus Sessions
- Collaborative Pomodoro sessions
- Real-time participant tracking via WebSocket
- Timer managed on backend (`timer-manager.ts`)
- Participants can join active sessions

### Spaces (Collaborative Workspaces)
- Real-time chat via WebSocket
- Members can join via invite links
- Shared quests and focus sessions

### Premium Subscription
- Stripe integration for payments
- Premium features gated by `isPremium` flag
- Premium middleware on protected routes
- Google Calendar sync is premium feature

### Gamification System
- XP and levels with rewards
- Virtual currencies: gold and gems
- Avatar customization (hats, glasses, shapes, colors, backgrounds)
- Shop for purchasing items
- Quest completion rewards

## Design Guidelines

- Always use Solar icon pack (imported from `solar-icon-set` in backoffice, use appropriate icon libraries in other apps)
- No gradients in designs
- No glow effects in designs
- Clean, modern UI with TailwindCSS

## Development Workflow

1. **Starting Development**:
   - Run backend: `cd backend && bun run dev`
   - Run frontend: `cd frontend && bun run dev`
   - Run mobile: `cd mobile && bun run start`
   - Run backoffice: `cd backoffice && bun run dev`

2. **Database Changes**:
   - Edit `backend/prisma/schema.prisma`
   - Run `cd backend && bun run migrate-dev` to create migration
   - Migration automatically runs `generate` to update Prisma client
   - Commit migration files to git

3. **API Changes**:
   - Add/modify routes in `backend/src/routes/`
   - Complex logic goes in `backend/src/services/`
   - Update types in `backend/src/lib/types.ts`
   - Ensure CORS origins in `webserver.ts` include your dev environment

4. **Frontend/Mobile Changes**:
   - Update stores in respective `/stores/` directories
   - API calls use centralized `requestApi()` or similar utilities
   - WebSocket updates handled by message handlers

## Environment Variables

Backend requires:
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret for JWT signing
- `OPENAI_API_KEY` - OpenAI API key for Marcel AI
- `QDRANT_URL`, `QDRANT_API_KEY` - Qdrant vector DB
- `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET` - Stripe integration
- Google OAuth: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`

Frontends need:
- API URL configuration (check respective `apiUrl.ts` files)

## Testing

- Mobile: `bun run test` (Jest with jest-expo preset)
- Backend: Use `bun test` for Bun's built-in test runner
- No comprehensive test suite currently exists

## Common Gotchas

- Bun is used throughout - don't use `npm` or `yarn` commands
- Prisma client is generated to custom path: `backend/src/lib/prisma/client`
- Backend and frontend both have separate git repositories (submodules) - check for `.git` directories
- WebSocket connections must match CORS allowed origins
- Mobile uses Expo Router v6 (file-based routing)
- Frontend uses SvelteKit (not Svelte standalone)
