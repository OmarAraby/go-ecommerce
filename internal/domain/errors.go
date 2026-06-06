package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")            // e.g. duplicate email
	ErrInvalidCredentials = errors.New("invalid credentials") // wrong password / bad token
	ErrInvalidFileType    = errors.New("invalid file type")    // unsupported image format
	ErrImageLimitReached  = errors.New("image limit reached")   // max 6 images per product
	ErrInsufficientStock  = errors.New("insufficient stock")    // not enough product stock
)
