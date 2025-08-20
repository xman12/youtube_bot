#!/bin/bash

# Загружаем переменные из .env
if [ -f .env ]; then
    source .env
else
    echo "❌ Файл .env не найден!"
    exit 1
fi

# Создаем папку для бэкапов если её нет
BACKUP_DIR="./backups"
mkdir -p "$BACKUP_DIR"

# Генерируем имя файла с датой и временем
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
BACKUP_FILE="$BACKUP_DIR/youtube_bot_$TIMESTAMP.sql"

echo "🗄️  Создание резервной копии MySQL базы данных..."

# Проверяем, что контейнер запущен
if ! docker-compose ps mysql | grep -q "Up"; then
    echo "❌ MySQL контейнер не запущен"
    echo "   Запустите: docker-compose up -d mysql"
    exit 1
fi

echo "✅ MySQL контейнер запущен"

# Создаем резервную копию
echo "📥 Создание дампа базы данных..."
if docker-compose exec mysql mysqldump -u root -p"$DB_PASSWORD" "$DB_NAME" > "$BACKUP_FILE"; then
    echo "✅ Резервная копия создана: $BACKUP_FILE"
    
    # Показываем размер файла
    FILE_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo "📊 Размер файла: $FILE_SIZE"
    
    # Создаем символическую ссылку на последний бэкап
    LATEST_BACKUP="$BACKUP_DIR/latest_backup.sql"
    ln -sf "$BACKUP_FILE" "$LATEST_BACKUP"
    echo "🔗 Создана ссылка: $LATEST_BACKUP -> $BACKUP_FILE"
    
    echo "🎉 Резервное копирование завершено успешно!"
else
    echo "❌ Ошибка при создании резервной копии"
    exit 1
fi
