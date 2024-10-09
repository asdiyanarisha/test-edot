// Code generated by mockery v2.46.2. DO NOT EDIT.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	models "test-edot/src/models"
)

// StockLevelRepositoryInterface is an autogenerated mock type for the StockLevelRepositoryInterface type
type StockLevelRepositoryInterface struct {
	mock.Mock
}

// Begin provides a mock function with given fields:
func (_m *StockLevelRepositoryInterface) Begin() *gorm.DB {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Begin")
	}

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// Create provides a mock function with given fields: tx, stockLevel
func (_m *StockLevelRepositoryInterface) Create(tx *gorm.DB, stockLevel *models.StockLevel) error {
	ret := _m.Called(tx, stockLevel)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *models.StockLevel) error); ok {
		r0 = rf(tx, stockLevel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOne provides a mock function with given fields: ctx, selectField, query, args
func (_m *StockLevelRepositoryInterface) FindOne(ctx context.Context, selectField string, query string, args ...any) (models.StockLevel, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, selectField, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 models.StockLevel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, ...any) (models.StockLevel, error)); ok {
		return rf(ctx, selectField, query, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, ...any) models.StockLevel); ok {
		r0 = rf(ctx, selectField, query, args...)
	} else {
		r0 = ret.Get(0).(models.StockLevel)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, ...any) error); ok {
		r1 = rf(ctx, selectField, query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneTx provides a mock function with given fields: tx, order, query, args
func (_m *StockLevelRepositoryInterface) FindOneTx(tx *gorm.DB, order string, query string, args ...interface{}) (models.StockLevelProduct, error) {
	var _ca []interface{}
	_ca = append(_ca, tx, order, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindOneTx")
	}

	var r0 models.StockLevelProduct
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string, ...interface{}) (models.StockLevelProduct, error)); ok {
		return rf(tx, order, query, args...)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string, ...interface{}) models.StockLevelProduct); ok {
		r0 = rf(tx, order, query, args...)
	} else {
		r0 = ret.Get(0).(models.StockLevelProduct)
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, string, string, ...interface{}) error); ok {
		r1 = rf(tx, order, query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTx provides a mock function with given fields: tx, order, query, args
func (_m *StockLevelRepositoryInterface) FindTx(tx *gorm.DB, order string, query string, args ...interface{}) ([]models.StockLevelProduct, error) {
	var _ca []interface{}
	_ca = append(_ca, tx, order, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindTx")
	}

	var r0 []models.StockLevelProduct
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string, ...interface{}) ([]models.StockLevelProduct, error)); ok {
		return rf(tx, order, query, args...)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string, ...interface{}) []models.StockLevelProduct); ok {
		r0 = rf(tx, order, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.StockLevelProduct)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, string, string, ...interface{}) error); ok {
		r1 = rf(tx, order, query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SumStockWarehouse provides a mock function with given fields: ctx, query, args
func (_m *StockLevelRepositoryInterface) SumStockWarehouse(ctx context.Context, query string, args ...any) (models.StockWarehouse, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for SumStockWarehouse")
	}

	var r0 models.StockWarehouse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...any) (models.StockWarehouse, error)); ok {
		return rf(ctx, query, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...any) models.StockWarehouse); ok {
		r0 = rf(ctx, query, args...)
	} else {
		r0 = ret.Get(0).(models.StockWarehouse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...any) error); ok {
		r1 = rf(ctx, query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOneTx provides a mock function with given fields: tx, updateStockLevel, selectFields, query, args
func (_m *StockLevelRepositoryInterface) UpdateOneTx(tx *gorm.DB, updateStockLevel *models.StockLevel, selectFields string, query string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, tx, updateStockLevel, selectFields, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOneTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *models.StockLevel, string, string, ...interface{}) error); ok {
		r0 = rf(tx, updateStockLevel, selectFields, query, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStockLevelRepositoryInterface creates a new instance of StockLevelRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStockLevelRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *StockLevelRepositoryInterface {
	mock := &StockLevelRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
