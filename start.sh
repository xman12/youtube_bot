#!/bin/bash

echo "🚀 Запуск YouTube Bot с Docker и мониторингом..."

# Проверяем наличие .env файла
if [ ! -f .env ]; then
    echo "📝 Создание .env файла..."
    chmod +x create-env.sh
    ./create-env.sh
    echo "⚠️  Пожалуйста, отредактируйте .env файл с вашими реальными значениями!"
    read -p "Нажмите Enter после редактирования .env файла..."
fi

# Проверяем наличие дампа базы данных
echo "🗄️  Проверка дампа базы данных..."
if [ -f "db_export/youtube_bot_clean_schema.sql" ]; then
    echo "✅ Дамп базы данных найден"
else
    echo "⚠️  Дамп базы данных не найден в db_export/"
    echo "   Система создаст пустую базу данных"
fi

# Останавливаем существующие контейнеры
echo "🛑 Остановка существующих контейнеров..."
docker-compose down

# Собираем и запускаем контейнеры
echo "🔨 Сборка и запуск контейнеров..."
docker-compose up -d --build

# Ждем запуска сервисов
echo "⏳ Ожидание запуска сервисов..."
sleep 10

# Проверяем статус
echo "📊 Статус сервисов:"
docker-compose ps

echo ""
echo "✅ Система запущена!"
echo ""
echo "🌐 Доступ к сервисам:"
echo "   YouTube Bot: https://localhost:443"
echo "   Prometheus:  https://localhost:9090 (HTTP: http://localhost:9091)"
echo "   Grafana:     https://localhost:3000 (HTTP: http://localhost:3001) (admin/admin)"
echo ""
echo "📝 Для просмотра логов используйте:"
echo "   docker-compose logs -f [service_name]"
echo ""
echo "🛑 Для остановки используйте:"
echo "   docker-compose down"
