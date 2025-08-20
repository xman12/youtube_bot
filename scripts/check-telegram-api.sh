#!/bin/bash

echo "🔍 Проверка статуса Telegram Bot API Server..."

# Проверяем, что контейнер запущен
if ! docker-compose ps telegram-api | grep -q "Up"; then
    echo "❌ Telegram API контейнер не запущен"
    echo "   Запустите: docker-compose up -d telegram-api"
    exit 1
fi

echo "✅ Telegram API контейнер запущен"

# Проверяем доступность API
echo "🌐 Проверка доступности API..."
if curl -s http://localhost:8081/api/status > /dev/null 2>&1; then
    echo "✅ Telegram API доступен на http://localhost:8081"
else
    echo "❌ Telegram API недоступен на http://localhost:8081"
    echo "   Проверьте логи: docker-compose logs telegram-api"
    exit 1
fi

# Проверяем логи
echo "📋 Последние логи Telegram API:"
docker-compose logs telegram-api --tail=10

echo ""
echo "🎯 Параметры подключения:"
echo "   URL: http://localhost:8081"
echo "   Внутренний URL: http://telegram-api:8081"
echo "   Порт: 8081"

# Проверяем переменные окружения
if [ -f .env ]; then
    source .env
    echo ""
    echo "⚙️  Настройки из .env:"
    echo "   TELEGRAM_API_ID: ${TELEGRAM_API_ID:-не установлен}"
    echo "   TELEGRAM_API_HASH: ${TELEGRAM_API_HASH:-не установлен}"
    echo "   TELEGRAM_API_TOKEN: ${TELEGRAM_API_TOKEN:0:10}..."
fi

echo ""
echo "💡 Для получения API_ID и API_HASH:"
echo "   1. Перейдите на https://my.telegram.org"
echo "   2. Войдите в свой аккаунт"
echo "   3. Перейдите в 'API development tools'"
echo "   4. Создайте новое приложение"
echo "   5. Скопируйте API_ID и API_HASH в .env файл"
