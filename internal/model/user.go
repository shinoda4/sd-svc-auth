/*
 * Copyright (c) 2025-11-20 shinoda4
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
