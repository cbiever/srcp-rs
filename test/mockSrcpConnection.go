package test

import (
	"testing"
)

type SendAndReceive func(request string) string

type MockSrcpConnection struct {
	t              *testing.T
	sendAndReceive SendAndReceive
}

func (mockSrcpConnection *MockSrcpConnection) SetSendAndReceive(sendAndReceive SendAndReceive) {
	mockSrcpConnection.sendAndReceive = sendAndReceive
}

func (mockSrcpConnection *MockSrcpConnection) Connect(url string) {
}

func (mockSrcpConnection *MockSrcpConnection) SendAndReceive(request string) string {
	if mockSrcpConnection.sendAndReceive != nil {
		return mockSrcpConnection.sendAndReceive(request)
	} else {
		mockSrcpConnection.t.Fatal("SendAndReceive used but not defined")
		return ""
	}
}

func (mockSrcpConnection *MockSrcpConnection) Receive() string {
	return ""
}

func (mockSrcpConnection *MockSrcpConnection) Close() {
}
