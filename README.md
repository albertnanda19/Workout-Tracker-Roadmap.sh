# Workout Tracker API

REST API for managing workout plans, schedules, and user authentication.

## Requirements

- Go 1.22+
- PostgreSQL

## Configuration

Create `.env` in project root:

```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=workout_tracker
DB_SSLMODE=disable
JWT_SECRET=supersecretkey
```

## Run

```
go run cmd/api/main.go
```

## Test

```
go test ./...
```

## API Documentation (OpenAPI)

OpenAPI spec file:

- `docs/openapi.yaml`

### View Documentation

- Open https://editor.swagger.io
- Upload atau paste isi `docs/openapi.yaml`

## Curl Examples

### Register

```
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"user@example.com","password":"password123"}'
```

### Login

```
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Create Workout Plan

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

### List Workouts (Pagination + Filter)

```
curl -X GET "http://localhost:8080/api/workouts?page=1&limit=10&name=push" \
  -H "Authorization: Bearer <TOKEN>"
```

### Schedule Workout

```
curl -X POST http://localhost:8080/api/workouts/schedule \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"workout_plan_id":"22222222-2222-2222-2222-222222222222","scheduled_date":"2026-02-20"}'
```
