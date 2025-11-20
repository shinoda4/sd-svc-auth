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

package email

import (
	"errors"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(from string, to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	if emailPassword == "" {
		err := errors.New("EMAIL_PASSWORD environment variable not set")
		return err
	}
	emailAddress := os.Getenv("EMAIL_ADDRESS")
	if emailAddress == "" {
		err := errors.New("EMAIL_ADDRESS environment variable not set")
		return err
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, emailAddress, emailPassword)

	return d.DialAndSend(m)
}
