package user

import (
	"context"

	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
)

type Service interface {
	// Auth
	Register(ctx context.Context, name, email, password string) (*entities.User, error)
	Login(ctx context.Context, email, password string) (*AuthResultDTO, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenPairDTO, error)
	Logout(ctx context.Context, refreshToken string) error

	// Password
	ForgotPassword(ctx context.Context, email string) (resetToken string, err error)
	ResetPassword(ctx context.Context, token, newPassword string) error

	// Profile
	GetProfile(ctx context.Context, userID int64) (*entities.User, error)
	UpdateProfile(ctx context.Context, userID int64, name string) (*entities.User, error)
	ChangeEmail(ctx context.Context, userID int64, newEmail, currentPassword string) (*entities.User, error)
	ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error
}
