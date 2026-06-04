-- name: GetProductImages :many
SELECT * FROM product_images
WHERE product_id = $1
ORDER BY is_main DESC, created_at ASC;

-- name: CountProductImages :one
SELECT COUNT(*) FROM product_images WHERE product_id = $1;

-- name: GetProductImage :one
SELECT * FROM product_images WHERE id = $1 AND product_id = $2;

-- name: AddProductImage :one
INSERT INTO product_images (product_id, url, is_main)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteProductImage :exec
DELETE FROM product_images WHERE id = $1 AND product_id = $2;

-- name: SetMainProductImage :exec
UPDATE product_images
SET is_main = (id = $1)
WHERE product_id = $2;
