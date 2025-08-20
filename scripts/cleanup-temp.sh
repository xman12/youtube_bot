#!/bin/bash

# Загружаем переменные из .env
if [ -f .env ]; then
    source .env
else
    echo "❌ Файл .env не найден!"
    exit 1
fi

# Настройки по умолчанию
CLEANUP_AGE_HOURS="${CLEANUP_AGE_HOURS:-2}"
CLEANUP_INTERVAL_SECONDS="${CLEANUP_INTERVAL_SECONDS:-3600}"

echo "🧹 Очистка временных файлов в контейнере youtube-bot..."
echo "⏰ Удаляем файлы старше $CLEANUP_AGE_HOURS часов"

# Проверяем, что контейнер запущен
if ! docker-compose ps youtube-bot | grep -q "Up"; then
    echo "❌ YouTube Bot контейнер не запущен"
    echo "   Запустите: docker-compose up -d youtube-bot"
    exit 1
fi

echo "✅ YouTube Bot контейнер запущен"

# Выполняем очистку в контейнере
echo "🗑️  Выполняем очистку..."

# Очищаем файлы в /tmp/videos
echo "📹 Очистка старых видео файлов..."
FOUND_VIDEOS=$(docker-compose exec youtube-bot find /tmp/videos -type f -mmin +120 2>/dev/null | wc -l)
echo "   Найдено $FOUND_VIDEOS файлов для удаления"
docker-compose exec youtube-bot find /tmp/videos -type f -mmin +120 -delete 2>/dev/null || true

# Очищаем файлы в /tmp/audio
echo "🎵 Очистка старых аудио файлов..."
FOUND_AUDIO=$(docker-compose exec youtube-bot find /tmp/audio -type f -mmin +120 2>/dev/null | wc -l)
echo "   Найдено $FOUND_AUDIO файлов для удаления"
docker-compose exec youtube-bot find /tmp/audio -type f -mmin +120 -delete 2>/dev/null || true

# Очищаем файлы в /tmp/images
echo "🖼️  Очистка старых изображений..."
FOUND_IMAGES=$(docker-compose exec youtube-bot find /tmp/images -type f -mmin +120 2>/dev/null | wc -l)
echo "   Найдено $FOUND_IMAGES файлов для удаления"
docker-compose exec youtube-bot find /tmp/images -type f -mmin +120 -delete 2>/dev/null || true

echo "✅ Очистка завершена!"

# Показываем статистику по папкам
echo ""
echo "📊 Статистика по папкам:"
echo "   /tmp/videos: $(docker-compose exec youtube-bot find /tmp/videos -type f 2>/dev/null | wc -l) файлов"
echo "   /tmp/audio: $(docker-compose exec youtube-bot find /tmp/audio -type f 2>/dev/null | wc -l) файлов"
echo "   /tmp/images: $(docker-compose exec youtube-bot find /tmp/images -type f 2>/dev/null | wc -l) файлов"

echo ""
echo "💡 Автоматическая очистка выполняется каждые $((CLEANUP_INTERVAL_SECONDS / 3600)) часов"
echo "   Настройки можно изменить в .env файле:"
echo "   CLEANUP_INTERVAL_SECONDS=$CLEANUP_INTERVAL_SECONDS"
echo "   CLEANUP_AGE_HOURS=$CLEANUP_AGE_HOURS"
