package repo

import (
	"context"
	"errors"
	"time"

	"github.com/shinoda4/sd-svc-auth/internal/service/entity"
)

type MockUser struct {
	ID               string
	Username         string
	Email            string
	Password         string
	VerifyToken      string
	EmailVerified    bool
	ResetTokenExpire time.Time
}

func (u *MockUser) GetResetTokenExpire() time.Time {
	return u.ResetTokenExpire
}

func (u *MockUser) GetEmailVerified() bool {
	return u.EmailVerified
}

func (u *MockUser) GetID() string               { return u.ID }
func (u *MockUser) GetEmail() string            { return u.Email }
func (u *MockUser) GetUsername() string         { return u.Username }
func (u *MockUser) CheckPassword(p string) bool { return u.Password == p }

type MockUserRepo struct {
	users map[string]*MockUser
}

func (r *MockUserRepo) GetUserByVerifyToken(ctx context.Context, token string) (entity.UserEntity, error) {
	for _, u := range r.users {
		if u.VerifyToken == token {
			return u, nil
		}
	}
	return nil, errors.New("invalid verify token")
}

func (r *MockUserRepo) SetEmailVerified(ctx context.Context, userID string) error {
	for _, u := range r.users {
		if u.ID == userID {
			u.EmailVerified = true
			return nil
		}
	}
	return errors.New("user not found")
}

func (r *MockUserRepo) SetVerifyToken(ctx context.Context, userID, token string) error {
	for _, u := range r.users {
		if u.ID == userID {
			u.VerifyToken = token
			return nil
		}
	}
	return errors.New("user not found")
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{users: make(map[string]*MockUser)}
}

func (r *MockUserRepo) CreateUser(ctx context.Context, email, username, password string) (entity.UserEntity, error) {
	u := &MockUser{
		ID:       "mock-" + email,
		Email:    email,
		Username: username,
		Password: password,
	}

	r.users[email] = u
	return u, nil
}

func (r *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (entity.UserEntity, error) {
	u, ok := r.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

type MockRedis struct {
	data map[string]string
}

func (r *MockRedis) SetBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	key := "blacklist:" + token
	r.data[key] = "1" // 仅表示存在即可
	return nil
}

func (r *MockRedis) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := "refresh_token:" + userID
	delete(r.data, key)
	return nil
}

func NewMockRedis() *MockRedis {
	return &MockRedis{data: make(map[string]string)}
}

func (r *MockRedis) StoreToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	r.data[userID] = token
	return nil
}

func (r *MockRedis) GetToken(ctx context.Context, userID string) (string, error) {
	v, ok := r.data[userID]
	if !ok {
		return "", errors.New("token not found")
	}
	return v, nil
}
