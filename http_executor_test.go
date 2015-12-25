package main

import (
	"testing"
	"fmt"
	"net/http"
	"io/ioutil"
)

const testKey1, testValue1 = "TEST_VALUE_1", "test-string-1"
const testKey2, testValue2 = "TEST_VALUE_2", "test-string-2"
var testingMap map[string]string

func init() {
	testingMap = make(map[string]string)
	testingMap[testKey1] = testValue1
	testingMap[testKey2] = testValue2
}

func TestBuildRegexRequest(t *testing.T) {
	testHeader := "test"
	headerValue := fmt.Sprintf("$(%s)", testKey2)
	request := GoltRequest{
		Payload: fmt.Sprintf("$(%s)", testKey1),
		Headers: map[string]*string{testHeader: &headerValue},
	}

	httpRequest := BuildRegexRequest(request, testingMap)

	body, _ := ioutil.ReadAll(httpRequest.Body)
	if string(body) != testValue1 {
		t.Fail()
	}
	if httpRequest.Header.Get(testHeader) != testValue2 {
		t.Fail()
	}
}

func TestBuildRequest(t *testing.T) {
	testHeaderKey, testHeaderValue, testPayload := "headerKey", "headerValue", "payload"
	request := GoltRequest{
		Payload: testPayload,
		Headers: map[string]*string{testHeaderKey: &testHeaderValue},
	}
	httpRequest := BuildRequest(request)
	body, _ := ioutil.ReadAll(httpRequest.Body)
	if string(body) != testPayload {
		t.Fail()
	}
	if httpRequest.Header.Get(testHeaderKey) != testHeaderValue {
		t.Fail()
	}

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
