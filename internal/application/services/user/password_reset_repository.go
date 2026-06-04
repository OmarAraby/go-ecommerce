package user

import (
	"context"
	"time"

	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	Get(ctx context.Context, token string) (*entities.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id int64) error
}
