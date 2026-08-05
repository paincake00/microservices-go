package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	signalGo "os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

const (
	grpcServerPort  = 50051
	defaultPartsNum = 5
)

var ErrNotFound = errors.New("not found")

type PartService struct {
	mx    sync.RWMutex
	parts map[string]*inventoryv1.Part
}

func NewPartService() *PartService {
	parts := make(map[string]*inventoryv1.Part, defaultPartsNum)

	categories, length := allCategories()

	for range defaultPartsNum {
		id, err := getNewUUID()
		if err != nil {
			log.Fatal(err)
		}

		//nolint:gosec // random is used only for test data generation
		category := categories[rand.Intn(length)]

		parts[id] = &inventoryv1.Part{
			Uuid:          id,
			Name:          gofakeit.Name(),
			Description:   gofakeit.ProductDescription(),
			Price:         gofakeit.Float64(),
			StockQuantity: gofakeit.Int64(),
			Category:      category,
			Dimensions: &inventoryv1.Dimensions{
				Height: gofakeit.Float64(),
				Width:  gofakeit.Float64(),
				Length: gofakeit.Float64(),
				Weight: gofakeit.Float64(),
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    gofakeit.Name(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: []string{
				gofakeit.Word(),
				gofakeit.Word(),
				gofakeit.Word(),
			},
			Metadata:  randomMetadata(),
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		}

		log.Printf("Add new part with UUID: %s", parts[id].Uuid)
	}

	return &PartService{
		parts: parts,
	}
}

type ListPartsFilter struct {
	Uuids                 map[string]struct{}
	Names                 map[string]struct{}
	Categories            map[inventoryv1.Category]struct{}
	ManufacturerCountries map[string]struct{}
	Tags                  map[string]struct{}
}

func NewListPartsFilter(
	uuids []string,
	names []string,
	categories []inventoryv1.Category,
	manufacturerCountries []string,
	tags []string,
) ListPartsFilter {
	return ListPartsFilter{
		Uuids:                 createSet(uuids),
		Names:                 createSet(names),
		Categories:            createSet(categories),
		ManufacturerCountries: createSet(manufacturerCountries),
		Tags:                  createSet(tags),
	}
}

func createSet[T comparable](items []T) map[T]struct{} {
	var set map[T]struct{}
	if len(items) != 0 {
		set = make(map[T]struct{}, len(items))
		for _, item := range items {
			set[item] = struct{}{}
		}
	}
	return set
}

func (f *ListPartsFilter) IsEmpty() bool {
	return len(f.Uuids) == 0 && len(f.Names) == 0 && len(f.Categories) == 0 && len(f.ManufacturerCountries) == 0 && len(f.Tags) == 0
}

func (f *ListPartsFilter) Match(part *inventoryv1.Part) bool {
	if _, ok := f.Uuids[part.Uuid]; len(f.Uuids) != 0 && !ok {
		return false
	}
	if _, ok := f.Names[part.Name]; len(f.Names) != 0 && !ok {
		return false
	}
	if _, ok := f.Categories[part.Category]; len(f.Categories) != 0 && !ok {
		return false
	}
	if _, ok := f.ManufacturerCountries[part.Manufacturer.Country]; len(f.ManufacturerCountries) != 0 && !ok {
		return false
	}
	if len(f.Tags) != 0 && !slices.ContainsFunc(
		part.Tags, func(s string) bool {
			_, ok := f.Tags[s]
			return ok
		},
	) {
		return false
	}
	return true
}

func (p *PartService) ListParts(filter ListPartsFilter) []*inventoryv1.Part {
	p.mx.RLock()
	defer p.mx.RUnlock()

	res := make([]*inventoryv1.Part, 0)

	emptyFilter := filter.IsEmpty()

	for _, part := range p.parts {
		if emptyFilter || filter.Match(part) {
			res = append(res, part)
		}
	}

	return res
}

func (p *PartService) GetPart(id string) (*inventoryv1.Part, error) {
	p.mx.RLock()
	defer p.mx.RUnlock()

	part, ok := p.parts[id]
	if !ok {
		return nil, ErrNotFound
	}
	return part, nil
}

func getNewUUID() (string, error) {
	uuidBytes, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return uuidBytes.String(), nil
}

func allCategories() ([]inventoryv1.Category, int) {
	res := make([]inventoryv1.Category, 0, len(inventoryv1.Category_name))

	for name := range inventoryv1.Category_name {
		if inventoryv1.Category(name) == inventoryv1.Category_CATEGORY_UNSPECIFIED {
			continue
		}

		res = append(res, inventoryv1.Category(name))
	}

	return res, len(res)
}

func randomMetadata() map[string]*inventoryv1.Value {
	return map[string]*inventoryv1.Value{
		"color": {
			Kind: &inventoryv1.Value_StringValue{
				StringValue: gofakeit.Color(),
			},
		},
		"count": {
			Kind: &inventoryv1.Value_Int64Value{
				Int64Value: gofakeit.Int64(),
			},
		},
		"rating": {
			Kind: &inventoryv1.Value_DoubleValue{
				DoubleValue: gofakeit.Float64(),
			},
		},
		"enabled": {
			Kind: &inventoryv1.Value_BoolValue{
				BoolValue: gofakeit.Bool(),
			},
		},
	}
}

type IPartService interface {
	GetPart(string) (*inventoryv1.Part, error)
	ListParts(ListPartsFilter) []*inventoryv1.Part
}

type InventoryHandler struct {
	inventoryv1.UnimplementedInventoryServiceServer

	partService IPartService
}

func NewInventoryHandler(partService IPartService) *InventoryHandler {
	return &InventoryHandler{
		partService: partService,
	}
}

func (i *InventoryHandler) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (
	*inventoryv1.GetPartResponse,
	error,
) {
	id := req.GetUuid()

	part, err := i.partService.GetPart(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Printf("Part %s not found", id)

			return nil, status.Errorf(codes.NotFound, "not found")
		}

		log.Printf("Error getting part %s: %s", id, err)

		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &inventoryv1.GetPartResponse{
		Part: part,
	}, nil
}

func (i *InventoryHandler) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	protoFilter := req.GetFilter()

	filter := NewListPartsFilter(
		protoFilter.GetUuids(), protoFilter.GetNames(), protoFilter.GetCategories(),
		protoFilter.GetManufacturerCountries(), protoFilter.GetTags(),
	)

	return &inventoryv1.ListPartsResponse{
		Parts: i.partService.ListParts(filter),
	}, nil
}

func main() {
	partService := NewPartService()
	inventoryHandler := NewInventoryHandler(partService)

	srv := grpc.NewServer()

	inventoryv1.RegisterInventoryServiceServer(srv, inventoryHandler)

	reflection.Register(srv)

	// Канал для уведомления об ошибках
	notify := make(chan error, 1)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcServerPort))
		if err != nil {
			log.Printf("failed to listen: %v", err)

			notify <- err
			close(notify)
			return
		}

		log.Printf("starting gRPC server on port %d", grpcServerPort)

		if err = srv.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)

			notify <- err
		}

		close(notify)
	}()

	interrupt := make(chan os.Signal, 1)
	signalGo.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		log.Printf("Received interrupting signal: %v", sig)
	case errNotify := <-notify:
		log.Printf("Received notify error: %v", errNotify)
	}

	log.Println("Shutting down gRPC server...")

	srv.GracefulStop()
}
