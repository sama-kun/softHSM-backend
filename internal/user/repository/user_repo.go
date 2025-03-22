package repository

import (
	"context"
	"errors"
	"fmt"
	"soft-hsm/internal/storage"
	"soft-hsm/internal/user/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	SaveUser(ctx context.Context, user *models.User) (*models.User, error)
	IsEmailTaken(ctx context.Context, email string) error
	ActiveUser(ctx context.Context, email string) error
	GetUserById(ctx context.Context, id int64) (*models.User, error)
	SetMasterPassword(ctx context.Context, id int64, hashedMasterPassword string) error
}

type UserRepository struct {
	db *storage.Postgres
}

func NewUserRepository(db *storage.Postgres) UserRepositoryInterface {
	return &UserRepository{db: db}
}

func (r *UserRepository) ActiveUser(ctx context.Context, email string) error {
	query := `
		UPDATE users
		SET is_verified = TRUE
		WHERE email = $1 AND is_deleted = FALSE
	`

	cmdTag, err := r.db.Conn().Exec(ctx, query, email)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with email: %s", email)
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, master_password, login, is_active, created_at, updated_at FROM users 
		WHERE email = $1 AND is_deleted = FALSE
	`

	var user models.User

	err := r.db.Conn().QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.MasterPassword,
		&user.Login,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, email, login, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_deleted = FALSE
	`

	var user models.User

	err := r.db.Conn().QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Email,
		&user.Login,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user *models.User) (*models.User, error) {
	const fn = "repository.user"

	query := `
		INSERT INTO users (email, password, login) 
		VALUES ($1, $2, $3) 
		RETURNING id, is_deleted, updated_at, created_at`
	err := r.db.Conn().QueryRow(ctx, query, user.Email, user.Password, user.Login).
		Scan(&user.Id, &user.IsDeleted, &user.UpdatedAt, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (r *UserRepository) IsEmailTaken(ctx context.Context, email string) error {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_deleted = FALSE)"
	err := r.db.Conn().QueryRow(ctx, query, email).Scan(&exists)

	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if exists {
		return errors.New("user already exists")
	}

	return nil
}

func (r *UserRepository) SetMasterPassword(ctx context.Context, id int64, hashedMasterPassword string) error {
	const fn = "userRepo.SetMasterPassword"

	query := `
	UPDATE users 
    SET 
        master_password = $2, 
        is_active_master = TRUE
				is_active = TRUE
    WHERE 
        id = $1 
        AND is_deleted = FALSE
				AND is_verified = TRUE
	`

	cmdTag, err := r.db.Conn().Exec(ctx, query, id, hashedMasterPassword)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id /%s/: %d", fn, id)
	}

	return nil
}
