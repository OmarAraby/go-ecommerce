package user

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/OmarAraby/go-ecommerce/internal/domain"
	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/auth"
)

var _ Service = (*service)(nil)

type service struct {
	users       UserRepository
	refreshRepo RefreshTokenRepository
	resetRepo   PasswordResetRepository
	jwtSecret   string
}

func NewService(
	users UserRepository,
	refreshRepo RefreshTokenRepository,
	resetRepo PasswordResetRepository,
	jwtSecret string,
) Service {
	return &service{users: users, refreshRepo: refreshRepo, resetRepo: resetRepo, jwtSecret: jwtSecret}
}

func (s *service) Register(ctx context.Context, name, email, password string) (*entities.User, error) {
	if _, err := s.users.GetByEmail(ctx, email); err == nil {
		return nil, fmt.Errorf("email already registered: %w", domain.ErrConflict)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	return s.users.Create(ctx, &entities.User{Name: name, Email: email, Password: string(hash), Role: "user"})
}

func (s *service) Login(ctx context.Context, email, password string) (*AuthResultDTO, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	return s.issueTokenPair(ctx, u)
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*TokenPairDTO, error) {
	rec, err := s.refreshRepo.Get(ctx, refreshToken)
	if err != nil || rec.Revoked || time.Now().After(rec.ExpiresAt) {
		return nil, domain.ErrInvalidCredentials
	}
	u, err := s.users.GetByID(ctx, rec.UserID)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	_ = s.refreshRepo.Revoke(ctx, refreshToken)
	accessToken, err := auth.GenerateAccessToken(u.ID, u.Email, u.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}
	newRefresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	if err := s.refreshRepo.Create(ctx, u.ID, newRefresh, time.Now().Add(auth.RefreshTokenDuration)); err != nil {
		return nil, err
	}
	return &TokenPairDTO{AccessToken: accessToken, RefreshToken: newRefresh}, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	return s.refreshRepo.Revoke(ctx, refreshToken)
}

func (s *service) ForgotPassword(ctx context.Context, email string) (string, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", nil
	}
	token, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", err
	}
	if err := s.resetRepo.Create(ctx, u.ID, token, time.Now().Add(1*time.Hour)); err != nil {
		return "", err
	}
	return token, nil
}

func (s *service) ResetPassword(ctx context.Context, token, newPassword string) error {
	rec, err := s.resetRepo.Get(ctx, token)
	if err != nil || rec.Used || time.Now().After(rec.ExpiresAt) {
		return domain.ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	if err := s.users.UpdatePassword(ctx, rec.UserID, string(hash)); err != nil {
		return err
	}
	_ = s.resetRepo.MarkUsed(ctx, rec.ID)
	_ = s.refreshRepo.RevokeAll(ctx, rec.UserID)
	return nil
}

func (s *service) GetProfile(ctx context.Context, userID int64) (*entities.User, error) {
	return s.users.GetByID(ctx, userID)
}

func (s *service) UpdateProfile(ctx context.Context, userID int64, name string) (*entities.User, error) {
	return s.users.UpdateProfile(ctx, userID, name)
}

func (s *service) ChangeEmail(ctx context.Context, userID int64, newEmail, currentPassword string) (*entities.User, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(currentPassword)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	return s.users.UpdateEmail(ctx, userID, newEmail)
}

func (s *service) ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return domain.ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(currentPassword)); err != nil {
		return domain.ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	return s.users.UpdatePassword(ctx, userID, string(hash))
}

func (s *service) issueTokenPair(ctx context.Context, u *entities.User) (*AuthResultDTO, error) {
	accessToken, err := auth.GenerateAccessToken(u.ID, u.Email, u.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	if err := s.refreshRepo.Create(ctx, u.ID, refreshToken, time.Now().Add(auth.RefreshTokenDuration)); err != nil {
		return nil, err
	}
	return &AuthResultDTO{AccessToken: accessToken, RefreshToken: refreshToken, User: u}, nil
}
