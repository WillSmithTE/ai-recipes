Go backend for AI-Recipes

## Features

- **Dashboard** — Main overview page providing a summary of AI recipe data and user activity (`/dashboard`)
- **Feedback Widget** — Embedded widget for collecting user feedback directly from the dashboard

## API Endpoints

| Method | Path         | Description                        |
|--------|--------------|------------------------------------|
| GET    | `/dashboard` | Returns dashboard summary data     |
| POST   | `/feedback`  | Submits user feedback from widget  |

## Getting Started

1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Start the server: `go run main.go`
4. Visit `http://localhost:8080/dashboard`
