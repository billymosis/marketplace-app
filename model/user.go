package model

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

const (
	devSaltRounds  = 8
	prodSaltRounds = 10
)

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
