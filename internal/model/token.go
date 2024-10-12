package model

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	UserID uint
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}
