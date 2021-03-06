package mocks

import "github.com/digitalocean/doctl/do"
import "github.com/stretchr/testify/mock"

// Generated: please do not edit by hand

type DriveActionsService struct {
	mock.Mock
}

func (_m *DriveActionsService) Attach(_a0 string, _a1 int) (*do.Action, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *do.Action
	if rf, ok := ret.Get(0).(func(string, int) *do.Action); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*do.Action)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *DriveActionsService) Detach(_a0 string) (*do.Action, error) {
	ret := _m.Called(_a0)

	var r0 *do.Action
	if rf, ok := ret.Get(0).(func(string) *do.Action); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*do.Action)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
