# API Test Cases — Go E-Commerce

> Base URL: `http://localhost:8080`
> After login, save `access_token` as `{{token}}` and use it in all protected endpoints.

---

## Variables

| Variable | Value |
|----------|-------|
| `{{base_url}}` | `http://localhost:8080` |
| `{{token}}` | from login response |
| `{{refresh_token}}` | from login response |
| `{{reset_token}}` | from forgot-password response |

---

## 1. Health

### GET /health
```
GET {{base_url}}/health
```
**Expected: 200**
```json
{ "status": "ok" }
```

---

## 2. Auth

### POST /auth/register — success
```
POST {{base_url}}/auth/register
Content-Type: application/json

{
  "name": "Omar Araby",
  "email": "omar@test.com",
  "password": "password123"
}
```
**Expected: 201**
```json
{
  "id": 1,
  "name": "Omar Araby",
  "email": "omar@test.com",
  "role": "user",
  "created_at": "...",
  "updated_at": "..."
}
```
> Note: no `password` field in response

---

### POST /auth/register — duplicate email
```
POST {{base_url}}/auth/register
Content-Type: application/json

{
  "name": "Omar Araby",
  "email": "omar@test.com",
  "password": "password123"
}
```
**Expected: 409**
```json
{ "code": "CONFLICT", "message": "email already registered" }
```

---

### POST /auth/register — validation errors
```
POST {{base_url}}/auth/register
Content-Type: application/json

{
  "name": "O",
  "email": "not-an-email",
  "password": "short"
}
```
**Expected: 422**
```json
{
  "code": "VALIDATION_ERROR",
  "message": "validation failed",
  "errors": {
    "name": "must be at least 2 characters",
    "email": "must be a valid email address",
    "password": "must be at least 8 characters"
  }
}
```

---

### POST /auth/register — missing fields
```
POST {{base_url}}/auth/register
Content-Type: application/json

{}
```
**Expected: 422**
```json
{
  "code": "VALIDATION_ERROR",
  "errors": {
    "name": "required",
    "email": "required",
    "password": "required"
  }
}
```

---

### POST /auth/login — success
```
POST {{base_url}}/auth/login
Content-Type: application/json

{
  "email": "omar@test.com",
  "password": "password123"
}
```
**Expected: 200**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "abc123...",
  "user": {
    "id": 1,
    "name": "Omar Araby",
    "email": "omar@test.com",
    "role": "user"
  }
}
```
> Save `access_token` as `{{token}}` and `refresh_token` as `{{refresh_token}}`

---

### POST /auth/login — wrong password
```
POST {{base_url}}/auth/login
Content-Type: application/json

{
  "email": "omar@test.com",
  "password": "wrongpassword"
}
```
**Expected: 401**
```json
{ "code": "UNAUTHORIZED", "message": "invalid email or password" }
```

---

### POST /auth/login — wrong email
```
POST {{base_url}}/auth/login
Content-Type: application/json

{
  "email": "notfound@test.com",
  "password": "password123"
}
```
**Expected: 401**
```json
{ "code": "UNAUTHORIZED", "message": "invalid email or password" }
```

---

### POST /auth/refresh — success
```
POST {{base_url}}/auth/refresh
Content-Type: application/json

{
  "refresh_token": "{{refresh_token}}"
}
```
**Expected: 200**
```json
{
  "access_token": "eyJ... (new)",
  "refresh_token": "xyz... (new)"
}
```
> Old tokens are revoked (token rotation). Save the new tokens.

---

### POST /auth/refresh — invalid token
```
POST {{base_url}}/auth/refresh
Content-Type: application/json

{
  "refresh_token": "invalid-token-here"
}
```
**Expected: 401**
```json
{ "code": "UNAUTHORIZED", "message": "invalid or expired refresh token" }
```

---

### POST /auth/logout
```
POST {{base_url}}/auth/logout
Content-Type: application/json

{
  "refresh_token": "{{refresh_token}}"
}
```
**Expected: 204** (no body)

---

### POST /auth/forgot-password — registered email
```
POST {{base_url}}/auth/forgot-password
Content-Type: application/json

