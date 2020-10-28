// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	requests "github.com/didil/kubexcloud/kxc-api/requests"
	mock "github.com/stretchr/testify/mock"

	responses "github.com/didil/kubexcloud/kxc-api/responses"
)

// AppSvc is an autogenerated mock type for the AppSvc type
type AppSvc struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, projectName, reqData
func (_m *AppSvc) Create(ctx context.Context, projectName string, reqData *requests.CreateApp) error {
	ret := _m.Called(ctx, projectName, reqData)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *requests.CreateApp) error); ok {
		r0 = rf(ctx, projectName, reqData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: ctx, projectName
func (_m *AppSvc) List(ctx context.Context, projectName string) (*responses.ListApp, error) {
	ret := _m.Called(ctx, projectName)

	var r0 *responses.ListApp
	if rf, ok := ret.Get(0).(func(context.Context, string) *responses.ListApp); ok {
		r0 = rf(ctx, projectName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*responses.ListApp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, projectName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Restart provides a mock function with given fields: ctx, projectName, appName
func (_m *AppSvc) Restart(ctx context.Context, projectName string, appName string) error {
	ret := _m.Called(ctx, projectName, appName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, projectName, appName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, projectName, appName, reqData
func (_m *AppSvc) Update(ctx context.Context, projectName string, appName string, reqData *requests.UpdateApp) error {
	ret := _m.Called(ctx, projectName, appName, reqData)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *requests.UpdateApp) error); ok {
		r0 = rf(ctx, projectName, appName, reqData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}