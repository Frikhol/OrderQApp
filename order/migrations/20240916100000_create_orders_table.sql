-- +goose Up
-- SQL in this section is executed when the migration is applied

CREATE TABLE IF NOT EXISTS orders (
    order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    order_address TEXT NOT NULL,
    order_location TEXT NOT NULL,
    order_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    order_time_gap INTERVAL NOT NULL,
    order_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(order_status);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back

DROP TABLE IF EXISTS orders; 