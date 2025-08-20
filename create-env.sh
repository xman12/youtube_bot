#!/bin/bash

if [ ! -f .env ]; then
    echo "Создание .env файла..."
    cat > .env << EOF
# Telegram Bot Token
TELEGRAM_API_TOKEN=your_telegram_bot_token_here

# Database Configuration (Local MySQL Container)
DB_LOGIN=youtube_user
DB_PASSWORD=your_secure_password
DB_HOST=mysql
DB_PORT=3306
DB_NAME=youtube_bot

# Proxy Configuration (optional)
PATH_TO_PROXY=./proxy_files/proxy.txt

# Video and Audio paths
PATH_TO_LOAD_VIDEO=/tmp/videos
PATH_TO_LOAD_AUDIO=/tmp/audio
PATH_TO_LOAD_IMG=/tmp/images

# Настройки очистки временных файлов
CLEANUP_INTERVAL_SECONDS=3600
CLEANUP_AGE_HOURS=2

# Telegram API Server Configuration
TELEGRAM_API_ID=111111111
TELEGRAM_API_HASH=telegram_api_hash
TELEGRAM_API_URL=http://telegram-api:8081

EOF
    echo ".env файл создан. Пожалуйста, отредактируйте его с вашими реальными значениями."
else
    echo ".env файл уже существует."
fi
