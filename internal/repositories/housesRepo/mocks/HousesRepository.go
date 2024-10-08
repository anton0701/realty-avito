// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	context "context"
	house "realty-avito/internal/repositories/housesRepo"

	mock "github.com/stretchr/testify/mock"
)

// HousesRepository is an autogenerated mock type for the HousesRepository type
type HousesRepository struct {
	mock.Mock
}

// CreateHouse provides a mock function with given fields: ctx, createHouseEntity
func (_m *HousesRepository) CreateHouse(ctx context.Context, createHouseEntity house.CreateHouseEntity) (*house.HouseEntity, error) {
	ret := _m.Called(ctx, createHouseEntity)

	var r0 *house.HouseEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, house.CreateHouseEntity) (*house.HouseEntity, error)); ok {
		return rf(ctx, createHouseEntity)
	}
	if rf, ok := ret.Get(0).(func(context.Context, house.CreateHouseEntity) *house.HouseEntity); ok {
		r0 = rf(ctx, createHouseEntity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*house.HouseEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, house.CreateHouseEntity) error); ok {
		r1 = rf(ctx, createHouseEntity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateHouseUpdatedAt provides a mock function with given fields: ctx, houseID
func (_m *HousesRepository) UpdateHouseUpdatedAt(ctx context.Context, houseID int64) error {
	ret := _m.Called(ctx, houseID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, houseID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewHousesRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewHousesRepository creates a new instance of HousesRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHousesRepository(t mockConstructorTestingTNewHousesRepository) *HousesRepository {
	mock := &HousesRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
