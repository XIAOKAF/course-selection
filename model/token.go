package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenClaims struct {
	UserId     string
	ExpireTime time.Time
	jwt.StandardClaims
}
