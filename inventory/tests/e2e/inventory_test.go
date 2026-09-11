//go:build e2e

package e2e

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/paincake00/microservices-go/platform/pkg/logger"
	inventoryv1 "github.com/paincake00/microservices-go/shared/pkg/proto/inventory/v1"
)

var _ = Describe(
	"Inventory Service", func() {
		var (
			ctx             context.Context
			cancel          context.CancelFunc
			conn            *grpc.ClientConn
			inventoryClient inventoryv1.InventoryServiceClient
		)

		BeforeEach(
			func() {
				ctx, cancel = context.WithCancel(suiteCtx)

				// Создаём gRPC клиент
				conn, err := grpc.NewClient(
					env.App.Address(),
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)
				Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

				inventoryClient = inventoryv1.NewInventoryServiceClient(conn)
			},
		)

		AfterEach(
			func() {
				// Чистим после теста
				cancel()

				if conn != nil {
					Expect(conn.Close()).To(Succeed())
				}
			},
		)

		Describe(
			"ListParts", func() {
				It(
					"Получить все детали", func() {
						parts, err := inventoryClient.ListParts(
							ctx, &inventoryv1.ListPartsRequest{
								Filter: &inventoryv1.PartsFilter{},
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(parts).ToNot(BeNil())
						Expect(parts.Parts).ToNot(BeEmpty())
						Expect(parts.Parts).To(HaveLen(5))

						logger.Info(ctx, "список деталей (uuids only)")
						for i, part := range parts.Parts {
							logger.Info(ctx, "деталь", zap.Int("номер", i+1), zap.String("part", part.Uuid))
						}
					},
				)
			},
		)
	},
)
