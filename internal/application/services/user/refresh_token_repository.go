package user

import (
	"context"
	"time"

	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	Get(ctx context.Context, token string) (*entities.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAll(ctx context.Context, userID int64) error
}
