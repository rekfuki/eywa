package types

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// Token represents a system token that is generated
// when issuing access tokens or creating sessions
type Token struct {
	ID             string `json:"session_id"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	UserEmail      string `json:"user_email"`
	UserAvatar     string `json:"user_avatar"`
	CompletionStep string `json:"completion_step"`
	jwt.StandardClaims
}

// SignToken creates a JWT token from the session
func SignToken(signingKey string, t *Token) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)

	signed, err := token.SignedString([]byte(signingKey))
	if err != nil {
		log.Errorf("Cannot sign token: %s", err)
		return "", err
	}

	return signed, nil
}

// ParseToken parses a token string into a decoded token
func ParseToken(signingKey, tokenStr string) (*Token, error) {
	session := &Token{}
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
