package main

import (
	"testing"
	"net/http"
	"os"
	"log"
)

const testKey1, testValue1 = "TEST_VALUE_1", "test-string-1"
const testKey2, testValue2 = "TEST_VALUE_2", "test-string-2"
var testingMap = map[string]string{testKey1: testValue1, testKey2: testKey2}

var countedRequest = 0

var jsonHeaders = http.Header{"Content-Type":[]string{"application/json"}}
var htmlHeaders = http.Header{"Content-Type":[]string{"text/html"}}
var assertionTestingTable = []struct{
	expectedSuccess bool
	response *http.Response
} {
	{expectedSuccess: true, response: &http.Response{StatusCode: 200, Header: jsonHeaders,}},
	{expectedSuccess: false, response: &http.Response{StatusCode: 404, Header: jsonHeaders,}},
	{expectedSuccess: false, response: &http.Response{StatusCode: 200, Header: htmlHeaders,}},
	{expectedSuccess: false, response: &http.Response{StatusCode: 404, Header: htmlHeaders,}},
}

func TestIsCallSuccessful(t *testing.T) {
	assert := GoltAssert{
		Type: "application/json",
		Status: 200,
	}
	for _, entry := range assertionTestingTable {
		if isCallSuccessful(assert, entry.response) != entry.expectedSuccess {
			t.Error("The assertion was not validated properly")
		}
	}
}

func TestExecuteRequestsSequence(t *testing.T) {
	requests := []GoltRequest{GoltRequest{}, GoltRequest{}}
	executor := GoltExecutor{
		Sender:	MockSender{},
		Logger: &GoltLogger{Logger: log.New(os.Stdout, "", 0),},
		SendingChannel: make(chan []byte, 1024),
	}
	executor.executeRequestsSequence(requests)
	if countedRequest != 2 {
		t.Error("Request sent should have been 2")
	}
}


type MockReader struct {}
func (reader MockReader) Read(p []byte) (int, error) {return 0, nil}

type MockWriter struct {}
func (writer MockWriter) Close() error {return nil}

type MockReadCloser struct {
	MockReader
	MockWriter
}

type MockSender struct {}

func (sender MockSender) Send(req *http.Request) (*http.Response, error) {
	countedRequest += 1
	return &http.Response{Body: MockReadCloser{}}, nil
}
