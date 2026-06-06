CREATE TABLE orders (
    id           BIGSERIAL        PRIMARY KEY,
    user_id      BIGINT           NOT NULL REFERENCES users(id),
    status       VARCHAR(20)      NOT NULL DEFAULT 'pending',
    total_amount DOUBLE PRECISION NOT NULL,
    created_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id           BIGSERIAL        PRIMARY KEY,
    order_id     BIGINT           NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id   BIGINT           NOT NULL REFERENCES products(id),
    product_name VARCHAR(200)     NOT NULL,  -- price/name snapshot at order time
    quantity     INT              NOT NULL CHECK(quantity > 0),
    unit_price   DOUBLE PRECISION NOT NULL,
    created_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id       ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
