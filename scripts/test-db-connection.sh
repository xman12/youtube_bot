#!/bin/bash

echo "🔍 Проверка подключения к внешней базе данных..."

# Загружаем переменные окружения
if [ -f .env ]; then
    source .env
else
    echo "❌ Файл .env не найден!"
    exit 1
fi

# Проверяем наличие необходимых переменных
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_LOGIN" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    echo "❌ Не все переменные базы данных установлены в .env файле!"
    echo "   Проверьте: DB_HOST, DB_PORT, DB_LOGIN, DB_PASSWORD, DB_NAME"
    exit 1
fi

echo "📊 Параметры подключения:"
echo "   Хост: $DB_HOST"
echo "   Порт: $DB_PORT"
echo "   Пользователь: $DB_LOGIN"
echo "   База данных: $DB_NAME"

# Проверяем доступность хоста
echo "🌐 Проверка доступности хоста..."
if ping -c 1 "$DB_HOST" > /dev/null 2>&1; then
    echo "✅ Хост $DB_HOST доступен"
else
    echo "❌ Хост $DB_HOST недоступен"
    echo "   Проверьте настройки сети и firewall"
    exit 1
fi

# Проверяем доступность порта
echo "🔌 Проверка доступности порта..."
if nc -z "$DB_HOST" "$DB_PORT" 2>/dev/null; then
    echo "✅ Порт $DB_PORT на $DB_HOST доступен"
else
    echo "❌ Порт $DB_PORT на $DB_HOST недоступен"
    echo "   Проверьте настройки MySQL и firewall"
    exit 1
fi

# Проверяем подключение к базе данных
echo "🗄️  Проверка подключения к базе данных..."
if mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_LOGIN" -p"$DB_PASSWORD" "$DB_NAME" -e "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ Подключение к базе данных успешно"
    
    # Проверяем наличие таблиц
    echo "📋 Проверка структуры базы данных..."
    TABLES=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_LOGIN" -p"$DB_PASSWORD" "$DB_NAME" -e "SHOW TABLES;" 2>/dev/null | wc -l)
    if [ "$TABLES" -gt 1 ]; then
        echo "✅ База данных содержит таблицы"
    else
        echo "⚠️  База данных пуста или таблицы не найдены"
    fi
else
    echo "❌ Не удалось подключиться к базе данных"
    echo "   Проверьте:"
    echo "   - Правильность логина и пароля"
    echo "   - Права доступа пользователя"
    echo "   - Настройки MySQL (bind-address, user permissions)"
    exit 1
fi

echo "🎉 Все проверки пройдены успешно!"
echo "   База данных готова к использованию с YouTube Bot"
