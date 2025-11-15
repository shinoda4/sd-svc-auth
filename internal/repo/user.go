package repo

import (
	"context"
	"fmt"

	"github.com/shinoda4/sd-svc-auth/internal/service"
)

func (u *User) GetID() string       { return u.ID }
func (u *User) GetEmail() string    { return u.Email }
func (u *User) GetUsername() string { return u.Username }

func (r *UserRepo) GetUserByVerifyToken(ctx context.Context, token string) (service.UserEntity, error) {
	u := &User{}
	err := r.db.GetContext(ctx, u,
		`SELECT id, email, username FROM users WHERE verify_token=$1`, token)
	if err != nil {
		return nil, fmt.Errorf("query user failed: %w", err)
	}
	return u, nil
}
func (r *UserRepo) SetEmailVerified(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET email_verified=true, verify_token=NULL WHERE id=$1`, userID)
	return err
}
