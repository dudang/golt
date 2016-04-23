package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

const testKey1, testValue1 = "TESTING_KEY_1", "TESTING_VALUE_1"
const testKey2, testValue2 = "TESTING_KEY_2", "TESTING_VALUE_2"

var generator *GoltGenerator

func init() {
	testingMap := make(map[string]string)
	testingMap[testKey1] = testValue1
	testingMap[testKey2] = testValue2
	generator = &GoltGenerator{
		RegexMap: testingMap,
	}
}

func TestBuildRequest(t *testing.T) {
	regexHeader := "test"
	headerValue := fmt.Sprintf("$(%s)", testKey2)
	regexRequest := GoltRequest{
		Payload: fmt.Sprintf("$(%s)", testKey1),
		Headers: map[string]*string{regexHeader: &headerValue},
	}
	httpRequest := generator.BuildRequest(true, regexRequest)
	regexBody, _ := ioutil.ReadAll(httpRequest.Body)
	if string(regexBody) != testValue1 || httpRequest.Header.Get(regexHeader) != testValue2 {
		t.Error("Regex request returned is not valid")
	}

	testHeaderKey, testHeaderValue, testPayload := "headerKey", "headerValue", "payload"
	request := GoltRequest{
		Payload: testPayload,
		Headers: map[string]*string{testHeaderKey: &testHeaderValue},
	}
	generatedRequest := generator.BuildRequest(false, request)
	generatedBody, _ := ioutil.ReadAll(generatedRequest.Body)
	if string(generatedBody) != testPayload || generatedRequest.Header.Get(testHeaderKey) != testHeaderValue {
		t.Error("Regular request returned is not valid")
	}

}
