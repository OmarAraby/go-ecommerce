package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")           // e.g. duplicate email
	ErrInvalidCredentials = errors.New("invalid credentials") // wrong password / bad token
)
