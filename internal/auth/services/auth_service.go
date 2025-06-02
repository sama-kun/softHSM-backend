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
		mailer:          mailer,
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

	fmt.Println("Token:", activationToken)

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

	if err != nil {
		return nil, fmt.Errorf("incorrect login or password: %w", err)
	}

	if !s.passwordService.CheckPassword(loginDto.Password, user.Password) {
		return nil, fmt.Errorf("incorrect login or password")
	}

	// if !user.IsActive {
	// 	return nil, fmt.Errorf("accout not active")
	// }

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

func (s *AuthService) SetMasterPassword(ctx context.Context, id int64, masterPassword string) (*dto.SetMasterPasswordResponseDTO, error) {
	hashedPassword, err := s.passwordService.HashPassword(masterPassword)
	if err != nil {
		return nil, fmt.Errorf("master password cannot hash: %w", err)
	}

	err = s.userRepo.SetMasterPassword(ctx, id, hashedPassword)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &dto.SetMasterPasswordResponseDTO{
		Success: "Master password was succesfully setted",
		Id:      id,
	}, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, userID int64, data dto.ResetPasswordDTO) bool {
	if data.NewPassword != data.ConfirmNewPassword {
		return false
	}

	user, err := s.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return false
	}

	if s.passwordService.CheckPassword(data.CurrentPassword, user.Password) {
		return false
	}

	newHashedPassword, err := s.passwordService.HashPassword(data.NewPassword)
	if err != nil {
		return false
	}
	check, _ := s.userRepo.SetPassword(ctx, userID, newHashedPassword)

	return check
}

func (s *AuthService) CheckMasterPassword(ctx context.Context, id int64, masterPassword string) (*dto.CheckMasterPasswordResponseDTO, error) {
	user, err := s.userRepo.GetUserById(ctx, id)

	if !s.passwordService.CheckPassword(masterPassword, user.MasterPassword) || err != nil {
		return nil, fmt.Errorf("incorrect login or password: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("User not active: %w", err)
	}

	sessionToken, err := s.cliamsService.GenerateBlockchainOTP(user.Id)

	if err != nil {
		return nil, fmt.Errorf("cannot generate session token: %w", err)
	}

	fmt.Println(sessionToken)

	return &dto.CheckMasterPasswordResponseDTO{
		SessionToken: sessionToken,
		Id:           id,
	}, nil
}
