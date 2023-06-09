// Code generated by mockery v2.0.3. DO NOT EDIT.

package mocks

import (
	dto "github.com/diarikom/running-app/running-app-api/internal/api/dto"
	mock "github.com/stretchr/testify/mock"

	nhttp "github.com/diarikom/running-app/running-app-api/pkg/nhttp"
)

// AssetService is an autogenerated mock type for the AssetService type
type AssetService struct {
	mock.Mock
}

// GetPublicUrl provides a mock function with given fields: assetType, fileName
func (_m *AssetService) GetPublicUrl(assetType int, fileName string) string {
	ret := _m.Called(assetType, fileName)

	var r0 string
	if rf, ok := ret.Get(0).(func(int, string) string); ok {
		r0 = rf(assetType, fileName)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetUploadRule provides a mock function with given fields: assetType
func (_m *AssetService) GetUploadRule(assetType int) (*nhttp.UploadRule, error) {
	ret := _m.Called(assetType)

	var r0 *nhttp.UploadRule
	if rf, ok := ret.Get(0).(func(int) *nhttp.UploadRule); ok {
		r0 = rf(assetType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nhttp.UploadRule)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(assetType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UploadFile provides a mock function with given fields: req
func (_m *AssetService) UploadFile(req dto.UploadReq) (*dto.UploadResp, error) {
	ret := _m.Called(req)

	var r0 *dto.UploadResp
	if rf, ok := ret.Get(0).(func(dto.UploadReq) *dto.UploadResp); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.UploadResp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dto.UploadReq) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
