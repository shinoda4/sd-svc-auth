package service

import "errors"

var ErrInvalidPassword = errors.New("invalid password")
var ErrInvalidToken = errors.New("invalid token")
var ErrEmailNotVerified = errors.New("email not verified")
var ErrUsernameNotValid = errors.New("username not valid")
