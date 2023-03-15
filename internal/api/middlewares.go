package api

import (
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"net/http"
)

type ValidateResetPasswordTokenFn func(token string) (*dto.ResetPasswordSession, error)
type ValidateResetPasswordUserFn func(session *dto.ResetPasswordSession) (string, error)

// / NewResetPasswordSessionMiddleware creates a middleware that validate a one time token for
// / password reset before calling handler function
func NewResetPasswordSessionMiddleware(vFn ValidateResetPasswordTokenFn, uFn ValidateResetPasswordUserFn, authKey string,
	logger nlog.Logger) nhttp.Middleware {
	// Return Middleware
	return func(next nhttp.Handler) nhttp.Handler {
		// Prepare function for  user auth handling
		fn := func(r *http.Request) (*nhttp.Success, error) {
			// Get token
			authValue := r.Header.Get(authKey)

			// Validate token and get session claims
			session, err := vFn(authValue)
			if err != nil {
				return nil, err
			}

			// Validate user signature and get user id
			userId, err := uFn(session)
			if err != nil {
				return nil, err
			}

			// Set user id, session id to header
			r.Header.Set(nhttp.KeyUserId, userId)

			// Call next handler
			return next.Fn(r)
		}

		return nhttp.Handler{Fn: fn, Logger: logger}
	}
}

type ValidateVerifyEmailTokenFn func(token string) (*dto.VerifyEmailSession, error)
type ValidateVerifyEmailUserFn func(session *dto.VerifyEmailSession) (string, error)

// / NewVerifyEmailSessionMiddleware creates a middleware that validate a one time token for
// / email verification before calling handler function
func NewVerifyEmailSessionMiddleware(vFn ValidateVerifyEmailTokenFn, uFn ValidateVerifyEmailUserFn, authKey string,
	logger nlog.Logger) nhttp.Middleware {
	// Return Middleware
	return func(next nhttp.Handler) nhttp.Handler {
		// Prepare function for  user auth handling
		fn := func(r *http.Request) (*nhttp.Success, error) {
			// Get token
			authValue := r.Header.Get(authKey)

			// Validate token and get session claims
			session, err := vFn(authValue)
			if err != nil {
				return nil, err
			}

			// Validate user signature and get user id
			userId, err := uFn(session)
			if err != nil {
				return nil, err
			}

			// Set user id, session id to header
			r.Header.Set(nhttp.KeyUserId, userId)

			// Call next handler
			return next.Fn(r)
		}

		return nhttp.Handler{Fn: fn, Logger: logger}
	}
}
