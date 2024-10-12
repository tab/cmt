package git

import (
	"context"
	"os/exec"
	"reflect"

	"go.uber.org/mock/gomock"
)

type MockExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockExecutorMockRecorder
}

type MockExecutorMockRecorder struct {
	mock *MockExecutor
}

func NewMockExecutor(ctrl *gomock.Controller) *MockExecutor {
	mock := &MockExecutor{ctrl: ctrl}
	mock.recorder = &MockExecutorMockRecorder{mock}
	return mock
}

func (m *MockExecutor) EXPECT() *MockExecutorMockRecorder {
	return m.recorder
}

func (m *MockExecutor) Run(ctx context.Context, name string, arg ...string) *exec.Cmd {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, name}
	for _, a := range arg {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Run", varargs...)
	ret0, _ := ret[0].(*exec.Cmd)
	return ret0
}

func (mr *MockExecutorMockRecorder) Run(ctx, name interface{}, arg ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, name}, arg...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockExecutor)(nil).Run), varargs...)
}
