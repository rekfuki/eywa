package auth

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	// OperatorUserID ...
	OperatorUserID = "00000000-0000-0000-0000-000000000000"
)

// Auth holds the request authorisation details from request header
type Auth struct {
	UserID     string
	RealUserID string
	UserAgent  string
}

// FromHeaders constructs an Auth object from a given set of headers
func FromHeaders(h http.Header) *Auth {
	return &Auth{
		UserID:     h.Get("X-User-Id"),
		RealUserID: h.Get("X-Real-User-Id"),
		UserAgent:  h.Get("User-Agent"),
	}
}

// IsOperator checks if user id header is equal to operator
func (a *Auth) IsOperator() bool {
	return a.RealUserID == OperatorUserID
}

// Check checks auth
func (a *Auth) Check() bool {
	if a.UserID == "" {
		log.Warnf("X-User-Id header not sent, ua: %s", a.UserAgent)
		return false
	}

	if _, err := uuid.FromString(a.UserID); err != nil {
		log.Warnf("Invalid X-User-Id header was sent: %s", a.UserID)
		return false
	}

	if a.RealUserID == "" {
		log.Warnf("X-Real-User-Id header not sent, ua: %s", a.UserAgent)
		a.RealUserID = a.UserID
	}

	if _, err := uuid.FromString(a.RealUserID); err != nil {
		log.Warnf("Invalid X-Real-User-Id header was sent: %s", a.RealUserID)
		return false
	}

	return true
}
