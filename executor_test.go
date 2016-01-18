package main

import (
	"testing"
	"net/http"
	"os"
	"log"
)

var countedRequest = 0

var jsonHeaders = http.Header{"Content-Type":[]string{"application/json"}}
var htmlHeaders = http.Header{"Content-Type":[]string{"text/html"}}
var utilGenerator = &GoltGenerator{}

var assertionTestingTable = []struct{
	expectedSuccess bool
	response *http.Response
} {
	{expectedSuccess: true, response: &http.Response{StatusCode: 200, Header: jsonHeaders,}},
	{expectedSuccess: false, response: &http.Response{StatusCode: 404, Header: jsonHeaders,}},
	{expectedSuccess: false, response: &http.Response{StatusCode: 200, Header: htmlHeaders,}},
	{expectedSuccess: false, response: &http.Response{StatusCode: 404, Header: htmlHeaders,}},
}

var requestTestingTable = []struct {
	extractor GoltExtractor
	response *http.Response
	expectedValue string
	extractionWasExecuted bool
} {
	{
		extractor: GoltExtractor{Field: "headers", Var: "CONTENT_TYPE", Regex: "text/html(.*)"},
		response: &http.Response{Header: htmlHeaders},
		expectedValue: "text/html",
		extractionWasExecuted: true,
	},
	{
		extractor: GoltExtractor{Field: "body", Var: "CONTENT_TYPE", Regex: "text/html(.*)"},
		response: &http.Response{Body: utilGenerator.buildRegularRequest(GoltRequest{Payload: "text/htmlABC"}).Body},
		expectedValue: "text/htmlABC",
		extractionWasExecuted: true,
	},
	{
		extractor: GoltExtractor{Field: "headers", Var: "CONTENT_TYPE", Regex: "text/html(.*)"},
		response: &http.Response{Header: jsonHeaders},
		expectedValue: "",
		extractionWasExecuted: false,
	},
	{
		extractor: GoltExtractor{Field: "headers", Var: "", Regex: "text/html(.*)"},
		response: &http.Response{Header: jsonHeaders},
		expectedValue: "",
		extractionWasExecuted: false,
	},
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

func TestHandleExtraction(t *testing.T) {
	for _, entry := range requestTestingTable {
		regexMap := make(map[string]string)
		extractionExecuted := handleExtraction(entry.extractor, entry.response, regexMap)
		if extractionExecuted != entry.extractionWasExecuted || regexMap["CONTENT_TYPE"] != entry.expectedValue {
			t.Error("Result returned by the method is not what was expected")
		}
	}
}
