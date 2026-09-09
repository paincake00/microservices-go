# ============================
# Stage 1: Build stage
# ============================

FROM golang:1.26.4-alpine as builder

# Устанавливаем git для загрузки зависимостей
RUN apk add --no-cache git

WORKDIR /app

# Копируем go.work и его контрольную сумму — это позволяет Go правильно связать модули
COPY go.work go.work.sum ./

# Копируем go.mod и go.sum всех зависимых модулей — нужно для кэширования go mod download
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY order/go.mod order/go.sum ./order/
COPY payment/go.mod payment/go.sum ./payment/

# Просто устанавливать зависимости
RUN go mod download

# Копируем исходный код всех модулей — включая зависимости, т.к. они нужны при компиляции
COPY platform ./platform
COPY shared ./shared
COPY inventory ./inventory
COPY order ./order
COPY payment ./payment

# Скачиваем grpc-health-probe — это утилита для проверки состояния gRPC-сервиса
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.57/grpc_health_probe-linux-amd64 ./grpc-health-probe

# Делаем скачанный файл исполняемым
RUN chmod +x grpc-health-probe

# Собираем бинарный файл Inventory-сервиса для Linux-архитектуры без CGO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-inventory ./inventory/cmd/inventoryapp/main.go


# ============================
# Stage 2: Final image
# ============================

# Используем минимальный стабильный образ Alpine Linux для финального контейнера
FROM alpine:3.21.3

# Тоже переходим в рабочую папку
WORKDIR /app

# Добавляем системную группу appgroup,
# а потом создаем непривилегированного пользователя appuser
# и добавляем в эту же группу. Чтобы не выполнять все из под root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Копируем скомпилированный бинарник из стадии builder
COPY --from=builder /app/app-inventory .

# Копируем бинарь grpc-health-probe в /bin
COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

# Меняем пользователя (chown - change owner) для этой директории рекурсивно (-R)
# (root->appuser)
RUN chown -R appuser:appgroup /app

# Запускаем приложение под non-root пользователем
USER appuser

# Экспонируем порт gRPC-сервиса Inventory
EXPOSE 50051

ENTRYPOINT ["./app-inventory"]