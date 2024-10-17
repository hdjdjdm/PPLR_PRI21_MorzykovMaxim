package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var csrfSecret = []byte("ugga-bugga-hihihihahaha")

type CSRFClaims struct {
	SessionID string `json:"session_id"`
	jwt.StandardClaims
}

func GenerateCSRFToken(sessionID string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &CSRFClaims{
		SessionID: sessionID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(csrfSecret)
}

func ValidateCSRFToken(tokenString string, r *http.Request) error {
	token, err := jwt.ParseWithClaims(tokenString, &CSRFClaims{}, func(token *jwt.Token) (interface{}, error) {
		return csrfSecret, nil
	})

	if claims, ok := token.Claims.(*CSRFClaims); ok && token.Valid {
		sessionID, err := r.Cookie("session_id")
		if err != nil || claims.SessionID != sessionID.Value {
			return errors.New("invalid CSRF token: session mismatch")
		}
		return nil
	} else {
		return err
	}
}