{
  "email": "omar@test.com"
}
```
**Expected: 200**
```json
{
  "message": "if this email is registered you will receive a reset link",
  "reset_token": "abc123..."
}
```
> Save `reset_token` as `{{reset_token}}`
> In production the token is sent via email only — not in the response

---

### POST /auth/forgot-password — unregistered email
```
POST {{base_url}}/auth/forgot-password
Content-Type: application/json

{
  "email": "ghost@test.com"
}
```
**Expected: 200** (same response — does not reveal if email exists)
```json
{ "message": "if this email is registered you will receive a reset link" }
```

---

### POST /auth/reset-password — success
```
POST {{base_url}}/auth/reset-password
Content-Type: application/json

{
  "token": "{{reset_token}}",
  "new_password": "newpassword123"
}
```
**Expected: 200**
```json
{ "message": "password reset successfully" }
```

---

### POST /auth/reset-password — token already used
```
POST {{base_url}}/auth/reset-password
Content-Type: application/json

{
  "token": "{{reset_token}}",
  "new_password": "anotherpassword"
}
```
**Expected: 400**
```json
{ "code": "BAD_REQUEST", "message": "invalid or expired reset token" }
```

---

## 3. Users — Protected (Authorization: Bearer {{token}})

### GET /users/me
```
GET {{base_url}}/users/me
Authorization: Bearer {{token}}
```
**Expected: 200**
```json
{
  "id": 1,
  "name": "Omar Araby",
  "email": "omar@test.com",
  "role": "user",
  "created_at": "...",
  "updated_at": "..."
}
```

---

### GET /users/me — no token
```
GET {{base_url}}/users/me
```
**Expected: 401**
```json
{ "code": "UNAUTHORIZED", "message": "missing or invalid token" }
```

---

### PUT /users/me — update name
```
PUT {{base_url}}/users/me
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "Omar Updated"
}
```
**Expected: 200** — user object with updated name

---

### PUT /users/me — validation error
```
PUT {{base_url}}/users/me
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "O"
}
```
**Expected: 422**
```json
{
  "code": "VALIDATION_ERROR",
  "errors": { "name": "must be at least 2 characters" }
}
```

---

### PUT /users/me/email — success
```
PUT {{base_url}}/users/me/email
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "new_email": "newemail@test.com",
  "current_password": "password123"
}
```
**Expected: 200** — user object with updated email

---

### PUT /users/me/email — wrong password
```
PUT {{base_url}}/users/me/email
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "new_email": "newemail@test.com",
  "current_password": "wrongpassword"
}
```
**Expected: 401**
```json
{ "code": "UNAUTHORIZED", "message": "incorrect password" }
```

---

### PUT /users/me/password — success
```
PUT {{base_url}}/users/me/password
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "current_password": "password123",
  "new_password": "newsecurepass123"
}
```
**Expected: 200**
```json
{ "message": "password changed successfully" }
```

---

### PUT /users/me/password — wrong current password
```
PUT {{base_url}}/users/me/password
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "current_password": "wrongpassword",
  "new_password": "newsecurepass123"
}
```
**Expected: 401**
```json
{ "code": "UNAUTHORIZED", "message": "incorrect current password" }
```

---

## 4. Products

### GET /products — default
```
GET {{base_url}}/products
```
**Expected: 200**
```json
{
  "data": [...],
  "page": 1,
  "limit": 20,
  "total": 5,
  "total_pages": 1
}
```

---

### GET /products — pagination
```
GET {{base_url}}/products?page=1&limit=2
```
**Expected: 200** — at most 2 products in `data`

---

### GET /products — filter by name
```
GET {{base_url}}/products?name=shirt
```
**Expected: 200** — products containing "shirt" in name (case-insensitive)

---

### GET /products — filter by price range
```
GET {{base_url}}/products?min_price=10&max_price=50
```
**Expected: 200** — products with price between 10 and 50

---

### GET /products — no results
```
GET {{base_url}}/products?min_price=99999
```
**Expected: 200**
```json
{ "data": [], "page": 1, "limit": 20, "total": 0, "total_pages": 0 }
```

---

### GET /products — sort by price ascending
```
GET {{base_url}}/products?sort=price&order=asc
```
**Expected: 200** — cheapest first

---

### GET /products — sort by price descending
```
GET {{base_url}}/products?sort=price&order=desc
```
**Expected: 200** — most expensive first

---

### GET /products — combined query
```
GET {{base_url}}/products?page=1&limit=5&sort=price&order=asc&min_price=10&name=go
```
**Expected: 200**

---

### GET /products/:id — success
```
GET {{base_url}}/products/1
```
**Expected: 200**
```json
{
  "id": 1,
  "name": "Go T-Shirt v2",
  "description": "Updated design",
  "price": 34.99,
  "stock": 78,
  "images": [
    { "id": 1, "url": "/uploads/products/1_123.png", "is_main": true }
  ],
  "created_at": "...",
  "updated_at": "..."
}
```

---

### GET /products/:id — not found
```
GET {{base_url}}/products/99999
```
**Expected: 404**
```json
{ "code": "NOT_FOUND", "message": "product not found" }
```

---

### GET /products/:id — invalid id
```
GET {{base_url}}/products/abc
```
**Expected: 400**
```json
{ "code": "BAD_REQUEST", "message": "invalid product id" }
```

---

### POST /products — success
```
POST {{base_url}}/products
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "Go Hoodie",
  "description": "Warm Go branded hoodie",
  "price": 59.99,
  "stock": 50
}
```
**Expected: 201**
```json
{
  "id": 2,
  "name": "Go Hoodie",
  "price": 59.99,
  "stock": 50,
  "images": []
}
```

---

### POST /products — validation errors
```
POST {{base_url}}/products
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "",
  "price": -5,
  "stock": -1
}
```
**Expected: 422**
```json
{
  "code": "VALIDATION_ERROR",
  "errors": {
    "name": "required",
    "price": "must be > 0",
    "stock": "must be >= 0"
  }
}
```

---

### PUT /products/:id — success
```
PUT {{base_url}}/products/1
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "Go T-Shirt v3",
  "description": "New design 2026",
  "price": 39.99,
  "stock": 100
}
```
**Expected: 200** — updated product

---

### PUT /products/:id — not found
```
PUT {{base_url}}/products/99999
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "Test",
  "price": 10,
  "stock": 1
}
```
**Expected: 404**

---

### DELETE /products/:id — success
```
DELETE {{base_url}}/products/2
Authorization: Bearer {{token}}
```
**Expected: 204** (no body)

---

### DELETE /products/:id — no token
```
DELETE {{base_url}}/products/1
```
**Expected: 401**

---

## 5. Product Images — Protected

### POST /products/:id/images — success (first image → auto main)
```
POST {{base_url}}/products/1/images
Authorization: Bearer {{token}}
Content-Type: multipart/form-data

