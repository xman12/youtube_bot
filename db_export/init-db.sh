#!/bin/bash

echo "🗄️  Инициализация базы данных..."

# Ждем, пока MySQL полностью запустится
echo "⏳ Ожидание запуска MySQL..."
while ! mysqladmin ping -h"localhost" --silent; do
    sleep 1
done

echo "✅ MySQL запущен"

# Проверяем, есть ли дамп для импорта
if [ -f "/docker-entrypoint-initdb.d/youtube_bot_clean_schema.sql" ]; then
    echo "📥 Импорт дампа базы данных..."
    
    # Импортируем дамп
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE" < /docker-entrypoint-initdb.d/youtube_bot_clean_schema.sql
    
    if [ $? -eq 0 ]; then
        echo "✅ Дамп успешно импортирован"
        
        # Проверяем количество таблиц
        TABLE_COUNT=$(mysql -u root -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE" -e "SHOW TABLES;" | wc -l)
        echo "📊 Импортировано таблиц: $((TABLE_COUNT - 1))"
    else
        echo "❌ Ошибка при импорте дампа"
        exit 1
    fi
else
    echo "⚠️  Дамп не найден, создаем пустую базу данных"
fi

       echo "🎉 База данных готова к использованию"
       
       # Исправляем аутентификацию для внешних клиентов
       echo "🔧 Настройка аутентификации для внешних клиентов..."
       mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "ALTER USER '$MYSQL_USER'@'%' IDENTIFIED WITH mysql_native_password BY '$MYSQL_PASSWORD'; FLUSH PRIVILEGES;" 2>/dev/null || true
       echo "✅ Аутентификация настроена для внешних клиентов"
