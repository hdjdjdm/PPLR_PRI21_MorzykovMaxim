package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Login string `json:"login"`
	Role  string `json:"role"`
	jwt.StandardClaims
}
