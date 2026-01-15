# {{.ProjectName}}

A mobile application built with Expo React Native and Go backend, using Supabase for the database.

## Project Structure

```
{{.ProjectName}}/
├── mobile/            # Expo React Native + NativeWind
│   ├── App.tsx        # Entry point
│   └── src/
│       └── lib/       # Utilities and API client
├── backend/           # Go Echo server
│   ├── cmd/server/    # Entry point
│   └── internal/      # Private packages
│       ├── config/    # Configuration
│       └── server/    # HTTP server
└── package.json       # Root workspace config
```

## Getting Started

### Prerequisites

- Node.js 18+
- Go 1.21+
- Expo CLI (`npm install -g expo-cli`)
- iOS Simulator (Mac) or Android Emulator
- Supabase account (for production)

### Setup

1. Install mobile dependencies:
   ```bash
   npm run install:mobile
   ```

2. Copy environment files:
   ```bash
   cp backend/.env.example backend/.env
   cp mobile/.env.example mobile/.env
   ```

3. Configure environment variables in both `.env` files.

### Development

Run both mobile app and backend:
```bash
npm run dev
```

Or run them separately:
```bash
# Terminal 1 - Backend
npm run dev:backend

# Terminal 2 - Mobile
npm run dev:mobile
```

Run on specific platform:
```bash
# iOS Simulator
npm run dev:ios

# Android Emulator
npm run dev:android
```

- Mobile App: Expo Go or Simulator
- Backend API: http://localhost:8080

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| * | /api/v1/* | API routes (add your endpoints here) |

## Architecture

```
Mobile App → Expo (React Native) → Go API (port 8080) → Supabase
```

The mobile app never connects directly to Supabase. All data flows through the Go backend, which:
- Validates requests
- Handles authentication
- Manages database operations
- Returns JSON responses

## Environment Variables

### Backend (.env)

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| SUPABASE_URL | Supabase project URL | - |
| SUPABASE_KEY | Supabase anon/service key | - |

### Mobile (.env)

| Variable | Description | Default |
|----------|-------------|---------|
| API_URL | Backend API URL | http://localhost:8080 |

## Building for Production

```bash
# Build mobile app
npm run build:mobile

# Build backend
cd backend && go build -o bin/server ./cmd/server
```

## Stack

- **Mobile**: Expo SDK 50, React Native, TypeScript, NativeWind
- **Backend**: Go, Echo, godotenv
- **Database**: Supabase (PostgreSQL)
