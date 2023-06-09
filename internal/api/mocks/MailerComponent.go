// Code generated by mockery v2.0.0-alpha.2. DO NOT EDIT.

package mocks

import (
	nmailgun "github.com/diarikom/running-app/running-app-api/pkg/nmailgun"
	mock "github.com/stretchr/testify/mock"
)

// MailerComponent is an autogenerated mock type for the MailerComponent type
type MailerComponent struct {
	mock.Mock
}

// GetDefaultSender provides a mock function with given fields:
func (_m *MailerComponent) GetDefaultSender() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetSender provides a mock function with given fields: name
func (_m *MailerComponent) GetSender(name string) string {
	ret := _m.Called(name)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Send provides a mock function with given fields: opt
func (_m *MailerComponent) Send(opt nmailgun.SendOpt) error {
	ret := _m.Called(opt)

	var r0 error
	if rf, ok := ret.Get(0).(func(nmailgun.SendOpt) error); ok {
		r0 = rf(opt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
