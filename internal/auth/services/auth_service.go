package services

import (
	"context"
	"fmt"
	"soft-hsm/internal/auth/dto"
	"soft-hsm/internal/auth/repository"
	"soft-hsm/internal/common/validators"
	"soft-hsm/internal/mailer"
	"soft-hsm/internal/user/models"
	userRepository "soft-hsm/internal/user/repository"
)

type AuthService struct {
	tokenRepo       repository.TokenRepositoryInterface
	cliamsService   *ClaimsService
	userRepo        userRepository.UserRepositoryInterface
	mailer          *mailer.Mailer
	passwordService *PasswordService
}

func NewAuthService(
	tokenRepo repository.TokenRepositoryInterface,
	claimsService *ClaimsService,
	userRepo userRepository.UserRepositoryInterface,
	mailer *mailer.Mailer,
	passwordService *PasswordService,
) *AuthService {
	return &AuthService{
		tokenRepo:       tokenRepo,
		cliamsService:   claimsService,
		userRepo:        userRepo,
		mailer:         mailer,
		passwordService: passwordService,
	}
}
func (s *AuthService) Register(ctx context.Context, registerDTO dto.RegisterDTO) (*dto.RegisterResponseDTO, error) {
	if err := validators.ValidateStruct(registerDTO); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := s.userRepo.IsEmailTaken(ctx, registerDTO.Email); err != nil {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := s.passwordService.HashPassword(registerDTO.Password)
	if err != nil {
		return nil, fmt.Errorf("password cannot hash: %w", err)
	}

	newUser, err := s.userRepo.SaveUser(ctx, &models.User{
		Email:    registerDTO.Email,
		Password: hashedPassword,
		Login:    registerDTO.Login,
	})
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

	// Send email asynchronously
	go func() {
		if err := s.mailer.SendActivationEmail(newUser.Email, activationToken); err != nil {
			fmt.Printf("Failed to send activation email to %s: %v\n", newUser.Email, err)
		}
	}()

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

	if !s.passwordService.CheckPassword(loginDto.Password, user.Password) || err != nil {
		return nil, fmt.Errorf("incorrect login or password: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("accout not active")
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


