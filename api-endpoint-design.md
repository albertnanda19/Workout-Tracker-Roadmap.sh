# API Endpoint Design

## Auth

```http
POST /auth/signup
POST /auth/login
POST /auth/logout
```

## Exercises

```http
GET /api/v1/exercises
GET /api/v1/exercises?category=strength
```

## Workout Plans

```http
POST /api/v1/workouts
GET  /api/v1/workouts
GET  /api/v1/workouts/:id
PUT  /api/v1/workouts/:id
DELETE /api/v1/workouts/:id
```

### Schedule

```http
POST /api/v1/workouts/:id/schedule
GET  /api/v1/schedules
```

### Session

```http
POST /api/v1/sessions/start
POST /api/v1/sessions/:id/finish
GET  /api/v1/sessions
```

### Reports

```http
GET /api/v1/reports/summary
GET /api/v1/reports/progress?exercise_id=
```
