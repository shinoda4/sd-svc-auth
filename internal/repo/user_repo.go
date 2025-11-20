package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/shinoda4/sd-svc-auth/internal/model"
	"github.com/shinoda4/sd-svc-auth/internal/service/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	Repo
}

func NewUserRepo(dsn string) (*UserRepo, error) {
	repo, err := NewPostgres(dsn)
	if err != nil {
		return nil, err
	}
	return &UserRepo{Repo: *repo}, nil
}

func (r *UserRepo) GetUserByVerifyToken(ctx context.Context, token string) (entity.UserEntity, error) {
	u := &model.User{}
	err := r.db.GetContext(ctx, u,
		`SELECT id, email, username, password_hash FROM users WHERE verify_token=$1`, token)
	if err != nil {
		return nil, fmt.Errorf("query user failed: %w", err)
	}
	return u, nil
}
func (r *UserRepo) SetEmailVerified(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET email_verified=true, verify_token=NULL WHERE id=$1`, userID)
	return err
}

func (r *UserRepo) CreateUser(ctx context.Context, email, username, password string) (entity.UserEntity, error) {
	var exists bool

	err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email)
	if err != nil {
		return nil, fmt.Errorf("check user existence: %w", err)
	}
	if exists {
		return nil, NewErrUserExists(email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	var id string
	err = r.db.GetContext(ctx, &id,
		`INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		email, username, string(hash))
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	user := &model.User{
		ID:           id,
		Email:        email,
		Username:     username,
		PasswordHash: string(hash),
	}
	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (entity.UserEntity, error) {
	u := &model.User{}
	err := r.db.GetContext(ctx, u, `SELECT id, email, username, password_hash, email_verified FROM users WHERE email=$1`, email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) SetVerifyToken(ctx context.Context, userID, token string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET verify_token=$1 WHERE id=$2`, token, userID)
	return err
}

func (r *UserRepo) SaveResetToken(ctx context.Context, userID, token string, expire time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET reset_token=$1, reset_token_expire=$2 WHERE id=$3`,
		token, expire, userID,
	)
	return err
}

func (r *UserRepo) GetUserByResetToken(ctx context.Context, token string) (entity.UserEntity, error) {
	u := &model.User{}
	err := r.db.GetContext(ctx, u,
		`SELECT id, email, reset_token_expire FROM users WHERE reset_token=$1`,
		token,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash=$1 WHERE id=$2`,
		hashed, userID,
	)
	return err
}

func (r *UserRepo) ClearResetToken(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET reset_token=NULL, reset_token_expire=NULL WHERE id=$1`,
		userID,
	)
	return err
}
