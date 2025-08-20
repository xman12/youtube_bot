#!/usr/bin/env bash
set -euo pipefail

YTDLP_BIN="/usr/bin/yt-dlp"
INTERVAL="${YTDLP_UPDATE_INTERVAL_SECONDS:-21600}" # 6h by default
CLEANUP_INTERVAL="${CLEANUP_INTERVAL_SECONDS:-3600}" # 1h by default
CLEANUP_AGE="${CLEANUP_AGE_HOURS:-2}" # 2h by default

update_loop() {
  while true; do
    if command -v "$YTDLP_BIN" >/dev/null 2>&1; then
      echo "[yt-dlp] Checking for updates..."
      apk update && apk upgrade yt-dlp || true
    else
      echo "[yt-dlp] Binary not found at $YTDLP_BIN"
    fi
    sleep "$INTERVAL"
  done
}

cleanup_loop() {
  while true; do
    echo "[cleanup] Starting cleanup of old temporary files..."
    
    # Очищаем файлы старше 2 часов в /tmp/videos
    if [ -d "/tmp/videos" ]; then
      find /tmp/videos -type f -mmin +120 -delete 2>/dev/null || true
      echo "[cleanup] Cleaned old video files"
    fi
    
    # Очищаем файлы старше 2 часов в /tmp/audio
    if [ -d "/tmp/audio" ]; then
      find /tmp/audio -type f -mmin +120 -delete 2>/dev/null || true
      echo "[cleanup] Cleaned old audio files"
    fi
    
    # Очищаем файлы старше 2 часов в /tmp/images
    if [ -d "/tmp/images" ]; then
      find /tmp/images -type f -mmin +120 -delete 2>/dev/null || true
      echo "[cleanup] Cleaned old image files"
    fi
    
    # Очищаем другие временные файлы в /tmp (кроме системных)
    find /tmp -type f -mmin +120 -not -path "/tmp/proxy_files/*" -not -path "/tmp/.X*" -delete 2>/dev/null || true
    echo "[cleanup] Cleaned other old temporary files"
    
    echo "[cleanup] Cleanup completed, sleeping for $CLEANUP_INTERVAL seconds"
    sleep "$CLEANUP_INTERVAL"
  done
}

# Start updater in background
update_loop &

# Start cleanup in background
cleanup_loop &

# Exec main process
exec "$@"


