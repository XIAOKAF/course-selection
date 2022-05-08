package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenClaims struct {
	Identify   string
	Variety    string
	ExpireTime time.Time
	jwt.StandardClaims
}
