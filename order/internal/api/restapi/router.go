package restapi

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/paincake00/microservices-go/order/internal/api/restapi/health"
	v1 "github.com/paincake00/microservices-go/order/internal/api/restapi/order/v1"
	"github.com/paincake00/microservices-go/order/internal/service"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
)

func NewRouter(orderService service.IOrderService, reqTimeout time.Duration) (*chi.Mux, error) {
	// Создаем роутер Chi
	r := chi.NewRouter()

	// Добавляем middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(reqTimeout))

	orderHandler := v1.NewOrderHandler(orderService)

	orderMux, err := orderv1.NewServer(orderHandler)
	if err != nil {
		return nil, fmt.Errorf("error create HTTP-mux from OpenAPI specs: %w", err)
	}

	// Добавляем health-эндпоинт
	r.Get("/health", health.CheckHealth)
	// Монтируем обработчики OpenAPI
	r.Mount("/", orderMux)

	return r, nil
}
