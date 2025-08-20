#!/bin/bash

# Загружаем переменные из .env
if [ -f .env ]; then
    source .env
else
    echo "❌ Файл .env не найден!"
    exit 1
fi

# Проверяем аргументы
if [ $# -eq 0 ]; then
    echo "📋 Использование: $0 <путь_к_файлу_бэкапа>"
    echo ""
    echo "📁 Доступные бэкапы:"
    if [ -d "./backups" ]; then
        ls -la ./backups/*.sql 2>/dev/null | head -10
        if [ -L "./backups/latest_backup.sql" ]; then
            echo ""
            echo "🔗 Последний бэкап: ./backups/latest_backup.sql"
        fi
    else
        echo "   Папка backups не найдена"
    fi
    exit 1
fi

BACKUP_FILE="$1"

# Проверяем существование файла
if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Файл бэкапа не найден: $BACKUP_FILE"
    exit 1
fi

echo "🗄️  Восстановление MySQL базы данных из резервной копии..."
echo "📁 Файл: $BACKUP_FILE"

# Проверяем, что контейнер запущен
if ! docker-compose ps mysql | grep -q "Up"; then
    echo "❌ MySQL контейнер не запущен"
    echo "   Запустите: docker-compose up -d mysql"
    exit 1
fi

echo "✅ MySQL контейнер запущен"

# Подтверждение пользователя
echo ""
echo "⚠️  ВНИМАНИЕ: Это действие перезапишет текущую базу данных!"
read -p "Продолжить? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Восстановление отменено"
    exit 1
fi

echo "📥 Восстановление базы данных..."
if docker-compose exec -T mysql mysql -u root -p"$DB_PASSWORD" "$DB_NAME" < "$BACKUP_FILE" 2>/dev/null; then
    echo "✅ База данных успешно восстановлена!"
    
    # Исправляем аутентификацию после восстановления
    echo "🔧 Настройка аутентификации..."
    docker-compose exec mysql mysql -u root -p"$DB_PASSWORD" -e "ALTER USER '$DB_LOGIN'@'%' IDENTIFIED WITH mysql_native_password BY '$DB_PASSWORD'; FLUSH PRIVILEGES;" 2>/dev/null || true
    echo "✅ Аутентификация настроена"
    
    echo "🎉 Восстановление завершено успешно!"
else
    echo "❌ Ошибка при восстановлении базы данных"
    exit 1
fi
