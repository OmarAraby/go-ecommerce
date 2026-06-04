package user

import (
	"context"

	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, u *entities.User) (*entities.User, error)
	GetByID(ctx context.Context, id int64) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	UpdateProfile(ctx context.Context, id int64, name string) (*entities.User, error)
	UpdateEmail(ctx context.Context, id int64, email string) (*entities.User, error)
	UpdatePassword(ctx context.Context, id int64, hashedPassword string) error
}
