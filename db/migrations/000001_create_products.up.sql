CREATE TABLE products (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    description TEXT       NOT NULL DEFAULT '',
    price      NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    stock      INT         NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
