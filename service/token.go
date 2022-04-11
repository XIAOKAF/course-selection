package service

import (
	"course-selection/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func CreateToken(id string, duration time.Duration) (error, string) {
	expireTime := time.Now().Add(duration * time.Minute)
	claims := model.TokenClaims{
		UserId:     id,
		ExpireTime: expireTime,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(token)
	if err != nil {
		return err, tokenString
	}
	return nil, tokenString
}