// Code generated by mockery v2.0.3. DO NOT EDIT.

package mocks

import (
	dto "github.com/diarikom/running-app/running-app-api/internal/api/dto"
	mock "github.com/stretchr/testify/mock"
)

// AdTagService is an autogenerated mock type for the AdTagService type
type AdTagService struct {
	mock.Mock
}

// GetAdTags provides a mock function with given fields:
func (_m *AdTagService) GetAdTags() (*dto.AdTagResp, error) {
	ret := _m.Called()

	var r0 *dto.AdTagResp
	if rf, ok := ret.Get(0).(func() *dto.AdTagResp); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.AdTagResp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
