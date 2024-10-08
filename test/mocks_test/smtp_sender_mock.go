// Code generated by MockGen. DO NOT EDIT.
// Source: ./src/notification/infrastructure/sender/smtp_client.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type MockSmtpClient struct {
	ctrl     *gomock.Controller
	recorder *MockSmtpClientMockRecorder
}

type MockSmtpClientMockRecorder struct {
	mock *MockSmtpClient
}

func NewMockSmtpClient(ctrl *gomock.Controller) *MockSmtpClient {
	mock := &MockSmtpClient{ctrl: ctrl}
	mock.recorder = &MockSmtpClientMockRecorder{mock}
	return mock
}

func (m *MockSmtpClient) EXPECT() *MockSmtpClientMockRecorder {
	return m.recorder
}

func (m *MockSmtpClient) Send(email, text string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", email, text)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockSmtpClientMockRecorder) Send(email, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSmtpClient)(nil).Send), email, text)
}
