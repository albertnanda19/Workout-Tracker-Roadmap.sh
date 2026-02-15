BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS exercises (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR NOT NULL,
  description TEXT,
  category VARCHAR,
  muscle_group VARCHAR,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS workout_plans (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  name VARCHAR NOT NULL,
  notes TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT workout_plans_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS workout_plan_exercises (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  workout_plan_id UUID NOT NULL,
  exercise_id UUID NOT NULL,
  sets INTEGER NOT NULL,
  reps INTEGER NOT NULL,
  weight NUMERIC(6,2),
  order_index INTEGER NOT NULL,
  CONSTRAINT workout_plan_exercises_workout_plan_id_fkey
    FOREIGN KEY (workout_plan_id) REFERENCES workout_plans(id) ON DELETE CASCADE,
  CONSTRAINT workout_plan_exercises_exercise_id_fkey
    FOREIGN KEY (exercise_id) REFERENCES exercises(id),
  CONSTRAINT workout_plan_exercises_sets_check CHECK (sets > 0),
  CONSTRAINT workout_plan_exercises_reps_check CHECK (reps > 0),
  CONSTRAINT workout_plan_exercises_order_index_check CHECK (order_index >= 0),
  CONSTRAINT workout_plan_exercises_weight_check CHECK (weight IS NULL OR weight >= 0),
  CONSTRAINT workout_plan_exercises_plan_order_unique UNIQUE (workout_plan_id, order_index)
);

CREATE TABLE IF NOT EXISTS scheduled_workouts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  workout_plan_id UUID NOT NULL,
  scheduled_at TIMESTAMP NOT NULL,
  status VARCHAR NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT scheduled_workouts_workout_plan_id_fkey
    FOREIGN KEY (workout_plan_id) REFERENCES workout_plans(id) ON DELETE CASCADE,
  CONSTRAINT scheduled_workouts_status_check
    CHECK (status IN ('pending', 'completed', 'canceled'))
);

CREATE TABLE IF NOT EXISTS workout_sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  workout_plan_id UUID,
  started_at TIMESTAMP NOT NULL,
  completed_at TIMESTAMP,
  notes TEXT,
  CONSTRAINT workout_sessions_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT workout_sessions_workout_plan_id_fkey
    FOREIGN KEY (workout_plan_id) REFERENCES workout_plans(id) ON DELETE SET NULL,
  CONSTRAINT workout_sessions_completed_after_started_check
    CHECK (completed_at IS NULL OR completed_at >= started_at)
);

CREATE TABLE IF NOT EXISTS workout_session_exercises (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  workout_session_id UUID NOT NULL,
  exercise_id UUID NOT NULL,
  sets INTEGER NOT NULL,
  reps INTEGER NOT NULL,
  weight NUMERIC(6,2),
  actual_reps INTEGER,
  actual_weight NUMERIC(6,2),
  CONSTRAINT workout_session_exercises_workout_session_id_fkey
    FOREIGN KEY (workout_session_id) REFERENCES workout_sessions(id) ON DELETE CASCADE,
  CONSTRAINT workout_session_exercises_exercise_id_fkey
    FOREIGN KEY (exercise_id) REFERENCES exercises(id),
  CONSTRAINT workout_session_exercises_sets_check CHECK (sets > 0),
  CONSTRAINT workout_session_exercises_reps_check CHECK (reps > 0),
  CONSTRAINT workout_session_exercises_weight_check CHECK (weight IS NULL OR weight >= 0),
  CONSTRAINT workout_session_exercises_actual_reps_check CHECK (actual_reps IS NULL OR actual_reps >= 0),
  CONSTRAINT workout_session_exercises_actual_weight_check CHECK (actual_weight IS NULL OR actual_weight >= 0)
);

CREATE INDEX IF NOT EXISTS idx_workout_plans_user_id ON workout_plans(user_id);

CREATE INDEX IF NOT EXISTS idx_workout_plan_exercises_workout_plan_id ON workout_plan_exercises(workout_plan_id);
CREATE INDEX IF NOT EXISTS idx_workout_plan_exercises_exercise_id ON workout_plan_exercises(exercise_id);

CREATE INDEX IF NOT EXISTS idx_scheduled_workouts_workout_plan_id ON scheduled_workouts(workout_plan_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_workouts_scheduled_at ON scheduled_workouts(scheduled_at);

CREATE INDEX IF NOT EXISTS idx_workout_sessions_user_id ON workout_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_workout_sessions_workout_plan_id ON workout_sessions(workout_plan_id);

CREATE INDEX IF NOT EXISTS idx_workout_session_exercises_workout_session_id ON workout_session_exercises(workout_session_id);
CREATE INDEX IF NOT EXISTS idx_workout_session_exercises_exercise_id ON workout_session_exercises(exercise_id);

COMMIT;
