package main
import (
	"testing"
	"fmt"
	"io/ioutil"
)

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

	httpRequest := buildRegexRequest(request, testingMap)

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
	httpRequest := buildRegularRequest(request)
	body, _ := ioutil.ReadAll(httpRequest.Body)
	if string(body) != testPayload {
		t.Fail()
	}
	if httpRequest.Header.Get(testHeaderKey) != testHeaderValue {
		t.Fail()
	}
}
