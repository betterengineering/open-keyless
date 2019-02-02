// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/scanner/scanner.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockScanner is a mock of Scanner interface
type MockScanner struct {
	ctrl     *gomock.Controller
	recorder *MockScannerMockRecorder
}

// MockScannerMockRecorder is the mock recorder for MockScanner
type MockScannerMockRecorder struct {
	mock *MockScanner
}

// NewMockScanner creates a new mock instance
func NewMockScanner(ctrl *gomock.Controller) *MockScanner {
	mock := &MockScanner{ctrl: ctrl}
	mock.recorder = &MockScannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockScanner) EXPECT() *MockScannerMockRecorder {
	return m.recorder
}

// Scan mocks base method
func (m *MockScanner) Scan() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Scan")
}

// Scan indicates an expected call of Scan
func (mr *MockScannerMockRecorder) Scan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockScanner)(nil).Scan))
}

// Done mocks base method
func (m *MockScanner) Done() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Done")
	ret0, _ := ret[0].(error)
	return ret0
}

// Done indicates an expected call of Done
func (mr *MockScannerMockRecorder) Done() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Done", reflect.TypeOf((*MockScanner)(nil).Done))
}
