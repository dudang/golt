package main

import (
"testing"
"net/http"
	"log"
	"os"
)

var countedRequest = 0

type MockSender struct {}
type MockReader struct {}
type MockWriter struct {}
type MockReadCloser struct {
	MockReader
	MockWriter
}

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
		Sender:	MockSender{},
		Logger: &GoltLogger{Logger: log.New(os.Stdout, "", 0),},
		SendingChannel: make(chan []byte, 1024),
	}

	executor.ExecuteHttpRequests()
	if countedRequest != 10 {
		t.Error("Amount of request sent is different than expected")
	}
}
