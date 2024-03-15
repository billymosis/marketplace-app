package auth

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type jwtCustomClaims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(id uint, userName string) (string, error) {
	now := time.Now()
	var expiration time.Time
	environment := os.Getenv("ENVIRONMENT")
	if environment == "production" {
		expiration = now.Add(time.Minute * 2)
	} else {
		expiration = now.Add(time.Hour * 1)
	}
	claims := &jwtCustomClaims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return t, nil
}

func GetUserId(ctx context.Context) (uint, error) {
	props, _ := ctx.Value("userAuthCtx").(jwt.MapClaims)
	logrus.Printf("%+v\n", props)

	userId, err := strconv.ParseInt(fmt.Sprintf("%v", props["user_id"]), 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(userId), nil
}
