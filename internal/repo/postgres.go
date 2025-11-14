package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shinoda4/sd-svc-auth/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, dsn string) (*UserRepo, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &UserRepo{pool: pool}, nil
}

func (r *UserRepo) Close() {
	r.pool.Close()
}

func (r *UserRepo) CreateUser(ctx context.Context, email, password string) error {

	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check user existence: %w", err)
	}
	if exists {
		return NewErrUserExists(email)
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2)`, email, string(hash))

	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (service.UserEntity, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, email, password_hash FROM users WHERE email=$1`, email)
	u := &User{}
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash); err != nil {
		return nil, ErrNotFound
	}
	return u, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
