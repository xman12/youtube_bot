FROM golang:1.23-alpine AS builder

WORKDIR /app

# Установка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения из папки cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Финальный образ
FROM alpine:latest

# Base packages and yt-dlp
RUN apk --no-cache add \
    ca-certificates \
    wget \
    bash \
    ffmpeg \
    yt-dlp

WORKDIR /root/

# Копирование бинарного файла из builder
COPY --from=builder /app/main .

# Создание пустого .env файла если его нет
RUN touch .env

# Entry and auto-update config
COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
ENV YTDLP_UPDATE_INTERVAL_SECONDS=21600

EXPOSE 8090

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["./main"]
