// Code generated by mockery v2.0.0-alpha.2. DO NOT EDIT.

package mocks

import (
	njwt "github.com/diarikom/running-app/running-app-api/internal/pkg/njwt"
	mock "github.com/stretchr/testify/mock"
)

// JWTIssuerComponent is an autogenerated mock type for the JWTIssuerComponent type
type JWTIssuerComponent struct {
	mock.Mock
}

// New provides a mock function with given fields: opt
func (_m *JWTIssuerComponent) New(opt njwt.ClaimOpt) (*njwt.Token, error) {
	ret := _m.Called(opt)

	var r0 *njwt.Token
	if rf, ok := ret.Get(0).(func(njwt.ClaimOpt) *njwt.Token); ok {
		r0 = rf(opt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*njwt.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(njwt.ClaimOpt) error); ok {
		r1 = rf(opt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Verify provides a mock function with given fields: input
func (_m *JWTIssuerComponent) Verify(input string) (*njwt.Claim, error) {
	ret := _m.Called(input)

	var r0 *njwt.Claim
	if rf, ok := ret.Get(0).(func(string) *njwt.Claim); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*njwt.Claim)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
