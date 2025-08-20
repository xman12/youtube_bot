#!/bin/bash

echo "🔍 Быстрая проверка доступности MySQL..."

# Проверяем, что порт 3306 открыт
if command -v nc >/dev/null 2>&1; then
    echo "📡 Проверка порта 3306..."
    if nc -z localhost 3306 2>/dev/null; then
        echo "✅ Порт 3306 открыт и доступен"
    else
        echo "❌ Порт 3306 недоступен"
        exit 1
    fi
elif command -v telnet >/dev/null 2>&1; then
    echo "📡 Проверка порта 3306 через telnet..."
    if timeout 5 telnet localhost 3306 2>/dev/null | grep -q "Connected"; then
        echo "✅ Порт 3306 открыт и доступен"
    else
        echo "❌ Порт 3306 недоступен"
        exit 1
    fi
else
    echo "⚠️  Не удалось проверить порт (установите netcat или telnet)"
fi

echo ""
echo "🎯 Параметры подключения для клиентов:"
echo "   Host: localhost"
echo "   Port: 3306"
echo "   Database: youtube_bot"
echo "   User: admin"
echo "   Password: admin"
echo "   Authentication: mysql_native_password"
echo ""
echo "🔗 Строка подключения:"
echo "   mysql://admin:admin@localhost:3306/youtube_bot"
echo ""
echo "✅ MySQL готов к подключению внешних клиентов!"
