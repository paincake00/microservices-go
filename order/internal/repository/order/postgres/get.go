package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/repository/model"
)

func (o *PostgresOrderStorage) Get(ctx context.Context, orderUUID string) (entity.Order, error) {
	const op = "postgres.Get"

	// Начало транзакции
	tx, err := o.pgClient.Pool.Begin(ctx)
	if err != nil {
		return entity.Order{}, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer func() {
		//nolint:gosec // G104: rollback error intentionally ignored
		_ = tx.Rollback(ctx)
	}()

	// Получение заказа по UUID
	builderOrderQuery := o.pgClient.Builder.
		Select("order_uuid", "user_uuid", "total_price", "transaction_uuid", "payment_method", "status").
		From(`"order"`).
		Where(squirrel.Eq{"order_uuid": orderUUID})

	orderQuery, args, err := builderOrderQuery.ToSql()
	if err != nil {
		return entity.Order{}, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	row := tx.QueryRow(ctx, orderQuery, args...)
	var order entity.Order
	err = row.Scan(
		&order.OrderUuid, &order.UserUuid, &order.TotalPrice, &order.TransactionUuid, &order.PaymentMethod,
		&order.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Order{}, model.ErrOrderNotFound
		}
		return entity.Order{}, fmt.Errorf("%s: failed to scan row: %w", op, err)
	}

	// Получение списка UUIDs деталей заказа
	builderOrderPartQuery := o.pgClient.Builder.
		Select("part_uuid").
		From("order_part").
		Where(squirrel.Eq{"order_uuid": orderUUID})

	orderQuery, args, err = builderOrderPartQuery.ToSql()
	if err != nil {
		return entity.Order{}, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	rows, err := tx.Query(ctx, orderQuery, args...)
	if err != nil {
		return entity.Order{}, err
	}
	defer rows.Close()

	parts := make([]string, 0)
	for rows.Next() {
		var partUuid string
		err = rows.Scan(&partUuid)
		if err != nil {
			return entity.Order{}, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}
		parts = append(parts, partUuid)
	}
	order.PartUuids = parts

	// Завершение транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return entity.Order{}, fmt.Errorf("%s: failed to commit: %w", op, err)
	}
	return order, nil
}
