package storage

import (
	"context"
	"errors"
	"fmt"

	"soft-hsm/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

var (
	ErrorURLNotFound = errors.New("URL not found")
	ErrorURLExists   = errors.New("URL exists")
)

type Postgres struct {
	conn *pgx.Conn
}

func NewPostgresDB(cfg config.DBConfig) (*Postgres, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при инициализации миграции: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("ошибка при применении миграций: %w", err)
	}

	fmt.Println("✅ Успешное подключение к PostgreSQL")
	return &Postgres{conn: conn}, nil
}

func (p *Postgres) Close() {
	p.conn.Close(context.Background())
	fmt.Println("🔌 Соединение с PostgreSQL закрыто")
}
