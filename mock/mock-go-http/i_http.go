// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pefish/go-http (interfaces: IHttp)

// Package mock_go_http is a generated GoMock package.
package mock_go_http

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	go_http "github.com/pefish/go-http"
)

// MockIHttp is a mock of IHttp interface.
type MockIHttp struct {
	ctrl     *gomock.Controller
	recorder *MockIHttpMockRecorder
}

// MockIHttpMockRecorder is the mock recorder for MockIHttp.
type MockIHttpMockRecorder struct {
	mock *MockIHttp
}

// NewMockIHttp creates a new mock instance.
func NewMockIHttp(ctrl *gomock.Controller) *MockIHttp {
	mock := &MockIHttp{ctrl: ctrl}
	mock.recorder = &MockIHttpMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIHttp) EXPECT() *MockIHttpMockRecorder {
	return m.recorder
}

// GetForBytes mocks base method.
func (m *MockIHttp) GetForBytes(arg0 *go_http.RequestParams) (*http.Response, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForBytes", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetForBytes indicates an expected call of GetForBytes.
func (mr *MockIHttpMockRecorder) GetForBytes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForBytes", reflect.TypeOf((*MockIHttp)(nil).GetForBytes), arg0)
}

// GetForString mocks base method.
func (m *MockIHttp) GetForString(arg0 *go_http.RequestParams) (*http.Response, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForString", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetForString indicates an expected call of GetForString.
func (mr *MockIHttpMockRecorder) GetForString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForString", reflect.TypeOf((*MockIHttp)(nil).GetForString), arg0)
}

// GetForStruct mocks base method.
func (m *MockIHttp) GetForStruct(arg0 *go_http.RequestParams, arg1 interface{}) (*http.Response, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForStruct", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetForStruct indicates an expected call of GetForStruct.
func (mr *MockIHttpMockRecorder) GetForStruct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForStruct", reflect.TypeOf((*MockIHttp)(nil).GetForStruct), arg0, arg1)
}

// PostForBytes mocks base method.
func (m *MockIHttp) PostForBytes(arg0 *go_http.RequestParams) (*http.Response, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostForBytes", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostForBytes indicates an expected call of PostForBytes.
func (mr *MockIHttpMockRecorder) PostForBytes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostForBytes", reflect.TypeOf((*MockIHttp)(nil).PostForBytes), arg0)
}

// PostForString mocks base method.
func (m *MockIHttp) PostForString(arg0 *go_http.RequestParams) (*http.Response, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostForString", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostForString indicates an expected call of PostForString.
func (mr *MockIHttpMockRecorder) PostForString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostForString", reflect.TypeOf((*MockIHttp)(nil).PostForString), arg0)
}

// PostForStruct mocks base method.
func (m *MockIHttp) PostForStruct(arg0 *go_http.RequestParams, arg1 interface{}) (*http.Response, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostForStruct", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostForStruct indicates an expected call of PostForStruct.
func (mr *MockIHttpMockRecorder) PostForStruct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostForStruct", reflect.TypeOf((*MockIHttp)(nil).PostForStruct), arg0, arg1)
}

// PostFormDataForStruct mocks base method.
func (m *MockIHttp) PostFormDataForStruct(arg0 *go_http.RequestParams, arg1 interface{}) (*http.Response, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostFormDataForStruct", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostFormDataForStruct indicates an expected call of PostFormDataForStruct.
func (mr *MockIHttpMockRecorder) PostFormDataForStruct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostFormDataForStruct", reflect.TypeOf((*MockIHttp)(nil).PostFormDataForStruct), arg0, arg1)
}

// PostMultipart mocks base method.
func (m *MockIHttp) PostMultipart(arg0 *go_http.PostMultipartParams) (*http.Response, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostMultipart", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostMultipart indicates an expected call of PostMultipart.
func (mr *MockIHttpMockRecorder) PostMultipart(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostMultipart", reflect.TypeOf((*MockIHttp)(nil).PostMultipart), arg0)
}

// PostMultipartForStruct mocks base method.
func (m *MockIHttp) PostMultipartForStruct(arg0 *go_http.PostMultipartParams, arg1 interface{}) (*http.Response, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostMultipartForStruct", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostMultipartForStruct indicates an expected call of PostMultipartForStruct.
func (mr *MockIHttpMockRecorder) PostMultipartForStruct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostMultipartForStruct", reflect.TypeOf((*MockIHttp)(nil).PostMultipartForStruct), arg0, arg1)
}
