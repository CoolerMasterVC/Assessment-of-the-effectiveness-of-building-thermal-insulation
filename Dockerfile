# Dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o lab1-app ./cmd/server

# Final stage
FROM alpine:latest

# Устанавливаем необходимые пакеты для работы
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем бинарник из стадии сборки
COPY --from=builder /app/lab1-app .

# Копируем шаблоны и статические файлы (исправляем пути)
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Копируем конфигурационные файлы (если есть)
COPY --from=builder /app/config/config.toml ./config.toml

# Создаем пустой .env если его нет (чтобы избежать ошибки)
RUN touch .env

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./lab1-app"]