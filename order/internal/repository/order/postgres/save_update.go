package memory

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/paincake00/microservices-go/order/internal/entity"
	"github.com/paincake00/microservices-go/order/internal/repository/model"
)

func (o *PostgresOrderStorage) Save(ctx context.Context, order entity.Order) (string, error) {
	const op = "postgres.Save"

	tx, err := o.pgClient.Pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("%s: error starting transaction: %w", op, err)
	}
	defer func() {
		//nolint:gosec // G104: rollback error intentionally ignored
		_ = tx.Rollback(ctx)
	}()

	// Создание заказа
	builderOrderQuery := o.pgClient.Builder.
		Insert(`"order"`).
		Columns("user_uuid", "total_price", "transaction_uuid", "payment_method", "status").
		Values(
			order.UserUuid, order.TotalPrice, order.TransactionUuid, order.PaymentMethod, order.Status,
		).
		Suffix("RETURNING order_uuid")

	orderQuery, args, err := builderOrderQuery.ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: error building order query: %w", op, err)
	}

	var orderUUID string
	err = tx.QueryRow(ctx, orderQuery, args...).Scan(&orderUUID)
	if err != nil {
		return "", fmt.Errorf("%s: error executing order query: %w", op, err)
	}

	// Сохранение UUIDs деталей заказа
	err = o.insertOrderParts(op, ctx, tx, orderUUID, order.PartUuids)
	if err != nil {
		return "", err
	}

	// Завершение транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("%s: error committing order transaction: %w", op, err)
	}
	return orderUUID, nil
}

func (o *PostgresOrderStorage) Update(ctx context.Context, order entity.Order) error {
	const op = "postgres.Update"

	tx, err := o.pgClient.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: error starting transaction: %w", op, err)
	}
	defer func() {
		//nolint:gosec // G104: rollback error intentionally ignored
		_ = tx.Rollback(ctx)
	}()

	// Обновление заказа
	builderOrderQuery := o.pgClient.Builder.
		Update(`"order"`).
		Set("user_uuid", order.UserUuid).
		Set("total_price", order.TotalPrice).
		Set("transaction_uuid", order.TransactionUuid).
		Set("payment_method", order.PaymentMethod).
		Set("status", order.Status).
		Where(squirrel.Eq{"order_uuid": order.OrderUuid})

	orderQuery, args, err := builderOrderQuery.ToSql()
	if err != nil {
		return fmt.Errorf("%s: error building order query: %w", op, err)
	}

	res, err := tx.Exec(ctx, orderQuery, args...)
	if err != nil {
		return fmt.Errorf("%s: error executing order query: %w", op, err)
	}
	if res.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	// Обновление деталей заказа

	// Удаление предыдущих UUIDs деталей заказа
	err = o.deleteOrderParts(op, ctx, tx, order.OrderUuid)
	if err != nil {
		return err
	}
	// Сохранение UUIDs деталей заказа
	err = o.insertOrderParts(op, ctx, tx, order.OrderUuid, order.PartUuids)
	if err != nil {
		return err
	}

	// Завершение транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: error committing order transaction: %w", op, err)
	}
	return nil
}

func (o *PostgresOrderStorage) insertOrderParts(
	op string, ctx context.Context,
	tx pgx.Tx,
	orderUUID string,
	partUuids []string,
) error {
	if len(partUuids) == 0 {
		return nil
	}

	builderOrderPartQuery := o.pgClient.Builder.
		Insert("order_part").
		Columns("order_uuid", "part_uuid")

	for _, partUuid := range partUuids {
		builderOrderPartQuery = builderOrderPartQuery.Values(orderUUID, partUuid)
	}

	orderPartQuery, args, err := builderOrderPartQuery.ToSql()
	if err != nil {
		return fmt.Errorf("%s: error building order part query: %w", op, err)
	}

	_, err = tx.Exec(ctx, orderPartQuery, args...)
	if err != nil {
		return fmt.Errorf("%s: error executing order part query: %w", op, err)
	}

	return nil
}

func (o *PostgresOrderStorage) deleteOrderParts(
	op string, ctx context.Context,
	tx pgx.Tx,
	orderUUID string,
) error {
	builderOrderPartQuery := o.pgClient.Builder.
		Delete("order_part").
		Where(squirrel.Eq{"order_uuid": orderUUID})

	orderPartQuery, args, err := builderOrderPartQuery.ToSql()
	if err != nil {
		return fmt.Errorf("%s: error building order part query: %w", op, err)
	}

	_, err = tx.Exec(ctx, orderPartQuery, args...)
	if err != nil {
		return fmt.Errorf("%s: error executing order part query: %w", op, err)
	}

	return nil
}
