// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"
	db "scanner/pkg/db"

	mock "github.com/stretchr/testify/mock"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"

	time "time"
)

// ResultStore is an autogenerated mock type for the ResultStore type
type ResultStore struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, v
func (_m *ResultStore) Add(ctx context.Context, v *db.Result) (primitive.ObjectID, error) {
	ret := _m.Called(ctx, v)

	var r0 primitive.ObjectID
	if rf, ok := ret.Get(0).(func(context.Context, *db.Result) primitive.ObjectID); ok {
		r0 = rf(ctx, v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(primitive.ObjectID)
		}
	}

	var r1 error
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
	if rf, ok := ret.Get(0).(func(context.Context, *db.FilterResult) []db.Result); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]db.Result)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, *db.FilterResult) int64); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *db.FilterResult) error); ok {
		r2 = rf(ctx, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UpdateFinding provides a mock function with given fields: ctx, id, findings
func (_m *ResultStore) UpdateFinding(ctx context.Context, id primitive.ObjectID, findings string) error {
	ret := _m.Called(ctx, id, findings)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, string) error); ok {
		r0 = rf(ctx, id, findings)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateFinishedAt provides a mock function with given fields: ctx, id, t
func (_m *ResultStore) UpdateFinishedAt(ctx context.Context, id primitive.ObjectID, t time.Time) error {
	ret := _m.Called(ctx, id, t)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, time.Time) error); ok {
		r0 = rf(ctx, id, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateQueuedAt provides a mock function with given fields: ctx, id, t
func (_m *ResultStore) UpdateQueuedAt(ctx context.Context, id primitive.ObjectID, t time.Time) error {
	ret := _m.Called(ctx, id, t)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, time.Time) error); ok {
		r0 = rf(ctx, id, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateScanningAt provides a mock function with given fields: ctx, id, t
func (_m *ResultStore) UpdateScanningAt(ctx context.Context, id primitive.ObjectID, t time.Time) error {
	ret := _m.Called(ctx, id, t)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, time.Time) error); ok {
		r0 = rf(ctx, id, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateStatus provides a mock function with given fields: ctx, id, status
func (_m *ResultStore) UpdateStatus(ctx context.Context, id primitive.ObjectID, status db.ResultStatus) error {
	ret := _m.Called(ctx, id, status)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, db.ResultStatus) error); ok {
		r0 = rf(ctx, id, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
