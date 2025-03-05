package services

import (
	"context"
	"errors"
	"fmt"
	"soft-hsm/internal/auth/dto"
	"soft-hsm/internal/auth/repository"
	"soft-hsm/internal/common/validators"
	"soft-hsm/internal/mailer"
	"soft-hsm/internal/user/models"
	userRepository "soft-hsm/internal/user/repository"
)

type AuthService struct {
	tokenRepo       repository.TokenRepository
	cliamsService   ClaimsService
	userRepo        userRepository.UserRepository
	mailer          mailer.Mailer
	passwordService PasswordService
}

func NewAuthService(tokenRepo repository.TokenRepository, claimsService ClaimsService, userRepo userRepository.UserRepository, mailer mailer.Mailer, passwordService PasswordService) *AuthService {
	return &AuthService{tokenRepo: tokenRepo, cliamsService: claimsService, userRepo: userRepo, mailer: mailer, passwordService: passwordService}
}

func (s *AuthService) Register(ctx context.Context, registerDTO dto.RegisterDTO) (*dto.RegisterResponseDTO, error) {
	if err := validators.ValidateStruct(registerDTO); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := s.userRepo.IsEmailTaken(ctx, registerDTO.Email); err == nil {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := s.passwordService.HashPassword(registerDTO.Password)

	if err != nil {
		return nil, fmt.Errorf("password cannot hash: %w", err)
	}

	newUser, err := s.userRepo.SaveUser(ctx, &models.User{Email: registerDTO.Email, Password: hashedPassword, Login: registerDTO.Login})
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	activationToken, err := s.cliamsService.GenerateActivationToken(newUser.Email)
	if err != nil {
		return nil, fmt.Errorf("cannot generate a token: %w", err)
	}

	if err := s.tokenRepo.SaveActivationToken(ctx, newUser.Email, activationToken, int64(s.cliamsService.cfg.JWTConfig.ActivationExpires)); err != nil {
		return nil, fmt.Errorf("cannot save activation token: %w", err)
	}

	err = s.mailer.SendActivationEmail(newUser.Email, activationToken)

	if err != nil {
		return nil, fmt.Errorf("cannot send to email: %w", err)
	}

	return &dto.RegisterResponseDTO{
		Success: "User successfully created and activation email sent",
		Email:   newUser.Email,
	}, nil

}

func (s *AuthService) Login(ctx context.Context, loginDto dto.LoginDTO) (*dto.LoginResponseDTO, error) {

	if err := validators.ValidateStruct(loginDto); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	user, err := s.userRepo.GetUserByEmail(ctx, loginDto.Email)

	if err != nil {
		if errors.Is(err, userRepository.ErrUserNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	token, err := s.cliamsService.GenerateToken(int(user.Id), loginDto.Email)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.tokenRepo.SaveToken(ctx, loginDto.Email, token, 86400); err != nil {
		return nil, err
	}

	return &dto.LoginResponseDTO{
		AccessToken: token,
		User:        user,
	}, nil
}

// Логаут: Удалить токен
func (s *AuthService) Logout(ctx context.Context, email string) error {
	return s.tokenRepo.DeleteToken(ctx, email)
}
