package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/internal/api/entity"
	"github.com/diarikom/running-app/running-app-api/internal/pkg/njwt"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"strings"
	"time"
)

type Authenticator struct {
	TokenIssuer           api.JWTIssuerComponent
	IdGen                 *api.SnowflakeGen
	Errors                *api.Errors
	Logger                nlog.Logger
	AppClientSecret       string
	DashboardClientSecret string
}

func (a *Authenticator) Init(app *api.Api) error {
	a.TokenIssuer = app.Components.JWTIssuer
	a.IdGen = app.Components.Id
	a.Errors = app.Components.Errors
	a.Logger = app.Logger
	a.AppClientSecret = app.Config.GetString(api.ConfAppClientSecret)
	a.DashboardClientSecret = app.Config.GetString(api.ConfDashboardClientSecret)
	return nil
}

func (a *Authenticator) ValidateVerifyEmailToken(bearer string) (*dto.VerifyEmailSession, error) {
	// Extract bearer token
	token, err := a.ExtractBearerToken(bearer)
	if err != nil {
		return nil, err
	}

	// Verify token
	claim, err := a.TokenIssuer.Verify(token)
	if err != nil {
		// Convert token error and return
		return nil, a.GetTokenError(err)
	}

	// Verify Purpose
	if claim.Purpose != api.JWTPurposeVerifyEmail {
		return nil, nhttp.ErrUnauthorized
	}

	resp := dto.VerifyEmailSession{
		RequestId:     claim.Session,
		Email:         claim.Subject,
		UserSignature: claim.Extra[api.UserSignatureKey],
	}
	return &resp, nil
}

func (a *Authenticator) ValidateResetPasswordToken(bearer string) (*dto.ResetPasswordSession, error) {
	// Extract bearer token
	token, err := a.ExtractBearerToken(bearer)
	if err != nil {
		return nil, err
	}

	// Verify token
	claim, err := a.TokenIssuer.Verify(token)
	if err != nil {
		// Convert token error and return
		return nil, a.GetTokenError(err)
	}

	// Validate purpose
	if claim.Purpose != api.JWTPurposeResetPassword {
		return nil, nhttp.ErrUnauthorized
	}

	resp := dto.ResetPasswordSession{
		RequestId:     claim.Session,
		Email:         claim.Subject,
		UserSignature: claim.Extra[api.UserSignatureKey],
	}
	return &resp, nil
}

func (a *Authenticator) SignMd5(req dto.SignatureReq) (string, error) {
	raw := fmt.Sprintf(req.Format, req.Args...)
	hasher := md5.New()
	_, err := hasher.Write([]byte(raw))
	if err != nil {
		a.Logger.Error("unable to sign md5", err)
		return "", err
	}
	signature := hex.EncodeToString(hasher.Sum(nil))
	return signature, nil
}

func (a *Authenticator) NewOneTimeToken(req dto.JWTOptReq) (*entity.AccessToken, error) {
	// Validate purpose
	switch req.Purpose {
	case api.JWTUser, api.JWTApp:
		return nil, errors.New("invalid purpose")
	}

	// Create token
	t, err := a.TokenIssuer.New(njwt.ClaimOpt{
		SessionId: req.SessionId,
		Subject:   req.Subject,
		Audience:  api.JWTAudienceApp,
		Lifetime:  time.Duration(req.Lifetime),
		Purpose:   req.Purpose,
		Extras:    req.Extras,
	})

	if err != nil {
		a.Logger.Error("unable to issue user one time token", err)
		return nil, err
	}

	return &entity.AccessToken{
		Token:     t.Encoded,
		ExpiredAt: t.ExpiredAt,
	}, err
}

func (a *Authenticator) GetTokenError(err error) error {
	// If error is jwt.ValidationError
	if e, ok := err.(*jwt.ValidationError); ok {
		// Converts to nhttp.Error
		switch e.Errors {
		case jwt.ValidationErrorExpired:
			return a.Errors.New("USR003")
		case jwt.ValidationErrorMalformed:
			return a.Errors.New("USR004")
		default:
			a.Logger.Error("unhandled jwt issuer error", err)
		}
	}

	// Else, return error
	return err
}

func (a *Authenticator) ExtractBearerToken(authString string) (result string, err error) {
	if authString == "" {
		a.Logger.Errorf("Authorization Header is empty")
		err = nhttp.ErrUnauthorized
		return
	}
	// Extract token
	splitToken := strings.Split(authString, " ")
	if len(splitToken) != 2 {
		a.Logger.Errorf("Bearer token is malformed")
		err = nhttp.ErrUnauthorized
		return
	}
	// Return token
	return splitToken[1], nil
}

func (a *Authenticator) ValidateUserAccess(bearer string) (sessionId string, userId string, err error) {
	// Extract bearer token
	token, err := a.ExtractBearerToken(bearer)
	if err != nil {
		return sessionId, userId, err
	}

	// Verify token
	claim, err := a.TokenIssuer.Verify(token)
	if err != nil {
		// Convert token error and return
		return sessionId, userId, a.GetTokenError(err)
	}

	// Verify purpose
	if claim.Purpose != api.JWTUser {
		return "", "", nhttp.ErrUnauthorized
	}

	// Get user id and subject
	sessionId = claim.Session
	userId = claim.Subject

	return sessionId, userId, nil
}

func (a *Authenticator) NewAccessToken(req dto.JWTOptReq) (*entity.AccessToken, error) {
	t, err := a.TokenIssuer.New(njwt.ClaimOpt{
		SessionId: req.SessionId,
		Subject:   req.Subject,
		Audience:  api.JWTAudienceUser,
		Lifetime:  time.Duration(req.Lifetime),
		Purpose:   api.JWTUser,
	})

	if err != nil {
		a.Logger.Error("unable to issue user access token", err)
		return nil, err
	}

	return &entity.AccessToken{
		Token:     t.Encoded,
		ExpiredAt: t.ExpiredAt,
	}, err
}

func (a *Authenticator) ValidateClient(secret string) (err error) {
	if a.AppClientSecret == secret {
		return nil
	}

	return nhttp.ErrUnauthorized
}

func (a *Authenticator) ValidateClientDashboard(secret string) (err error) {
	if a.DashboardClientSecret == secret {
		return nil
	}

	return nhttp.ErrUnauthorized
}
