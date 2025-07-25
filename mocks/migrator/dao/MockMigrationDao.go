// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	"github.com/wizards-0/go-pins/migrator/dao"
	"github.com/wizards-0/go-pins/migrator/types"
)

var TYPE_MIGRATION_LOG = mock.AnythingOfType("types.MigrationLog")
var TYPE_MIGRATION = mock.AnythingOfType("types.Migration")

var passThroughMap = map[string]func(mockDao *MockMigrationDao){
	"SetupMigrationTable": func(mockDao *MockMigrationDao) {
		mockDao.EXPECT().SetupMigrationTable().RunAndReturn(func() error {
			return mockDao.realDao.SetupMigrationTable()
		}).Once()
	},
	"GetMigrationLogs": func(mockDao *MockMigrationDao) {
		mockDao.EXPECT().GetMigrationLogs().RunAndReturn(func() ([]types.MigrationLog, error) {
			return mockDao.realDao.GetMigrationLogs()
		}).Once()
	},
	"InsertMigrationLog": func(mockDao *MockMigrationDao) {
		mockDao.EXPECT().InsertMigrationLog(TYPE_MIGRATION_LOG).RunAndReturn(func(mLog types.MigrationLog) error {
			return mockDao.realDao.InsertMigrationLog(mLog)
		}).Once()
	},
	"UpdateMigrationStatus": func(mockDao *MockMigrationDao) {
		mockDao.EXPECT().UpdateMigrationStatus(TYPE_MIGRATION_LOG).RunAndReturn(func(mLog types.MigrationLog) error {
			return mockDao.realDao.UpdateMigrationStatus(mLog)
		}).Once()
	},
	"DeleteMigrationLog": func(mockDao *MockMigrationDao) {
		mockDao.EXPECT().DeleteMigrationLog(TYPE_MIGRATION_LOG).RunAndReturn(func(mLog types.MigrationLog) error {
			return mockDao.realDao.DeleteMigrationLog(mLog)
		}).Once()
	},
	"ExecuteQuery": func(mockDao *MockMigrationDao) {
		mockDao.EXPECT().ExecuteQuery(TYPE_MIGRATION).RunAndReturn(func(m types.Migration) error {
			return mockDao.realDao.ExecuteQuery(m)
		}).Once()
	},
	"ExecuteRollback": func(mockDao *MockMigrationDao) {
		mockDao.EXPECT().ExecuteRollback(TYPE_MIGRATION).RunAndReturn(func(m types.Migration) error {
			return mockDao.realDao.ExecuteRollback(m)
		}).Once()
	},
}

// NewMockMigrationDao creates a new instance of MockMigrationDao. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMigrationDao(orig dao.MigrationDao, t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMigrationDao {
	mock := &MockMigrationDao{
		realDao: orig,
	}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockMigrationDao is an autogenerated mock type for the MigrationDao type
type MockMigrationDao struct {
	mock.Mock
	realDao dao.MigrationDao
}

type MockMigrationDao_Expecter struct {
	mock *mock.Mock
}

func (mockDao *MockMigrationDao) PassThrough(methodNames ...string) {
	for _, name := range methodNames {
		fn, exists := passThroughMap[name]
		if exists {
			fn(mockDao)
		}
	}
}

func (_m *MockMigrationDao) EXPECT() *MockMigrationDao_Expecter {
	return &MockMigrationDao_Expecter{mock: &_m.Mock}
}

// DeleteMigrationLog provides a mock function for the type MockMigrationDao
func (_mock *MockMigrationDao) DeleteMigrationLog(mLog types.MigrationLog) error {
	ret := _mock.Called(mLog)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMigrationLog")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(types.MigrationLog) error); ok {
		r0 = returnFunc(mLog)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockMigrationDao_DeleteMigrationLog_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteMigrationLog'
type MockMigrationDao_DeleteMigrationLog_Call struct {
	*mock.Call
}

// DeleteMigrationLog is a helper method to define mock.On call
//   - mLog types.MigrationLog
func (_e *MockMigrationDao_Expecter) DeleteMigrationLog(mLog interface{}) *MockMigrationDao_DeleteMigrationLog_Call {
	return &MockMigrationDao_DeleteMigrationLog_Call{Call: _e.mock.On("DeleteMigrationLog", mLog)}
}

func (_c *MockMigrationDao_DeleteMigrationLog_Call) Run(run func(mLog types.MigrationLog)) *MockMigrationDao_DeleteMigrationLog_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 types.MigrationLog
		if args[0] != nil {
			arg0 = args[0].(types.MigrationLog)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockMigrationDao_DeleteMigrationLog_Call) Return(err error) *MockMigrationDao_DeleteMigrationLog_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockMigrationDao_DeleteMigrationLog_Call) RunAndReturn(run func(mLog types.MigrationLog) error) *MockMigrationDao_DeleteMigrationLog_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteQuery provides a mock function for the type MockMigrationDao
