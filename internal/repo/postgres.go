package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shinoda4/sd-svc-auth/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type UserRepo struct {
	db *sqlx.DB
}

func NewPostgres(dsn string) (*UserRepo, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &UserRepo{db: db}, nil
}

func (r *UserRepo) Close() {
	err := r.db.Close()
	if err != nil {
		return
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, email, username, password string) error {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email)
	if err != nil {
		return fmt.Errorf("check user existence: %w", err)
	}
	if exists {
		return NewErrUserExists(email)
	}

	// hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)`, email, username, string(hash))

	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (service.UserEntity, error) {
	u := &User{}
	err := r.db.GetContext(ctx, u, `SELECT id, email, password_hash FROM users WHERE email=$1`, email)
	if err != nil {
		return nil, ErrNotFound
	}
	return u, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
