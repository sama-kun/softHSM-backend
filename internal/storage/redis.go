package storage

import (
	"context"
	"fmt"
	"time"

	"soft-hsm/internal/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // Пароль, если есть
		DB:       cfg.DB,       // Номер базы (по умолчанию 0)
	})

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к Redis: %w", err)
	}

	return &Redis{client: client}, nil
}

// Close закрывает соединение с Redis
func (r *Redis) Close() error {
	return r.client.Close()
}
func (r *Redis) Save(ctx context.Context, key, value string, expiresInMinute int64) error {
	return r.client.Set(ctx, key, value, time.Duration(expiresInMinute)*time.Minute).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete удаляет ключ из Redis
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