func (_mock *MockMigrationDao) ExecuteQuery(m types.Migration) error {
	ret := _mock.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteQuery")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(types.Migration) error); ok {
		r0 = returnFunc(m)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockMigrationDao_ExecuteQuery_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteQuery'
type MockMigrationDao_ExecuteQuery_Call struct {
	*mock.Call
}

// ExecuteQuery is a helper method to define mock.On call
//   - m types.Migration
func (_e *MockMigrationDao_Expecter) ExecuteQuery(m interface{}) *MockMigrationDao_ExecuteQuery_Call {
	return &MockMigrationDao_ExecuteQuery_Call{Call: _e.mock.On("ExecuteQuery", m)}
}

func (_c *MockMigrationDao_ExecuteQuery_Call) Run(run func(m types.Migration)) *MockMigrationDao_ExecuteQuery_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 types.Migration
		if args[0] != nil {
			arg0 = args[0].(types.Migration)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockMigrationDao_ExecuteQuery_Call) Return(err error) *MockMigrationDao_ExecuteQuery_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockMigrationDao_ExecuteQuery_Call) RunAndReturn(run func(m types.Migration) error) *MockMigrationDao_ExecuteQuery_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteRollback provides a mock function for the type MockMigrationDao
func (_mock *MockMigrationDao) ExecuteRollback(m types.Migration) error {
	ret := _mock.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteRollback")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(types.Migration) error); ok {
		r0 = returnFunc(m)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockMigrationDao_ExecuteRollback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteRollback'
type MockMigrationDao_ExecuteRollback_Call struct {
	*mock.Call
}

// ExecuteRollback is a helper method to define mock.On call
//   - m types.Migration
func (_e *MockMigrationDao_Expecter) ExecuteRollback(m interface{}) *MockMigrationDao_ExecuteRollback_Call {
	return &MockMigrationDao_ExecuteRollback_Call{Call: _e.mock.On("ExecuteRollback", m)}
}

func (_c *MockMigrationDao_ExecuteRollback_Call) Run(run func(m types.Migration)) *MockMigrationDao_ExecuteRollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 types.Migration
		if args[0] != nil {
			arg0 = args[0].(types.Migration)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockMigrationDao_ExecuteRollback_Call) Return(err error) *MockMigrationDao_ExecuteRollback_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockMigrationDao_ExecuteRollback_Call) RunAndReturn(run func(m types.Migration) error) *MockMigrationDao_ExecuteRollback_Call {
	_c.Call.Return(run)
	return _c
}

// GetMigrationLogs provides a mock function for the type MockMigrationDao
func (_mock *MockMigrationDao) GetMigrationLogs() ([]types.MigrationLog, error) {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMigrationLogs")
	}

	var r0 []types.MigrationLog
	var r1 error
	if returnFunc, ok := ret.Get(0).(func() ([]types.MigrationLog, error)); ok {
		return returnFunc()
	}
	if returnFunc, ok := ret.Get(0).(func() []types.MigrationLog); ok {
		r0 = returnFunc()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.MigrationLog)
		}
	}
	if returnFunc, ok := ret.Get(1).(func() error); ok {
		r1 = returnFunc()
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockMigrationDao_GetMigrationLogs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMigrationLogs'
type MockMigrationDao_GetMigrationLogs_Call struct {
	*mock.Call
}

// GetMigrationLogs is a helper method to define mock.On call
func (_e *MockMigrationDao_Expecter) GetMigrationLogs() *MockMigrationDao_GetMigrationLogs_Call {
	return &MockMigrationDao_GetMigrationLogs_Call{Call: _e.mock.On("GetMigrationLogs")}
}

func (_c *MockMigrationDao_GetMigrationLogs_Call) Run(run func()) *MockMigrationDao_GetMigrationLogs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMigrationDao_GetMigrationLogs_Call) Return(migrationLogs []types.MigrationLog, err error) *MockMigrationDao_GetMigrationLogs_Call {
	_c.Call.Return(migrationLogs, err)
	return _c
}

