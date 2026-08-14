-- +goose Down
DROP TABLE IF EXISTS order;
DROP TABLE IF EXISTS order_part;
DROP TYPE order_status_enum;
DROP TYPE payment_method_enum;