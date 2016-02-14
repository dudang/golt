package main

import (
	"testing"
	"net/http"
	"errors"
)

var countedRequest = 0
var successfulRequest = 0

// Mock struct needed for a GoltSender
type MockSender struct{}
type MockErrorSender struct{}
type MockReader struct{}
type MockWriter struct{}
type MockReadCloser struct {
	MockReader
	MockWriter
}

// Mock implementations for GoltSender methods
func (reader MockReader) Read(p []byte) (int, error) {
	return 0, nil
}
func (writer MockWriter) Close() error {
	return nil
}
func (sender MockSender) Send(req *http.Request) (*http.Response, error) {
	countedRequest += 1
	return &http.Response{Body: MockReadCloser{}}, nil
}
func (sender MockErrorSender) Send(req *http.Request) (*http.Response, error) {
	countedRequest += 1
	return nil, errors.New("Error sending request")
}

type MockLogger struct{}

func (logger MockLogger) SetOutputFile(filename string) {}
func (logger MockLogger) Finish() {}
func (logger MockLogger) Log(message LogMessage) {
	if message.Success == true {
		successfulRequest += 1
	}
}

func TestExecuteHttpRequests(t *testing.T) {
	threadGroup := GoltThreadGroup{
		Threads: 5,
		Timeout: 500,
		Repetitions: 5,
		Stage: 1,
		Requests: []GoltRequest{GoltRequest{}, GoltRequest{}},
	}

	executor := GoltExecutor{
		ThreadGroup: threadGroup,
		Sender:    MockSender{},
		Logger:    MockLogger{},
		SendingChannel: make(chan []byte, 1024),
	}

	executorWithError := GoltExecutor{
		ThreadGroup: threadGroup,
		Sender:    MockErrorSender{},
		Logger:    MockLogger{},
		SendingChannel: make(chan []byte, 1024),
	}

	if !resetAndTest(executor, 10, 10) || !resetAndTest(executorWithError, 10, 0) {
		t.Error("Error executing the test plan")
	}
}

func resetAndTest(executor GoltExecutor, expectedRequests int, expectedSuccess int) bool {
	successfulRequest = 0
	countedRequest = 0
	executor.ExecuteHttpRequests()
	if countedRequest != expectedRequests  || successfulRequest != expectedSuccess {
		return false
	}
	return true
}
