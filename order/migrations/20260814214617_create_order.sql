-- +goose Up
CREATE TYPE payment_method_enum AS ENUM (
    'PAYMENT_METHOD_UNSPECIFIED',
    'PAYMENT_METHOD_CARD',
    'PAYMENT_METHOD_SBP',
    'PAYMENT_METHOD_CREDIT_CARD',
    'PAYMENT_METHOD_INVESTOR_MONEY'
);

CREATE TYPE order_status_enum AS ENUM (
    'PENDING_PAYMENT',
    'PAID',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS "order"
(
    order_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    total_price NUMERIC NOT NULL,
    transaction_uuid UUID,
    payment_method payment_method_enum,
    status order_status_enum NOT NULL DEFAULT 'PENDING_PAYMENT'
);

CREATE TABLE IF NOT EXISTS order_part
(
    order_uuid UUID NOT NULL REFERENCES "order"(order_uuid) ON DELETE CASCADE,
    part_uuid UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS order_part;
DROP TABLE IF EXISTS "order";
DROP TYPE payment_method_enum;
DROP TYPE order_status_enum;