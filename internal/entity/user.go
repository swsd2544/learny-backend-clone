package entity

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	STUDENT = "student"
	TEACHER = "teacher"
)

const (
	SALT = 12
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	HashPassword []byte    `json:"-"`
	Coin         int64     `json:"coin"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	Version      int64     `json:"-"`
	CharacterID  int64     `json:"character_id"`
}

func (u *User) SetPassword(plaintextPassword string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), SALT)
	if err != nil {
		return err
	}
	u.HashPassword = hashPassword
	return nil
}

func (u *User) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.HashPassword, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
