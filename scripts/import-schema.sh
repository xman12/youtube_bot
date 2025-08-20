#!/bin/bash

# Скрипт для импорта схемы в базу данных

echo "🗄️  Импорт схемы в базу данных..."

# Проверяем, что мы в правильной директории
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Запустите скрипт из корневой директории проекта"
    exit 1
fi

# Выбираем файл схемы
if [ $# -eq 1 ]; then
    SCHEMA_FILE="$1"
else
    # Ищем файлы схем
    SCHEMA_FILES=$(find db_export -name "*schema*.sql" -type f | head -5)
    
    if [ -z "$SCHEMA_FILES" ]; then
        echo "❌ Файлы схем не найдены в папке db_export"
        echo "💡 Сначала создайте схему: ./scripts/create-clean-schema.sh"
        exit 1
    fi
    
    echo "📁 Найденные файлы схем:"
    echo "$SCHEMA_FILES"
    echo ""
    
    # Берем самый новый файл
    SCHEMA_FILE=$(ls -t db_export/*schema*.sql | head -1)
fi

if [ ! -f "$SCHEMA_FILE" ]; then
    echo "❌ Файл $SCHEMA_FILE не найден"
    exit 1
fi

echo "📤 Импорт схемы из файла: $SCHEMA_FILE"

# Проверяем, что контейнеры запущены
if ! docker ps | grep -q "youtube-bot-mysql"; then
    echo "⚠️  MySQL контейнер не запущен. Запускаем контейнеры..."
    docker-compose up -d mysql
    echo "⏳ Ждем запуска MySQL..."
    sleep 10
fi

# Проверяем подключение к MySQL
echo "🔍 Проверка подключения к MySQL..."
if ! docker exec youtube-bot-mysql mysqladmin ping -h localhost --silent; then
    echo "❌ Не удается подключиться к MySQL"
    exit 1
fi

echo "✅ MySQL доступен"

# Импортируем схему
echo "📥 Импорт схемы в базу данных..."
docker exec -i youtube-bot-mysql mysql \
    -u root \
    -p"${DB_PASSWORD:-admin}" \
    youtube_bot < "$SCHEMA_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Схема успешно импортирована"
    
    # Проверяем количество таблиц
    TABLE_COUNT=$(docker exec youtube-bot-mysql mysql \
        -u root \
        -p"${DB_PASSWORD:-admin}" \
        youtube_bot \
        -e "SHOW TABLES;" | wc -l)
    
    echo "📊 Таблиц в базе данных: $((TABLE_COUNT - 1))"
    
    # Показываем список таблиц
    echo "📋 Список таблиц:"
    docker exec youtube-bot-mysql mysql \
        -u root \
        -p"${DB_PASSWORD:-admin}" \
        youtube_bot \
        -e "SHOW TABLES;" | tail -n +2
    
else
    echo "❌ Ошибка при импорте схемы"
    exit 1
fi

echo ""
echo "🎉 Импорт схемы завершен!"
echo "💡 База данных готова к использованию"

