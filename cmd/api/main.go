package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpdelivery "workout-tracker/internal/delivery/http"
	"workout-tracker/internal/infrastructure"
	"workout-tracker/internal/infrastructure/auth"
	"workout-tracker/internal/infrastructure/migration"
	"workout-tracker/internal/infrastructure/repository"
	"workout-tracker/internal/infrastructure/seeder"
	"workout-tracker/internal/usecase"
)

func main() {
	cfg, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := infrastructure.NewPostgresDB(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database connected")

	if err := migration.RunMigrations(db); err != nil {
		log.Fatal(err)
	}
	log.Println("Migration completed")

	if err := seeder.RunSeeders(db); err != nil {
		log.Fatal(err)
	}
	log.Println("Seeding completed")

	userRepo := repository.NewPostgresUserRepository(db)
	workoutRepo := repository.NewPostgresWorkoutRepository(db)
	exerciseRepo := repository.NewPostgresExerciseRepository(db)

	jwtSvc := auth.NewJWTService(cfg.JWTSecret)
	userUC := usecase.NewUserUsecase(userRepo, jwtSvc)
	workoutUC := usecase.NewWorkoutUsecase(workoutRepo)
	exerciseUC := usecase.NewExerciseUsecase(exerciseRepo)
	handler := httpdelivery.NewHandler(userUC, workoutUC, exerciseUC)
	router := httpdelivery.NewRouter(handler, jwtSvc)
	_, _, _ = exerciseRepo, workoutUC, router

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Server running on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
