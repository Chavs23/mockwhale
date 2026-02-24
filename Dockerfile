# Этап 1: Сборка бинарника
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем статический бинарник без CGO
RUN CGO_ENABLED=0 GOOS=linux go build -o mockwhale ./cmd/api/main.go

# Этап 2: Финальный образ
FROM alpine:latest

WORKDIR /root/

# Копируем только бинарник и необходимые файлы из этапа сборки
COPY --from=builder /app/mockwhale .
COPY --from=builder /app/migrations ./migrations

# Открываем порт
EXPOSE 3000

# Запускаем приложение
CMD ["./mockwhale"]