field name : "image"
file       : any .png or .jpg file
```
**Expected: 200**
```json
{
  "id": 1,
  "images": [
    { "id": 1, "url": "/uploads/products/1_123.png", "is_main": true }
  ]
}
```

---

### POST /products/:id/images — second image (not main)
```
POST {{base_url}}/products/1/images
Authorization: Bearer {{token}}
Content-Type: multipart/form-data

field: "image" = second image file
```
**Expected: 200**
```json
{
  "images": [
    { "id": 1, "url": "...", "is_main": true },
    { "id": 2, "url": "...", "is_main": false }
  ]
}
```

---

### POST /products/:id/images — wrong file type
```
POST {{base_url}}/products/1/images
Authorization: Bearer {{token}}
Content-Type: multipart/form-data

field: "image" = a .txt or .pdf file
```
**Expected: 400**
```json
{ "code": "BAD_REQUEST", "message": "only jpeg, png, webp, gif images are allowed" }
```

---

### POST /products/:id/images — missing field
```
POST {{base_url}}/products/1/images
Authorization: Bearer {{token}}
Content-Type: multipart/form-data

(no field named "image")
```
**Expected: 400**
```json
{ "code": "BAD_REQUEST", "message": "field 'image' is required" }
```

---

### POST /products/:id/images — exceed 6 image limit
> Upload 7 images to the same product (run this after already having 6)

**Expected: 400** (on the 7th)
```json
{ "code": "BAD_REQUEST", "message": "image limit reached" }
```

---

### PUT /products/:id/images/:imageId/main — set main image
```
PUT {{base_url}}/products/1/images/2/main
Authorization: Bearer {{token}}
```
**Expected: 200** — image 2 is `is_main: true`, image 1 is `is_main: false`

---

### PUT /products/:id/images/:imageId/main — image not found
```
PUT {{base_url}}/products/1/images/9999/main
Authorization: Bearer {{token}}
```
**Expected: 404**
```json
{ "code": "NOT_FOUND", "message": "image not found" }
```

---

### DELETE /products/:id/images/:imageId — success
```
DELETE {{base_url}}/products/1/images/1
Authorization: Bearer {{token}}
```
**Expected: 200** — product with remaining images

---

### DELETE /products/:id/images/:imageId — deleting the main image
> Delete the image that has `is_main: true`

**Expected: 200** — next remaining image is automatically promoted to main

---

## 6. Orders — Protected

### POST /orders — success
```
POST {{base_url}}/orders
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "items": [
    { "product_id": 1, "quantity": 2 },
    { "product_id": 2, "quantity": 1 }
  ]
}
```
**Expected: 201**
```json
{
  "id": 1,
  "status": "pending",
  "total_amount": 129.97,
  "items": [
    {
      "id": 1,
      "product_id": 1,
      "product_name": "Go T-Shirt v2",
      "quantity": 2,
      "unit_price": 34.99,
      "subtotal": 69.98
    },
    {
      "id": 2,
      "product_id": 2,
      "product_name": "Go Hoodie",
      "quantity": 1,
      "unit_price": 59.99,
      "subtotal": 59.99
    }
  ],
  "created_at": "..."
}
```
> Verify that product stock decreased after this request

---

### POST /orders — insufficient stock
```
POST {{base_url}}/orders
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "items": [
    { "product_id": 1, "quantity": 9999 }
  ]
}
```
**Expected: 400**
```json
{ "code": "BAD_REQUEST", "message": "insufficient stock for one or more items" }
```

---

### POST /orders — product not found
```
POST {{base_url}}/orders
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "items": [
    { "product_id": 99999, "quantity": 1 }
  ]
}
```
**Expected: 404**
```json
{ "code": "NOT_FOUND", "message": "one or more products not found" }
```

---

### POST /orders — empty items
```
POST {{base_url}}/orders
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "items": []
}
```
**Expected: 422**
```json
{
  "code": "VALIDATION_ERROR",
  "errors": { "items": "must be >= 1" }
}
```

---

### POST /orders — invalid quantity
```
POST {{base_url}}/orders
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "items": [
    { "product_id": 1, "quantity": 0 }
  ]
}
```
**Expected: 422**
```json
{
  "code": "VALIDATION_ERROR",
  "errors": { "items[0].quantity": "must be >= 1" }
}
```

---

### POST /orders — no token
```
POST {{base_url}}/orders
Content-Type: application/json

