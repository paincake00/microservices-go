package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/paincake00/microservices-go/order/internal/mapper"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

const (
	grpcInventoryAddress = "localhost:50051"
	grpcPaymentAddress   = "localhost:50052"

	httpPort          = 8080
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrPartsNotFound = errors.New("parts not found")
	ErrOrderConflict = errors.New("order already paid")
)

type InMemoryStorage struct {
	mx     sync.RWMutex
	orders map[string]orderv1.OrderByUUIDResponse
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		orders: make(map[string]orderv1.OrderByUUIDResponse),
	}
}

func (o *InMemoryStorage) Save(order orderv1.OrderByUUIDResponse) {
	o.mx.Lock()
	defer o.mx.Unlock()

	o.orders[order.GetOrderUUID()] = order
}

func (o *InMemoryStorage) Get(orderUUID string) (orderv1.OrderByUUIDResponse, error) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	if order, ok := o.orders[orderUUID]; ok {
		return order, nil
	}
	return orderv1.OrderByUUIDResponse{}, ErrOrderNotFound
}

func (o *InMemoryStorage) Update(order orderv1.OrderByUUIDResponse) {
	o.Save(order)
}

type IOrderStorage interface {
	Save(orderv1.OrderByUUIDResponse)
	Get(string) (orderv1.OrderByUUIDResponse, error)
	Update(orderv1.OrderByUUIDResponse)
}

type OrderService struct {
	orderStorage     IOrderStorage
	inventoryService inventoryv1.InventoryServiceClient
	paymentService   paymentv1.PaymentServiceClient
}

func NewOrderService(
	orderStorage IOrderStorage,
	inventoryService inventoryv1.InventoryServiceClient,
	paymentService paymentv1.PaymentServiceClient,
) *OrderService {
	return &OrderService{
		orderStorage:     orderStorage,
		inventoryService: inventoryService,
		paymentService:   paymentService,
	}
}

func (or *OrderService) Create(ctx context.Context, userUUID string, partUUIDs []string) (
	orderv1.CreateOrderResponse,
	error,
) {
	resp, err := or.inventoryService.ListParts(
		ctx, &inventoryv1.ListPartsRequest{
			Filter: &inventoryv1.PartsFilter{
				Uuids: partUUIDs,
			},
		},
	)
	if err != nil {
		return orderv1.CreateOrderResponse{}, err
	}
	parts := resp.GetParts()
	if len(parts) < len(partUUIDs) {
		return orderv1.CreateOrderResponse{}, ErrPartsNotFound
	}

	orderUuid, err := getNewUUID()
	if err != nil {
		return orderv1.CreateOrderResponse{}, err
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	order := orderv1.OrderByUUIDResponse{
		OrderUUID:  orderUuid,
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     orderv1.OrderStatusEnumPENDINGPAYMENT,
	}

	log.Printf("Create order: %+v", order)

	or.orderStorage.Save(order)

	return orderv1.CreateOrderResponse{
		OrderUUID:  orderUuid,
		TotalPrice: totalPrice,
	}, nil
}

func (or *OrderService) GetByUuid(ctx context.Context, orderUUID string) (
	orderv1.OrderByUUIDResponse,
	error,
) {
	order, err := or.orderStorage.Get(orderUUID)
	if err != nil {
		return orderv1.OrderByUUIDResponse{}, err
	}
	return order, nil
}

func (or *OrderService) Pay(
	ctx context.Context,
	orderUUID string,
	payMethod orderv1.PaymentMethodEnum,
) (orderv1.PayOrderResponse, error) {
	order, err := or.orderStorage.Get(orderUUID)
	if err != nil {
		return orderv1.PayOrderResponse{}, err
	}

	payResponse, err := or.paymentService.PayOrder(
		ctx, &paymentv1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      order.UserUUID,
			PaymentMethod: mapper.MapPaymentMethodToGrpc(payMethod),
		},
	)
	if err != nil {
		return orderv1.PayOrderResponse{}, err
	}

	// Меняем поля у копии
	order.TransactionUUID.SetTo(payResponse.GetTransactionUuid())
	order.PaymentMethod.SetTo(payMethod)
	order.Status = orderv1.OrderStatusEnumPAID

	// Обновляем в карте заказ с новыми полями
	or.orderStorage.Update(order)

	log.Printf("Pay order: %+v", order)

	return orderv1.PayOrderResponse{TransactionUUID: order.TransactionUUID.Value}, nil
}

func (or *OrderService) Cancel(ctx context.Context, orderUUID string) error {
	order, err := or.orderStorage.Get(orderUUID)
	if err != nil {
		return err
	}

	if order.Status == orderv1.OrderStatusEnumPAID {
		return ErrOrderConflict
	}

	if order.Status == orderv1.OrderStatusEnumPENDINGPAYMENT {
		order.Status = orderv1.OrderStatusEnumCANCELLED
		or.orderStorage.Update(order)
	}
	return nil
}

