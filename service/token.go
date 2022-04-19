package service

import (
	"course-selection/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("www..xyz.com")

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
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return err, tokenString
	}
	return nil, tokenString
}

func ParseToken(tokenString string) (*model.TokenClaims, error) {
	tokenClaims := &model.TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return tokenClaims, err
	}
	claims, ok := token.Claims.(*model.TokenClaims)
	if !ok {
		return claims, errors.New("token解析失败")
	}
	return claims, nil
}
