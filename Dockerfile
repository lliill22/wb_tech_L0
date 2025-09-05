# ---------- Стадия сборки ----------
FROM golang:1.24.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN go build -o main ./cmd

# ---------- Стадия выполнения ----------
FROM debian:12-slim

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/main .

# Копируем шаблоны (и при необходимости статику)
COPY --from=builder /app/templates ./templates
# COPY --from=builder /app/static ./static  # если понадобится

CMD ["./main"]
