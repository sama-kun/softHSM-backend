package services

import (
	"context"
	"fmt"
	"soft-hsm/internal/user/models"
	"soft-hsm/internal/user/repository"
)

type UserService struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserService(userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Me(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.userRepo.GetUserById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("cannot found user: %w", err)
	}

	return user, nil
}