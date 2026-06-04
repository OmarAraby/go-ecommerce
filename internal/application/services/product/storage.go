package product

import (
	"context"
	"io"
)

// FileStorage is the contract the infrastructure layer must satisfy for file persistence.
// Defined here (at the consumer) — not at the implementation.
type FileStorage interface {
	Save(ctx context.Context, dir, filename string, r io.Reader) (url string, err error)
	Delete(ctx context.Context, url string) error
}