type IOrderService interface {
	Create(context.Context, string, []string) (
		orderv1.CreateOrderResponse,
		error,
	)
	GetByUuid(context.Context, string) (
		orderv1.OrderByUUIDResponse,
		error,
	)
	Pay(
		context.Context,
		string,
		orderv1.PaymentMethodEnum,
	) (orderv1.PayOrderResponse, error)
	Cancel(context.Context, string) error
}

type OrderHandler struct {
	orderService IOrderService
}

func NewOrderHandler(orderService IOrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (oh *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (
	orderv1.CreateOrderRes,
	error,
) {
	resp, err := oh.orderService.Create(ctx, req.GetUserUUID(), req.GetPartUuids())
	if err != nil {
		log.Printf("Create order failed: %v", err)

		if errors.Is(err, ErrPartsNotFound) {
			return &orderv1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		return &orderv1.InternalError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &resp, nil
}

func (oh *OrderHandler) GetOrderByUuid(
	ctx context.Context,
	params orderv1.GetOrderByUuidParams,
) (orderv1.GetOrderByUuidRes, error) {
	order, err := oh.orderService.GetByUuid(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderv1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		return &orderv1.InternalError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &order, nil
}

func (oh *OrderHandler) PayOrder(
	ctx context.Context,
	req *orderv1.PayOrderRequest,
	params orderv1.PayOrderParams,
) (orderv1.PayOrderRes, error) {
	responseWithTransaction, err := oh.orderService.Pay(ctx, params.OrderUUID, req.GetPaymentMethod())
	if err != nil {
		log.Printf("Pay order failed: %v", err)

		if errors.Is(err, ErrOrderNotFound) {
			return &orderv1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		return &orderv1.InternalError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &responseWithTransaction, nil
}

func (oh *OrderHandler) CancelOrderByUuid(
	ctx context.Context,
	params orderv1.CancelOrderByUuidParams,
) (orderv1.CancelOrderByUuidRes, error) {
	err := oh.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("Cancel order failed: %v", err)

		if errors.Is(err, ErrOrderNotFound) {
			return &orderv1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, ErrOrderConflict) {
			return &orderv1.ConflictError{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		}

		return &orderv1.InternalError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &orderv1.CancelOrderByUuidNoContent{}, nil
}

// NewError создает новую ошибку в формате GenericError
func (oh *OrderHandler) NewError(_ context.Context, err error) *orderv1.GenericErrorStatusCode {
	return &orderv1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderv1.GenericError{
			Code:    orderv1.NewOptInt(http.StatusInternalServerError),
			Message: orderv1.NewOptString(err.Error()),
		},
	}
}

func getNewUUID() (string, error) {
	uuidBytes, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return uuidBytes.String(), nil
}

func createNewGrpcClient(address string) (*grpc.ClientConn, func() error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to gRPC Server with address %s: %v", address, err)
	}
	return conn, func() error { return conn.Close() }
}

func main() {
	// Создаем подключения к gRPC-серверам с функциями закрытия
	connInventory, closeConnInventory := createNewGrpcClient(grpcInventoryAddress)
	connPayment, closeConnPayment := createNewGrpcClient(grpcPaymentAddress)

	// Создаем сервисы-клиенты для доступа к gRPC-методам
	inventoryService := inventoryv1.NewInventoryServiceClient(connInventory)
	paymentService := paymentv1.NewPaymentServiceClient(connPayment)

	// Внедряем зависимости
	orderStorage := NewInMemoryStorage()
	orderService := NewOrderService(orderStorage, inventoryService, paymentService)
	orderHandler := NewOrderHandler(orderService)

	orderMux, err := orderv1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	// Создаем роутер Chi
	r := chi.NewRouter()

	// Добавляем middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderMux)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	notify := make(chan error, 1)

	go func() {
		log.Printf("Starting HTTP Server on port %d", httpPort)

		if err = srv.ListenAndServe(); err != nil {
			log.Printf("HTTP server get suddenly error: %v", err)

			notify <- err
		}

		close(notify)
	}()

	// WAITING

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-notify:
		if err != nil {
			log.Printf("Received suddenly error from notify: %v", err)
		}
	case sig := <-interrupt:
		log.Printf("Received interrupting signal: %v", sig)
	}

	// SHUTDOWN

	ctxShutDown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Printf("Shutting down server...")

	err = srv.Shutdown(ctxShutDown)
	if err != nil {
		log.Printf("Shutdown HTTP server failed: %v", err)
	}
	err = closeConnInventory()
	if err != nil {
		log.Printf("Close connection inventory failed: %v", err)
	}
	err = closeConnPayment()
	if err != nil {
		log.Printf("Close connection payment failed: %v", err)
	}

	log.Printf("Server shutdown.")
}