func (_c *MockMigrationDao_GetMigrationLogs_Call) RunAndReturn(run func() ([]types.MigrationLog, error)) *MockMigrationDao_GetMigrationLogs_Call {
	_c.Call.Return(run)
	return _c
}

// InsertMigrationLog provides a mock function for the type MockMigrationDao
func (_mock *MockMigrationDao) InsertMigrationLog(mLog types.MigrationLog) error {
	ret := _mock.Called(mLog)

	if len(ret) == 0 {
		panic("no return value specified for InsertMigrationLog")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(types.MigrationLog) error); ok {
		r0 = returnFunc(mLog)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockMigrationDao_InsertMigrationLog_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertMigrationLog'
type MockMigrationDao_InsertMigrationLog_Call struct {
	*mock.Call
}

// InsertMigrationLog is a helper method to define mock.On call
//   - mLog types.MigrationLog
func (_e *MockMigrationDao_Expecter) InsertMigrationLog(mLog interface{}) *MockMigrationDao_InsertMigrationLog_Call {
	return &MockMigrationDao_InsertMigrationLog_Call{Call: _e.mock.On("InsertMigrationLog", mLog)}
}

func (_c *MockMigrationDao_InsertMigrationLog_Call) Run(run func(mLog types.MigrationLog)) *MockMigrationDao_InsertMigrationLog_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 types.MigrationLog
		if args[0] != nil {
			arg0 = args[0].(types.MigrationLog)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockMigrationDao_InsertMigrationLog_Call) Return(err error) *MockMigrationDao_InsertMigrationLog_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockMigrationDao_InsertMigrationLog_Call) RunAndReturn(run func(mLog types.MigrationLog) error) *MockMigrationDao_InsertMigrationLog_Call {
	_c.Call.Return(run)
	return _c
}

// SetupMigrationTable provides a mock function for the type MockMigrationDao
func (_mock *MockMigrationDao) SetupMigrationTable() error {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for SetupMigrationTable")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func() error); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockMigrationDao_SetupMigrationTable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetupMigrationTable'
type MockMigrationDao_SetupMigrationTable_Call struct {
	*mock.Call
}

// SetupMigrationTable is a helper method to define mock.On call
func (_e *MockMigrationDao_Expecter) SetupMigrationTable() *MockMigrationDao_SetupMigrationTable_Call {
	return &MockMigrationDao_SetupMigrationTable_Call{Call: _e.mock.On("SetupMigrationTable")}
}

func (_c *MockMigrationDao_SetupMigrationTable_Call) Run(run func()) *MockMigrationDao_SetupMigrationTable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMigrationDao_SetupMigrationTable_Call) Return(err error) *MockMigrationDao_SetupMigrationTable_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockMigrationDao_SetupMigrationTable_Call) RunAndReturn(run func() error) *MockMigrationDao_SetupMigrationTable_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMigrationStatus provides a mock function for the type MockMigrationDao
func (_mock *MockMigrationDao) UpdateMigrationStatus(mLog types.MigrationLog) error {
	ret := _mock.Called(mLog)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMigrationStatus")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(types.MigrationLog) error); ok {
		r0 = returnFunc(mLog)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockMigrationDao_UpdateMigrationStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMigrationStatus'
type MockMigrationDao_UpdateMigrationStatus_Call struct {
	*mock.Call
}

// UpdateMigrationStatus is a helper method to define mock.On call
//   - mLog types.MigrationLog
func (_e *MockMigrationDao_Expecter) UpdateMigrationStatus(mLog interface{}) *MockMigrationDao_UpdateMigrationStatus_Call {
	return &MockMigrationDao_UpdateMigrationStatus_Call{Call: _e.mock.On("UpdateMigrationStatus", mLog)}
}

func (_c *MockMigrationDao_UpdateMigrationStatus_Call) Run(run func(mLog types.MigrationLog)) *MockMigrationDao_UpdateMigrationStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 types.MigrationLog
		if args[0] != nil {
			arg0 = args[0].(types.MigrationLog)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockMigrationDao_UpdateMigrationStatus_Call) Return(err error) *MockMigrationDao_UpdateMigrationStatus_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockMigrationDao_UpdateMigrationStatus_Call) RunAndReturn(run func(mLog types.MigrationLog) error) *MockMigrationDao_UpdateMigrationStatus_Call {
	_c.Call.Return(run)
	return _c
}
