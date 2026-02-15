package response

import (
	"time"

	"workout-tracker/internal/domain"
)

type WorkoutPlanDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ScheduledWorkoutDTO struct {
	ID            string    `json:"id"`
	WorkoutPlanID string    `json:"workout_plan_id"`
	ScheduledDate time.Time `json:"scheduled_date"`
	CreatedAt     time.Time `json:"created_at"`
}

func ToWorkoutPlanDTO(p domain.WorkoutPlan) WorkoutPlanDTO {
	return WorkoutPlanDTO{
		ID:        p.ID,
		Name:      p.Name,
		Notes:     p.Notes,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToScheduledWorkoutDTO(sw domain.ScheduledWorkout) ScheduledWorkoutDTO {
	return ScheduledWorkoutDTO{
		ID:            sw.ID,
		WorkoutPlanID: sw.WorkoutPlanID,
		ScheduledDate: sw.ScheduledDate,
		CreatedAt:     sw.CreatedAt,
	}
}
