package nhttp

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"net/http"
)

const (
	KeyUserId    = "AUTH_USER_ID"
	KeySessionId = "AUTH_SESSION_ID"
)

// validateSessionFn is a function to validate session that has been extracted and make sure session is belongs to user
type ValidateSessionFn func(sessionId, userId string) (err error)

// validateUserTokenFn is a function to validate token and extract user id and session id
type ValidateUserTokenFn func(token string) (sessionId, userId string, err error)

// validateUserTokenFn is a contract function to validate token
type ValidateTokenFn func(token string) (err error)

// Middleware is a function that is able to chain between Handlers
type Middleware func(h Handler) Handler

// NewClientAuthMiddleware creates a middleware that validate a Client by its secret key
func NewClientAuthMiddleware(vFn ValidateTokenFn, authKey string, logger nlog.Logger) Middleware {
	// Return Middleware
	return func(next Handler) Handler {
		// Prepare function for app authentication handling
		fn := func(r *http.Request) (*Success, error) {
			// Get token
			authValue := r.Header.Get(authKey)

			// Validate token
			err := vFn(authValue)
			if err != nil {
				return nil, err
			}

			// Call next handler
			return next.Fn(r)
		}

		return Handler{Fn: fn, Logger: logger}
	}
}

// NewUserAuthMiddleware creates a middleware that validate user access token before calling handler function
func NewUserAuthMiddleware(vFn ValidateUserTokenFn, sFn ValidateSessionFn, authKey string,
	logger nlog.Logger) Middleware {
	// Return Middleware
	return func(next Handler) Handler {
		// Prepare function for  user auth handling
		fn := func(r *http.Request) (*Success, error) {
			// Get token
			authValue := r.Header.Get(authKey)

			// Validate token
			sessionId, userId, err := vFn(authValue)
			if err != nil {
				return nil, err
			}

			// Validate session
			err = sFn(sessionId, userId)
			if err != nil {
				return nil, err
			}

			// Set user id, session id to header
			r.Header.Set(KeyUserId, userId)
			r.Header.Set(KeySessionId, sessionId)

			// Call next handler
			return next.Fn(r)
		}

		return Handler{Fn: fn, Logger: logger}
	}
}
