package services

import (
	"context"
	"crypto/subtle"
	"fmt"
	"math/rand"
	"soft-hsm/internal/auth/dto"
	"soft-hsm/internal/auth/repository"
	"soft-hsm/internal/common/validators"
	"soft-hsm/internal/mailer"
	"soft-hsm/internal/user/models"
	userRepository "soft-hsm/internal/user/repository"
	"strings"
	"time"
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
	fmt.Println("Hashed Pass", hashedPassword)

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

	// is_verified | is_active_master | is_active

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

func generateOtp() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
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
		return nil, fmt.Errorf("incorrect login or password: %w", err)
	}

	// cpG3Fa7kiyU8aPgg7GtR/Q==$xgMmWlA7EVmgNqh8OxVmGwyya32MmgHPl3WDQBN23A4=

	// if !user.IsActive {
	// 	return nil, fmt.Errorf("accout not active")
	// }

	otp := generateOtp()

	fmt.Println(otp)

	token, err := s.cliamsService.GenerateOTP(int64(user.Id), otp, loginDto.Email)

	fmt.Println(token)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.tokenRepo.SaveToken(ctx, loginDto.Email, token, 86400); err != nil {
		return nil, err
	}

	go func() {
		if err := s.mailer.SendOTPEmail(user.Email, otp); err != nil {
			fmt.Printf("Failed to send activation email to %s: %v\n", user.Email, err)
		}
	}()

	return &dto.LoginResponseDTO{
		SessionToken: token,
		User:         user,
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

func (s *AuthService) ResetPassword(ctx context.Context, userID int64, data dto.ResetPasswordDTO) (bool, error) {
	fmt.Println("DTO:", data)
	if data.NewPassword != data.ConfirmNewPassword {
		return false, nil
	}

	user, err := s.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return false, err
	}

	fmt.Println("PASS:", data.CurrentPassword, user.Password)
	fmt.Println("CHECK:", s.passwordService.CheckPassword(data.CurrentPassword, user.Password))

	if !s.passwordService.CheckPassword(data.CurrentPassword, user.Password) {
		return false, err
	}

	newHashedPassword, err := s.passwordService.HashPassword(data.NewPassword)
	if err != nil {
		return false, err
	}

	fmt.Println("User: ", user)
	check, err := s.userRepo.SetPassword(ctx, userID, newHashedPassword)

	if err != nil {
		return false, err
	}

	return check, nil
}

func (s *AuthService) CheckMasterPassword(ctx context.Context, sessionToken string, otp string) (*dto.CheckMasterPasswordResponseDTO, error) {
	// user, err := s.userRepo.GetUserById(ctx, id)

	// if err != nil {
	// 	return nil, fmt.Errorf("cannot generate access token: %w", err)
	// }

	// if !user.IsActive {
	// 	return nil, fmt.Errorf("User not active: %w", err)
	// }

	payload, _ := s.cliamsService.ValidateOTPToken(sessionToken)

	fmt.Println("From DTO", otp)
	fmt.Println("OTP", payload.Otp)
	fmt.Println("ID", payload.Id)
	fmt.Println("Email", payload.Email)
	fmt.Println("Valid", subtle.ConstantTimeCompare([]byte(strings.TrimSpace(payload.Otp)), []byte(strings.TrimSpace(otp))) != 1)

	if subtle.ConstantTimeCompare([]byte(strings.TrimSpace(payload.Otp)), []byte(strings.TrimSpace(otp))) != 1 {
		return nil, fmt.Errorf("cannot validate session token")
	}

	accessToken, err := s.cliamsService.GenerateToken(int(payload.Id), payload.Email)
	if err != nil {
		return nil, fmt.Errorf("cannot generate access token: %w", err)
	}

	user, err := s.userRepo.GetUserByEmail(ctx, payload.Email)

	if err != nil {
		return nil, fmt.Errorf("incorrect login or password: %w", err)
	}

	return &dto.CheckMasterPasswordResponseDTO{
		AccessToken: accessToken,
		User:        user,
	}, nil
}
