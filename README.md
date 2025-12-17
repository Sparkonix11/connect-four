# Connect Four - Real-time Multiplayer Game

A full-stack implementation of the classic Connect Four game with real-time multiplayer support, intelligent bot opponent, and game analytics.

## What's This?

This is a real-time version of Connect Four (or 4 in a Row) where you can play against another player or challenge an AI bot. The game uses WebSockets for instant updates, so you'll see your opponent's moves in real time.

If no one joins your game within 10 seconds, the bot will jump in to play with you. And yes, the bot actually tries to win - it's not just making random moves.

## Tech Stack

**Backend:**
- Go (Golang) - chose this for better concurrency and performance
- Gorilla WebSocket - real-time communication
- PostgreSQL - game history and leaderboard persistence
- Kafka - analytics event streaming (optional)

**Frontend:**
- React with TypeScript
- Vite - fast dev server and builds
- Tailwind CSS - styling
- WebSocket client for real-time updates

**Infrastructure:**
- Docker & Docker Compose
- Air - live reload for Go development

## Features

Real-time multiplayer gameplay via WebSockets  
Smart bot opponent with strategic decision-making  
Player matchmaking with 10-second timeout  
Reconnection support (30 seconds to rejoin)  
Leaderboard tracking wins  
Game analytics via Kafka  
Persistent game history  

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local frontend development)

### Quick Start with Docker

1. Clone the repo and set up environment:

```bash
git clone https://github.com/Sparkonix11/connect-four.git
cd connect-four
cp .env.example .env
```

2. Start everything (without Kafka):

```bash
docker compose up
```

The app will be available at:
- Frontend: http://localhost
- Backend API: http://localhost:8080
- Health check: http://localhost:8080/health

3. If you want to enable analytics with Kafka:

```bash
docker compose --profile kafka up
```

That's it. The database will auto-migrate on startup.

### Local Development

**Backend:**

```bash
# Install Air for hot reload
go install github.com/air-verse/air@latest

# Run with hot reload
air
```

**Frontend:**

```bash
cd frontend
npm install
npm run dev
```

Frontend will run on http://localhost:5173

### Environment Variables

Check `.env.example` for all available options. Key ones:

- `SERVER_PORT` - API server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `KAFKA_ENABLED` - Enable/disable Kafka analytics (true/false)
- `MATCHMAKING_TIMEOUT_SECONDS` - Wait time before bot joins (default: 10)
- `RECONNECT_TIMEOUT_SECONDS` - Time to rejoin after disconnect (default: 30)
- `BOT_MOVE_DELAY_MS` - Bot thinking time for realism (default: 300ms)

## How to Play

1. Open the app and enter your username
2. You'll be matched with another player (or bot after 10s)
3. Click on any column to drop your disc
4. First to connect 4 wins!

If you disconnect, you can rejoin the same game within 30 seconds by entering the same username.

## Project Structure

```
.
├── cmd/server/          # Application entry point
├── internal/
│   ├── api/            # HTTP routes
│   ├── bot/            # AI bot strategy
│   ├── database/       # Database connection
│   ├── game/           # Core game logic
│   ├── kafka/          # Kafka producer/consumer
│   ├── matchmaking/    # Player queue
│   ├── models/         # Data models
│   ├── repository/     # Database queries
│   └── websocket/      # WebSocket handlers
├── frontend/           # React app
├── api/               # OpenAPI spec
└── docker-compose.yml
```

## About the Bot

The bot isn't just making random moves. It follows this priority:

1. **Win** - if it can connect 4, it will
2. **Block** - if you're about to win, it'll try to stop you
3. **Build** - otherwise, it works on creating its own winning positions

Check out `internal/bot/strategy.go` if you want to see how it thinks.

## Kafka Analytics

When enabled, the system tracks these events via Kafka:

- Game start/end
- Player moves
- Game outcomes (win/loss/draw)
- Player disconnections

A consumer service processes these events to calculate:
- Average game duration
- Win rates
- Games per hour/day
- Player-specific stats

See `KAFKA_SETUP.md` for detailed Kafka configuration.

## API Endpoints

- `GET /health` - Health check
- `GET /api/leaderboard` - Get top players
- `WS /ws` - WebSocket connection for gameplay

## Testing

Backend tests:

```bash
go test ./...
```

Frontend type checking:

```bash
cd frontend
npm run typecheck
```