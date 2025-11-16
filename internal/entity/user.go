package entity

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID            string
	Email         string
	Username      string
	PasswordHash  string
	EmailVerified bool
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
