package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	userapp "github.com/OmarAraby/go-ecommerce/internal/application/services/user"
	"github.com/OmarAraby/go-ecommerce/internal/domain"
	"github.com/OmarAraby/go-ecommerce/internal/domain/entities"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/db"
)

var _ userapp.UserRepository = (*UserRepo)(nil)
var _ userapp.RefreshTokenRepository = (*RefreshTokenRepo)(nil)
var _ userapp.PasswordResetRepository = (*PasswordResetRepo)(nil)

// ── User ──────────────────────────────────────────────────────────────────────

type UserRepo struct{ q *db.Queries }

func NewUserRepo(pool *pgxpool.Pool) *UserRepo { return &UserRepo{q: db.New(pool)} }

func (r *UserRepo) Create(ctx context.Context, u *entities.User) (*entities.User, error) {
	row, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Name: u.Name, Email: u.Email, Password: u.Password, Role: u.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}
	return userToDomain(row), nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*entities.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}
	return userToDomain(row), nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("GetUserByEmail: %w", err)
	}
	return userToDomain(row), nil
}

func (r *UserRepo) UpdateProfile(ctx context.Context, id int64, name string) (*entities.User, error) {
	row, err := r.q.UpdateUserProfile(ctx, db.UpdateUserProfileParams{ID: id, Name: name})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("UpdateProfile: %w", err)
	}
	return userToDomain(row), nil
}

func (r *UserRepo) UpdateEmail(ctx context.Context, id int64, email string) (*entities.User, error) {
	row, err := r.q.UpdateUserEmail(ctx, db.UpdateUserEmailParams{ID: id, Email: email})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("UpdateEmail: %w", err)
	}
	return userToDomain(row), nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	return r.q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{ID: id, Password: hashedPassword})
}

func userToDomain(u db.User) *entities.User {
	return &entities.User{
		ID: u.ID, Name: u.Name, Email: u.Email, Password: u.Password, Role: u.Role,
		CreatedAt: u.CreatedAt.Time, UpdatedAt: u.UpdatedAt.Time,
	}
}

// ── Refresh Token ─────────────────────────────────────────────────────────────

type RefreshTokenRepo struct{ q *db.Queries }

func NewRefreshTokenRepo(pool *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{q: db.New(pool)}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	_, err := r.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID: userID, Token: token, ExpiresAt: pgxTimestamp(expiresAt),
	})
	return err
}

func (r *RefreshTokenRepo) Get(ctx context.Context, token string) (*entities.RefreshToken, error) {
	row, err := r.q.GetRefreshToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("GetRefreshToken: %w", err)
	}
	return &entities.RefreshToken{
		ID: row.ID, UserID: row.UserID, Token: row.Token,
		ExpiresAt: row.ExpiresAt.Time, Revoked: row.Revoked, CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *RefreshTokenRepo) Revoke(ctx context.Context, token string) error {
	return r.q.RevokeRefreshToken(ctx, token)
}

func (r *RefreshTokenRepo) RevokeAll(ctx context.Context, userID int64) error {
	return r.q.RevokeAllUserRefreshTokens(ctx, userID)
}

// ── Password Reset ────────────────────────────────────────────────────────────

type PasswordResetRepo struct{ q *db.Queries }

func NewPasswordResetRepo(pool *pgxpool.Pool) *PasswordResetRepo {
	return &PasswordResetRepo{q: db.New(pool)}
}

func (r *PasswordResetRepo) Create(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	_, err := r.q.CreatePasswordResetToken(ctx, db.CreatePasswordResetTokenParams{
		UserID: userID, Token: token, ExpiresAt: pgxTimestamp(expiresAt),
	})
	return err
}

func (r *PasswordResetRepo) Get(ctx context.Context, token string) (*entities.PasswordResetToken, error) {
	row, err := r.q.GetPasswordResetToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("GetPasswordResetToken: %w", err)
	}
	return &entities.PasswordResetToken{
		ID: row.ID, UserID: row.UserID, Token: row.Token,
		ExpiresAt: row.ExpiresAt.Time, Used: row.Used, CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *PasswordResetRepo) MarkUsed(ctx context.Context, id int64) error {
	return r.q.MarkPasswordResetTokenUsed(ctx, id)
}
