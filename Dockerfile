# 1. Используем минимальный образ с Go
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 2. Устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# 3. Копируем исходный код и компилируем бинарник
COPY . .
RUN go build -o main ./cmd/main.go

# 4. Создаем финальный образ
FROM alpine:latest

WORKDIR /app

# 5. Копируем бинарник из builder'а
COPY --from=builder /app/main .
COPY ./migrations ./migrations

# 6. Настраиваем переменные окружения
# ENV POSTGRES_HOST="postgres"
# ENV POSTGRES_PORT="5432"
# ENV POSTGRES_USER="postgres"
# ENV POSTGRES_PASSWORD="12345678"
# ENV POSTGRES_DB="soft_hsm"
# ENV REDIS_ADDR="redis:6379"

# 7. Запуск приложения
CMD ["/app/main"]