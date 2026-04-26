# Scrum Poker

A minimal Persian RTL Scrum Poker app with live team voting.

## Tech Stack

- Frontend: React, TypeScript, Vite, Tailwind CSS
- Backend: Go, Fiber
- Database: MySQL
- Realtime: Redis, WebSocket

## Run with Docker

```bash
docker compose up --build
```

## Backend (local)

```bash
cd backend
cp .env.example .env
go run ./cmd/server
```

## Frontend (local)

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

## Test and Build

### Backend

```bash
cd backend
go test ./...
go build ./...
```

### Frontend

```bash
cd frontend
npm test
npm run build
```
