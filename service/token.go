package service

import (
	"course-selection/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("www..xyz.com")

func RememberStatus(id string, duration time.Duration) (error, string) {
	return CreateToken(id, duration)
}

func CreateToken(id string, duration time.Duration) (error, string) {
	expireTime := time.Now().Add(duration * time.Minute)
	claims := model.TokenClaims{
		Identify:   id,
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
	var tokenClaims model.TokenClaims
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*model.TokenClaims)
	if !ok {
		return nil, errors.New("fail to parse token")
	}
	err = token.Claims.Valid()
	if err != nil {
		return nil, err
	}
	return claims, nil
}
