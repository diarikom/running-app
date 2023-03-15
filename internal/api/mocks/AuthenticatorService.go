// Code generated by mockery v2.0.3. DO NOT EDIT.

package mocks

import (
	dto "github.com/diarikom/running-app/running-app-api/internal/api/dto"
	entity "github.com/diarikom/running-app/running-app-api/internal/api/entity"

	mock "github.com/stretchr/testify/mock"
)

// AuthenticatorService is an autogenerated mock type for the AuthenticatorService type
type AuthenticatorService struct {
	mock.Mock
}

// NewAccessToken provides a mock function with given fields: req
func (_m *AuthenticatorService) NewAccessToken(req dto.JWTOptReq) (*entity.AccessToken, error) {
	ret := _m.Called(req)

	var r0 *entity.AccessToken
	if rf, ok := ret.Get(0).(func(dto.JWTOptReq) *entity.AccessToken); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.AccessToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dto.JWTOptReq) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOneTimeToken provides a mock function with given fields: req
func (_m *AuthenticatorService) NewOneTimeToken(req dto.JWTOptReq) (*entity.AccessToken, error) {
	ret := _m.Called(req)

	var r0 *entity.AccessToken
	if rf, ok := ret.Get(0).(func(dto.JWTOptReq) *entity.AccessToken); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.AccessToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dto.JWTOptReq) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignMd5 provides a mock function with given fields: req
func (_m *AuthenticatorService) SignMd5(req dto.SignatureReq) (string, error) {
	ret := _m.Called(req)

	var r0 string
	if rf, ok := ret.Get(0).(func(dto.SignatureReq) string); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dto.SignatureReq) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateClient provides a mock function with given fields: secret
func (_m *AuthenticatorService) ValidateClient(secret string) error {
	ret := _m.Called(secret)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(secret)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateClientDashboard provides a mock function with given fields: secret
func (_m *AuthenticatorService) ValidateClientDashboard(secret string) error {
	ret := _m.Called(secret)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(secret)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateResetPasswordToken provides a mock function with given fields: token
func (_m *AuthenticatorService) ValidateResetPasswordToken(token string) (*dto.ResetPasswordSession, error) {
	ret := _m.Called(token)

	var r0 *dto.ResetPasswordSession
	if rf, ok := ret.Get(0).(func(string) *dto.ResetPasswordSession); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.ResetPasswordSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateUserAccess provides a mock function with given fields: bearer
func (_m *AuthenticatorService) ValidateUserAccess(bearer string) (string, string, error) {
	ret := _m.Called(bearer)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(bearer)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(bearer)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(bearer)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ValidateVerifyEmailToken provides a mock function with given fields: token
func (_m *AuthenticatorService) ValidateVerifyEmailToken(token string) (*dto.VerifyEmailSession, error) {
	ret := _m.Called(token)

	var r0 *dto.VerifyEmailSession
	if rf, ok := ret.Get(0).(func(string) *dto.VerifyEmailSession); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.VerifyEmailSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
