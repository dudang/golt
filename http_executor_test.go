package main

import (
	"testing"
	"net/http"
)

const testKey1, testValue1 = "TEST_VALUE_1", "test-string-1"
const testKey2, testValue2 = "TEST_VALUE_2", "test-string-2"
var testingMap map[string]string
var countedRequest = 0

func init() {
	testingMap = make(map[string]string)
	testingMap[testKey1] = testValue1
	testingMap[testKey2] = testValue2
}

func TestIsCallSuccessful(t *testing.T) {
	assert := GoltAssert{
		Type: "application/json",
		Status: 200,
	}

	var jsonHeaders http.Header
	jsonHeaders = http.Header{}
	jsonHeaders.Set("Content-Type", "application/json")

	var htmlHeaders http.Header
	htmlHeaders = http.Header{}
	htmlHeaders.Set("Content-Type", "text/html")

	validResponse := &http.Response{
		StatusCode: 200,
		Header: jsonHeaders,
	}

	wrongStatusCodeResponse := &http.Response{
		StatusCode: 404,
		Header: jsonHeaders,
	}

	wrongContentTypeResponse := &http.Response{
		StatusCode: 200,
		Header: htmlHeaders,
	}

	wrongResponse := &http.Response{
		StatusCode: 404,
		Header: htmlHeaders,
	}

	if isCallSuccessful(assert, validResponse) != true {
		t.Fail()
	}

	if isCallSuccessful(assert, wrongStatusCodeResponse) == true {
		t.Fail()
	}

	if isCallSuccessful(assert, wrongContentTypeResponse) == true {
		t.Fail()
	}

	if isCallSuccessful(assert, wrongResponse) == true {
		t.Fail()
	}
}

func TestExecuteRequestsSequence(t *testing.T) {
	requests := []GoltRequest{GoltRequest{}}
	sender := MockSender{}
	executeRequestsSequence(requests, sender, 1, 1)
	if countedRequest != 1 {
		t.Error("Request sent should have been 1")
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
