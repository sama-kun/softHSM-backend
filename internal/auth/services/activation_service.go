package services

import (
	"context"
	"fmt"
	"soft-hsm/internal/auth/dto"
	"soft-hsm/internal/auth/repository"
	userRepository "soft-hsm/internal/user/repository"
)

type ActivationService struct {
	userRepo userRepository.UserRepositoryInterface
	tokenRepo repository.TokenRepositoryInterface
	claimsService ClaimsService
}

func NewActivationService(userRepo userRepository.UserRepositoryInterface, tokenRepo repository.TokenRepositoryInterface, claimsService *ClaimsService) *ActivationService {
	return &ActivationService{userRepo: userRepo, tokenRepo: tokenRepo, claimsService: *claimsService}
}

func (s *ActivationService) ActiveUser(ctx context.Context, token string) (*dto.RegisterResponseDTO,error) {
	claims, err := s.claimsService.ValidateActivationToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid activation token: %w", err)
	}
	tokenFromRedis, err := s.tokenRepo.GetActivationToken(ctx, claims.Email)

	if err != nil || token != tokenFromRedis {
		return nil, fmt.Errorf("invalid activation token: %w", err)
	}



	if err := s.userRepo.ActiveUser(ctx, claims.Email); err != nil {
		return nil, fmt.Errorf("user can not update: %w", err)
	}

	return &dto.RegisterResponseDTO{
		Email: claims.Email,
		Success: "User successfully activated",
	}, nil
}