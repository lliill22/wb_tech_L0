# ---------- Стадия сборки ----------
FROM golang:1.24.5 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum (для кэширования зависимостей)
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем бинарник
RUN go build -o main ./cmd

# ---------- Стадия выполнения ----------
FROM debian:12-slim

WORKDIR /app

# Копируем бинарник из стадии сборки
COPY --from=builder /app/main .

# Запускаем приложение
CMD ["./main"]
