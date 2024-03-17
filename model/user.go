package model

import (
	"context"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

const (
	devSaltRounds  = 8
	prodSaltRounds = 10
)

type UserStore interface {
	GetById(ctx context.Context, id uint) (*User, error)
	GetByUsername(ctx context.Context, name string) (*User, error)
	CreateUser(context.Context, *User) (*User, error)
	GetValidator() *validator.Validate
}

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (user *User) HashPassword() error {
	saltRound, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), saltRound)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return nil
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
