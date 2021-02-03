package types

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// Session represents a user session
type Session struct {
	ID     string `json:"session_id"`
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// SignToken creates a JWT token from the session
func SignToken(signingKey string, claims *Session) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(signingKey))
	if err != nil {
		log.Errorf("Cannot sign token: %s", err)
		return "", err
	}

	return signed, nil
}

// ParseToken parses a token string into a session
func ParseToken(signingKey, tokenStr string) (*Session, error) {
	session := &Session{}
	_, err := jwt.ParseWithClaims(tokenStr, session, func(token *jwt.Token) (interface{}, error) {
		switch token.Method.(type) {
		case *jwt.SigningMethodHMAC:
			return []byte(signingKey), nil
		default:
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	})
	if err != nil {
		return nil, err
	}
	return session, nil
}
