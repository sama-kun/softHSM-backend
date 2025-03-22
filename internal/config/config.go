package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env          string `yaml:"env" env-default:"local"` //env-required: "true" in prod
	Database     DBConfig
	HTTPServer   `yaml:"http_server"`
	RedisConfig  `yaml:"redis"`
	JWTConfig    `yaml:"jwt_config"`
	MailerConfig `yaml:"mailer"`
}

type MailerConfig struct {
	From     string `yaml:"from" env-default:"samgar.robot@gmail.com"`
	Password string `yaml:"password" validate:"required,email"`
	SMTPHost string `yaml:"smtp_host" validate:"required"`
	SMTPPort int    `yaml:"smtp_port" validate:"required,numeric"`
}

type JWTConfig struct {
	Secret            string `yaml:"secret" env-default:"ono"`
	Expires           int    `yaml:"expires" env-default:"1440"`
	ActivationSecret  string `yaml:"activation_secret"`
	ActivationExpires int    `yaml:"activation_expires"`
	SessionSecret     string `yaml:"session_secret"`
	SessionExpires    string `yaml:"session_expires"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode" env-default:"false"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"4s"`
}

func MustLoad() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️ Не удалось загрузить .env файл, используем системные переменные")
	}
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIGPATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
