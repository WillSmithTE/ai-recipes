Go backend for AI-Recipes

## Features

- **Dashboard** (`/dashboard`) — Central hub for viewing and managing AI recipes.
- **Feedback Widget** — Integrated in-app feedback widget for collecting user feedback.

## Endpoints

### GET /api/dashboard
Returns dashboard data including recipe stats and feedback summary.

### POST /api/feedback
Submit user feedback with rating (1-4), comment, and page.

## Running

```bash
go run main.go
```

## Testing

```bash
go test -v ./...
```
