#!/bin/bash

set -e  # Остановить скрипт при ошибке

# Настройки базы данных
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="12345678"
DB_NAME="soft_hsm"

# 1. Применяем миграции
echo "🔄 Выполнение миграций..."
for file in ./migrations/*.up.sql; do
    echo "▶️ Применение миграции: $file"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f "$file"
done

echo "✅ Миграции успешно применены."

# 2. Компиляция и запуск приложения
echo "🚀 Запуск приложения..."
go run cmd/main.go