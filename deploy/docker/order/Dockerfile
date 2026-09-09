FROM golang:1.26.4-alpine as builder

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

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-order ./order/cmd/orderapp/main.go


FROM alpine:3.21.3

RUN apk add --no-cache curl

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder app/app-order .
COPY --from=builder app/order/migrations ./migrations

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["./app-order"]