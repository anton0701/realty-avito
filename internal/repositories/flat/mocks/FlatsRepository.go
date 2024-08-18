// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	context "context"
	flat "realty-avito/internal/repositories/flat"

	mock "github.com/stretchr/testify/mock"
)

// FlatsRepository is an autogenerated mock type for the FlatsRepository type
type FlatsRepository struct {
	mock.Mock
}

// CreateFlat provides a mock function with given fields: ctx, flatModel
func (_m *FlatsRepository) CreateFlat(ctx context.Context, flatModel flat.CreateFlatEntity) (*flat.FlatEntity, error) {
	ret := _m.Called(ctx, flatModel)

	var r0 *flat.FlatEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flat.CreateFlatEntity) (*flat.FlatEntity, error)); ok {
		return rf(ctx, flatModel)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flat.CreateFlatEntity) *flat.FlatEntity); ok {
		r0 = rf(ctx, flatModel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flat.FlatEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flat.CreateFlatEntity) error); ok {
		r1 = rf(ctx, flatModel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetApprovedFlatsByHouseID provides a mock function with given fields: ctx, houseID
func (_m *FlatsRepository) GetApprovedFlatsByHouseID(ctx context.Context, houseID int64) ([]flat.FlatEntity, error) {
	ret := _m.Called(ctx, houseID)

	var r0 []flat.FlatEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]flat.FlatEntity, error)); ok {
		return rf(ctx, houseID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []flat.FlatEntity); ok {
		r0 = rf(ctx, houseID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flat.FlatEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, houseID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFlatByFlatID provides a mock function with given fields: ctx, flatID
func (_m *FlatsRepository) GetFlatByFlatID(ctx context.Context, flatID int64) (*flat.FlatEntity, error) {
	ret := _m.Called(ctx, flatID)

	var r0 *flat.FlatEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*flat.FlatEntity, error)); ok {
		return rf(ctx, flatID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *flat.FlatEntity); ok {
		r0 = rf(ctx, flatID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flat.FlatEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, flatID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFlatsByHouseID provides a mock function with given fields: ctx, houseID
func (_m *FlatsRepository) GetFlatsByHouseID(ctx context.Context, houseID int64) ([]flat.FlatEntity, error) {
	ret := _m.Called(ctx, houseID)

	var r0 []flat.FlatEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]flat.FlatEntity, error)); ok {
		return rf(ctx, houseID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []flat.FlatEntity); ok {
		r0 = rf(ctx, houseID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flat.FlatEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, houseID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateFlat provides a mock function with given fields: ctx, updateFlatModel
func (_m *FlatsRepository) UpdateFlat(ctx context.Context, updateFlatModel flat.UpdateFlatEntity) (*flat.FlatEntity, error) {
	ret := _m.Called(ctx, updateFlatModel)

	var r0 *flat.FlatEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flat.UpdateFlatEntity) (*flat.FlatEntity, error)); ok {
		return rf(ctx, updateFlatModel)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flat.UpdateFlatEntity) *flat.FlatEntity); ok {
		r0 = rf(ctx, updateFlatModel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flat.FlatEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flat.UpdateFlatEntity) error); ok {
		r1 = rf(ctx, updateFlatModel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFlatsRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewFlatsRepository creates a new instance of FlatsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFlatsRepository(t mockConstructorTestingTNewFlatsRepository) *FlatsRepository {
	mock := &FlatsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
