// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/gpt/client.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/gpt/client.go -destination=internal/app/gpt/client_mock.go -package=gpt
//

// Package gpt is a generated GoMock package.
package gpt

import (
	reflect "reflect"

	resty "github.com/go-resty/resty/v2"
	gomock "go.uber.org/mock/gomock"
)

// MockHTTPClient is a mock of HTTPClient interface.
type MockHTTPClient struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPClientMockRecorder
	isgomock struct{}
}

// MockHTTPClientMockRecorder is the mock recorder for MockHTTPClient.
type MockHTTPClientMockRecorder struct {
	mock *MockHTTPClient
}

// NewMockHTTPClient creates a new mock instance.
func NewMockHTTPClient(ctrl *gomock.Controller) *MockHTTPClient {
	mock := &MockHTTPClient{ctrl: ctrl}
	mock.recorder = &MockHTTPClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHTTPClient) EXPECT() *MockHTTPClientMockRecorder {
	return m.recorder
}

// R mocks base method.
func (m *MockHTTPClient) R() *resty.Request {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "R")
	ret0, _ := ret[0].(*resty.Request)
	return ret0
}

// R indicates an expected call of R.
func (mr *MockHTTPClientMockRecorder) R() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "R", reflect.TypeOf((*MockHTTPClient)(nil).R))
}

// SetBaseURL mocks base method.
func (m *MockHTTPClient) SetBaseURL(url string) *resty.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBaseURL", url)
	ret0, _ := ret[0].(*resty.Client)
	return ret0
}

// SetBaseURL indicates an expected call of SetBaseURL.
func (mr *MockHTTPClientMockRecorder) SetBaseURL(url any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBaseURL", reflect.TypeOf((*MockHTTPClient)(nil).SetBaseURL), url)
}

// SetHeader mocks base method.
func (m *MockHTTPClient) SetHeader(header, value string) *resty.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", header, value)
	ret0, _ := ret[0].(*resty.Client)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockHTTPClientMockRecorder) SetHeader(header, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockHTTPClient)(nil).SetHeader), header, value)
}

// SetRetryCount mocks base method.
func (m *MockHTTPClient) SetRetryCount(count int) *resty.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRetryCount", count)
	ret0, _ := ret[0].(*resty.Client)
	return ret0
}

// SetRetryCount indicates an expected call of SetRetryCount.
func (mr *MockHTTPClientMockRecorder) SetRetryCount(count any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRetryCount", reflect.TypeOf((*MockHTTPClient)(nil).SetRetryCount), count)
}
