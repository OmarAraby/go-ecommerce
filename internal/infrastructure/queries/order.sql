-- name: CreateOrder :one
INSERT INTO orders (user_id, status, total_amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: GetUserOrder :one
SELECT * FROM orders WHERE id = $1 AND user_id = $2;

-- name: ListUserOrders :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserOrders :one
SELECT COUNT(*) FROM orders WHERE user_id = $1;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, product_name, quantity, unit_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetOrderItems :many
SELECT * FROM order_items WHERE order_id = $1 ORDER BY id;
