package entity

import (
	"errors"
	"github.com/swsd2544/learny-backend-clone/internal/validator"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	RoleStudent = "student"
	RoleTeacher = "teacher"
)

const (
	SALT = 12
)

type password struct {
	plaintext *string
	Hash      []byte
}

type User struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Email       string    `json:"email"`
	Password    password  `json:"-"`
	Coin        int64     `json:"coin"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	Version     int64     `json:"-"`
	CharacterID int64     `json:"character_id"`
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), SALT)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.Hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
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

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "name", "must be provided")
	v.Check(len(user.Username) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(user.Firstname != "", "name", "must be provided")
	v.Check(len(user.Firstname) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(user.Lastname != "", "name", "must be provided")
	v.Check(len(user.Lastname) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(user.Coin >= 0, "coin", "must be a positive number")
	v.Check(validator.PermittedValue(user.Role, RoleStudent, RoleTeacher), "role", "must be either student or teacher")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePassword(v, *user.Password.plaintext)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}

func ValidatePassword(v *validator.Validator, plaintextPassword string) {
	const entropy = 60
	err := passwordvalidator.Validate(plaintextPassword, entropy)
	if err != nil {
		v.AddError("password", err.Error())
	}
}