{ "items": [{ "product_id": 1, "quantity": 1 }] }
```
**Expected: 401**

---

### GET /orders — list my orders
```
GET {{base_url}}/orders
Authorization: Bearer {{token}}
```
**Expected: 200**
```json
{
  "data": [...],
  "page": 1,
  "limit": 20,
  "total": 1,
  "total_pages": 1
}
```

---

### GET /orders — with pagination
```
GET {{base_url}}/orders?page=1&limit=5
Authorization: Bearer {{token}}
```
**Expected: 200**

---

### GET /orders/:id — success
```
GET {{base_url}}/orders/1
Authorization: Bearer {{token}}
```
**Expected: 200** — order with full items detail

---

### GET /orders/:id — not found or belongs to another user
```
GET {{base_url}}/orders/99999
Authorization: Bearer {{token}}
```
**Expected: 404**
```json
{ "code": "NOT_FOUND", "message": "order not found" }
```

---

## 7. Static Files

### GET /uploads/... — access uploaded image
```
GET {{base_url}}/uploads/products/1_123456789.png
```
**Expected: 200** — the actual image binary

---

## End-to-End Flow

```
1.  POST /auth/register          → create account
2.  POST /auth/login             → get tokens
3.  POST /products               → create a product (as admin/user)
4.  POST /products/1/images      → upload product image
5.  PUT  /products/1/images/2/main → set main image
6.  GET  /products               → browse products (public)
7.  GET  /products?name=go&sort=price&order=asc → filtered browse
8.  POST /orders                 → place order (stock decrements)
9.  GET  /orders                 → view my orders
10. GET  /orders/1               → view order detail
11. POST /auth/refresh           → refresh expired token
12. PUT  /users/me               → update profile
13. PUT  /users/me/password      → change password
14. POST /auth/logout            → end session
```
