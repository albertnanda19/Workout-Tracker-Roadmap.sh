# Workout Tracker API

## Project Reference

This project is built as part of the roadmap.sh Backend Project:

https://roadmap.sh/projects/fitness-workout-tracker

## Overview

Workout Tracker API is a REST API for managing workout plans, workout schedules, and user authentication.

Main capabilities:

- User Authentication (JWT)
- Workout Plan CRUD
- Scheduled Workouts
- Pagination & Filtering
- OpenAPI Documentation
- Unit Testing
- Observability (Logging + Error Handling)

## Tech Stack

- Go (Golang)
- PostgreSQL
- Clean Architecture
- JWT Authentication
- OpenAPI 3.0 Specification
- Testify (Unit Testing)
- slog (Structured Logging)

## Requirements

- Go 1.22+
- PostgreSQL

## Running the Project

1. Clone repository

2. Copy environment file

```
cp .env.example .env
```

3. Set environment variables in `.env` (PostgreSQL connection and JWT secret)

4. Run the API server

```
go run cmd/api/main.go
```

## Running Tests

```
go test ./...
go test ./... -cover
```

## API Documentation

OpenAPI specification is available at:

- `docs/openapi.yaml`

To view the documentation:

1. Open https://editor.swagger.io
2. Upload or paste the contents of `docs/openapi.yaml`

## Example Usage

### 1) Register

```
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"user@example.com","password":"your-password"}'
```

### 2) Login

```
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"your-password"}'
```

### 3) Create Workout

```
curl -X POST http://localhost:8080/api/workouts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "name": "Push Day",
    "notes": "chest + triceps",
    "exercises": [
      {"exercise_id":"11111111-1111-1111-1111-111111111111","sets":3,"reps":10,"weight":60,"order_index":0}
    ]
  }'
```

### 4) List Workouts

```
curl -X GET "http://localhost:8080/api/workouts?page=1&limit=10&name=push" \
  -H "Authorization: Bearer <TOKEN>"
```

### 5) Schedule Workout

```
curl -X POST http://localhost:8080/api/workouts/schedule \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"workout_plan_id":"22222222-2222-2222-2222-222222222222","scheduled_date":"2026-02-20"}'
```

## Project Structure

```
cmd/
internal/
  domain/
  repository/
  usecase/
  delivery/
docs/
```

- **cmd/**
  Application entrypoints.
- **internal/domain/**
  Core business entities, errors, and shared types (e.g., pagination).
- **internal/repository/**
  Repository interfaces and contracts used by the use cases.
- **internal/usecase/**
  Application business logic and orchestration.
- **internal/delivery/**
  HTTP delivery layer (handlers, middleware, request/response DTO).
- **docs/**
  OpenAPI specification (`openapi.yaml`).

## License

This project is for educational purposes as part of roadmap.sh backend projects.
