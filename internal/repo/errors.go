package repo

import (
	"errors"
	"fmt"
)

type ErrUserExists struct {
	Email string
}

func (e *ErrUserExists) Error() string {
	return fmt.Sprintf("user already exists: %s", e.Email)
}

func NewErrUserExists(email string) error {
	return &ErrUserExists{Email: email}
}

var ErrNotFound = errors.New("not found")
