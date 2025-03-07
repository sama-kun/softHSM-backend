package repository

import (
	"context"

	"soft-hsm/internal/storage"
)

type TokenRepositoryInterface interface {
	SaveToken(ctx context.Context, email, token string, expiresInMinute int64) error
	GetToken(ctx context.Context, email string) (string, error)
	DeleteToken(ctx context.Context, email string) error
	SaveActivationToken(ctx context.Context, email, token string, expiresIn int64) error
	GetActivationToken(ctx context.Context, email string) (string, error)
	DeleteActivationToken(ctx context.Context, email string) error
}

type TokenRepository struct {
	redis *storage.Redis
}

func NewTokenRepository(redis *storage.Redis) TokenRepositoryInterface {
	return &TokenRepository{redis: redis}
}

func (r *TokenRepository) SaveToken(ctx context.Context, email, token string, expiresInMinute int64) error {
	return r.redis.Save(ctx, "token:"+email, token, expiresInMinute)
}

func (r *TokenRepository) GetToken(ctx context.Context, email string) (string, error) {
	return r.redis.Get(ctx, "token:"+email)
}

func (r *TokenRepository) DeleteToken(ctx context.Context, email string) error {
	return r.redis.Delete(ctx, "token:"+email)
}

func (r *TokenRepository) SaveActivationToken(ctx context.Context, email, token string, expiresIn int64) error {
	return r.redis.Save(ctx, "activationToken:"+email, token, expiresIn)
}

func (r *TokenRepository) GetActivationToken(ctx context.Context, email string) (string, error) {
	return r.redis.Get(ctx, "activationToken:"+email)
}

func (r *TokenRepository) DeleteActivationToken(ctx context.Context, email string) error {
	return r.redis.Delete(ctx, "activationToken:"+email)
}
