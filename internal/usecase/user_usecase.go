package usecase

import "workout-tracker/internal/domain"

type UserUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) *UserUsecase {
	return &UserUsecase{repo: r}
}
