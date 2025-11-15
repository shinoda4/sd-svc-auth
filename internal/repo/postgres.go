package repo

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/shinoda4/sd-svc-auth/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Email        string
	Username     string
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

func (r *UserRepo) CreateUser(ctx context.Context, email, username, password string) (service.UserEntity, error) {
	var exists bool
	// 先检查用户是否已存在
	err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email)
	if err != nil {
		return nil, fmt.Errorf("check user existence: %w", err)
	}
	if exists {
		return nil, NewErrUserExists(email)
	}

	// hash 密码
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// 插入用户，同时返回 id
	var id string
	err = r.db.GetContext(ctx, &id,
		`INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		email, username, string(hash))
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	// 构造并返回 User 对象
	user := &User{
		ID:           id,
		Email:        email,
		Username:     username,
		PasswordHash: string(hash),
	}
	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (service.UserEntity, error) {
	u := &User{}
	err := r.db.GetContext(ctx, u, `SELECT id, email, username, password_hash FROM users WHERE email=$1`, email)
	if err != nil {
		return nil, ErrNotFound
	}
	return u, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
func (r *UserRepo) SetVerifyToken(ctx context.Context, userID, token string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET verify_token=$1 WHERE id=$2`, token, userID)
	return err
}
