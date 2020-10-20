// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ServiceRequester is an autogenerated mock type for the ServiceRequester type
type ServiceRequester struct {
	mock.Mock
}

// Request provides a mock function with given fields: _a0
func (_m *ServiceRequester) Request(_a0 context.Context) (int, error) {
	ret := _m.Called(_a0)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
