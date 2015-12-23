package main

import (
	"testing"
	"fmt"
	"net/http"
)

const testKey1, testValue1 = "TEST_VALUE_1", "test-string-1"
const testKey2, testValue2 = "TEST_VALUE_2", "test-string-2"

func TestReplaceRegex(t *testing.T) {
	testingMap := make(map[string]string)
	testingMap[testKey1] = testValue1
	testingMap[testKey2] = testValue2

	replaceString1 := fmt.Sprintf("This is a test to replace the regex $(%s)", testKey1)
	replaceRegex(r, &replaceString1, testingMap)
	expectedString1 := fmt.Sprintf("This is a test to replace the regex %s", testValue1)
	if replaceString1 != expectedString1 {
		t.Fail()
	}

	replaceString2 := fmt.Sprintf("This is a test $(%s) to replace the regex $(%s)", testKey2, testKey2)
	replaceRegex(r, &replaceString2, testingMap)
	expectedString2 := fmt.Sprintf("This is a test %s to replace the regex %s", testValue2, testValue2)
	if replaceString2 != expectedString2 {
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
