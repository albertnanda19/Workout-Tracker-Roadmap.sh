package domain

import "time"

type ScheduledWorkout struct {
	ID            string
	UserID        string
	WorkoutPlanID string
	ScheduledDate time.Time
	CreatedAt     time.Time
}
