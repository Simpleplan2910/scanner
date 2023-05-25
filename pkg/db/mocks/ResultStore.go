// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	context "context"
	db "scanner/pkg/db"

	mock "github.com/stretchr/testify/mock"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// ResultStore is an autogenerated mock type for the ResultStore type
type ResultStore struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, v
func (_m *ResultStore) Add(ctx context.Context, v *db.Result) (primitive.ObjectID, error) {
	ret := _m.Called(ctx, v)

	var r0 primitive.ObjectID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *db.Result) (primitive.ObjectID, error)); ok {
		return rf(ctx, v)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *db.Result) primitive.ObjectID); ok {
		r0 = rf(ctx, v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(primitive.ObjectID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *db.Result) error); ok {
		r1 = rf(ctx, v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Filter provides a mock function with given fields: ctx, filter
func (_m *ResultStore) Filter(ctx context.Context, filter *db.FilterResult) ([]db.Result, int64, error) {
	ret := _m.Called(ctx, filter)

	var r0 []db.Result
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *db.FilterResult) ([]db.Result, int64, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *db.FilterResult) []db.Result); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]db.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *db.FilterResult) int64); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, *db.FilterResult) error); ok {
		r2 = rf(ctx, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewResultStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewResultStore creates a new instance of ResultStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewResultStore(t mockConstructorTestingTNewResultStore) *ResultStore {
	mock := &ResultStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
