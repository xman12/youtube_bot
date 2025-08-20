#!/bin/bash

# Скрипт для проверки структуры проекта и Docker конфигурации

echo "🔍 Проверка структуры проекта..."

# Проверяем основные файлы
echo "📁 Проверка основных файлов:"
[ -f "Dockerfile" ] && echo "✅ Dockerfile найден" || echo "❌ Dockerfile не найден"
[ -f "docker-compose.yml" ] && echo "✅ docker-compose.yml найден" || echo "❌ docker-compose.yml не найден"
[ -f "go.mod" ] && echo "✅ go.mod найден" || echo "❌ go.mod не найден"
[ -f "cmd/main.go" ] && echo "✅ cmd/main.go найден" || echo "❌ cmd/main.go не найден"

# Проверяем entrypoint скрипты
echo ""
echo "📁 Проверка entrypoint скриптов:"
[ -f "docker/entrypoint.sh" ] && echo "✅ docker/entrypoint.sh найден" || echo "❌ docker/entrypoint.sh не найден"
[ -f "docker/nginx-entrypoint.sh" ] && echo "✅ docker/nginx-entrypoint.sh найден" || echo "❌ docker/nginx-entrypoint.sh не найден"

# Проверяем права на выполнение
echo ""
echo "🔐 Проверка прав на выполнение:"
[ -x "docker/entrypoint.sh" ] && echo "✅ docker/entrypoint.sh исполняемый" || echo "❌ docker/entrypoint.sh не исполняемый"
[ -x "docker/nginx-entrypoint.sh" ] && echo "✅ docker/nginx-entrypoint.sh исполняемый" || echo "❌ docker/nginx-entrypoint.sh не исполняемый"

# Проверяем структуру Go проекта
echo ""
echo "🐹 Проверка структуры Go проекта:"
[ -d "cmd" ] && echo "✅ папка cmd существует" || echo "❌ папка cmd не существует"
[ -d "internal" ] && echo "✅ папка internal существует" || echo "❌ папка internal не существует"
[ -d "internal/app" ] && echo "✅ папка internal/app существует" || echo "❌ папка internal/app не существует"

# Проверяем nginx конфигурацию
echo ""
echo "🌐 Проверка nginx конфигурации:"
[ -f "nginx/nginx.conf" ] && echo "✅ nginx/nginx.conf найден" || echo "❌ nginx/nginx.conf не найден"

# Проверяем переменные окружения
echo ""
echo "🔧 Проверка переменных окружения:"
[ -f ".env" ] && echo "✅ .env файл найден" || echo "⚠️  .env файл не найден (создайте из env.example)"
[ -f "env.example" ] && echo "✅ env.example найден" || echo "❌ env.example не найден"

echo ""
echo "🎉 Проверка завершена!"

