FROM golang:1.24.5-alpine

WORKDIR /app

# Копируем только файлы модулей сначала
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/server

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]