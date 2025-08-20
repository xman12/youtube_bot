#!/bin/bash

echo "🔍 Проверка импорта базы данных..."

# Загружаем переменные окружения
if [ -f .env ]; then
    source .env
else
    echo "❌ Файл .env не найден!"
    exit 1
fi

# Проверяем, запущен ли MySQL контейнер
if ! docker-compose ps mysql | grep -q "Up"; then
    echo "❌ MySQL контейнер не запущен"
    echo "   Запустите: docker-compose up -d mysql"
    exit 1
fi

echo "✅ MySQL контейнер запущен"

# Ждем, пока MySQL будет готов
echo "⏳ Ожидание готовности MySQL..."
sleep 10

# Проверяем подключение к базе данных
echo "🗄️  Проверка подключения к базе данных..."
if docker-compose exec mysql mysql -u "$DB_LOGIN" -p"$DB_PASSWORD" "$DB_NAME" -e "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ Подключение к базе данных успешно"
else
    echo "❌ Не удалось подключиться к базе данных"
    exit 1
fi

# Проверяем количество таблиц
echo "📋 Проверка структуры базы данных..."
TABLE_COUNT=$(docker-compose exec mysql mysql -u "$DB_LOGIN" -p"$DB_PASSWORD" "$DB_NAME" -e "SHOW TABLES;" 2>/dev/null | wc -l)

if [ "$TABLE_COUNT" -gt 1 ]; then
    echo "✅ База данных содержит $((TABLE_COUNT - 1)) таблиц"
    
    # Показываем список таблиц
    echo "📊 Список таблиц:"
    docker-compose exec mysql mysql -u "$DB_LOGIN" -p"$DB_PASSWORD" "$DB_NAME" -e "SHOW TABLES;" 2>/dev/null | tail -n +2
    
    # Проверяем количество записей в основных таблицах
    echo "📈 Статистика записей:"
    
    # Проверяем таблицу users
    if docker-compose exec mysql mysql -u "$DB_LOGIN" -p"$DB_PASSWORD" "$DB_NAME" -e "DESCRIBE users;" > /dev/null 2>&1; then
        USER_COUNT=$(docker-compose exec mysql mysql -u "$DB_LOGIN" -p"$DB_PASSWORD" "$DB_NAME" -e "SELECT COUNT(*) FROM users;" 2>/dev/null | tail -n 1)
        echo "   👥 Пользователей: $USER_COUNT"
    fi
    
    # Проверяем таблицу videos
    if docker-compose exec mysql mysql -u "$DB_LOGIN" -p"$DB_PASSWORD" "$DB_NAME" -e "DESCRIBE videos;" > /dev/null 2>&1; then
        VIDEO_COUNT=$(docker-compose exec mysql mysql -u "$DB_LOGIN" -p"$DB_PASSWORD" "$DB_NAME" -e "SELECT COUNT(*) FROM videos;" 2>/dev/null | tail -n 1)
        echo "   🎥 Видео: $VIDEO_COUNT"
    fi
    
else
    echo "⚠️  База данных пуста или содержит только системные таблицы"
    echo "   Возможно, дамп не был импортирован"
fi

echo "🎉 Проверка завершена"
