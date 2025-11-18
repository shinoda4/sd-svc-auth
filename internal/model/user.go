package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               string    `db:"id"`
	Email            string    `db:"email"`
	Username         string    `db:"username"`
	PasswordHash     string    `db:"password_hash"`
	EmailVerified    bool      `db:"email_verified"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	ResetTokenExpire time.Time `db:"reset_token_expire"`
}

func (u *User) GetID() string       { return u.ID }
func (u *User) GetEmail() string    { return u.Email }
func (u *User) GetUsername() string { return u.Username }
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
func (u *User) GetEmailVerified() bool {
	return u.EmailVerified
}

func (u *User) GetResetTokenExpire() time.Time {
	return u.ResetTokenExpire
}
