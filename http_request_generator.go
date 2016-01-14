package main
import (
	"bytes"
	"regexp"
	"strings"
	"net/http"
)


const regexForVariable = "\\$\\(.*?\\)"
var r, _ = regexp.Compile(regexForVariable)

func BuildRequest(shouldRegex bool, request GoltRequest, regexMap map[string]string) *http.Request{
	if shouldRegex {
		return buildRegexRequest(request, regexMap)
	} else {
		return buildRegularRequest(request)
	}
}

func buildRegexRequest(request GoltRequest, extractorMap map[string]string) *http.Request{
	payloadString := generatePayload(request, extractorMap)
	payload := []byte(payloadString)

	req, _ := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(payload))

	headers := generateHeaders(request, extractorMap)
	for k, v := range headers {
		req.Header.Set(k, *v)
	}
	return req
}

func buildRegularRequest(request GoltRequest) *http.Request {
	payload := []byte(request.Payload)
	req, _ := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(payload))
	for k, v := range request.Headers {
		req.Header.Set(k, *v)
	}
	return req
}

func generatePayload(request GoltRequest, extractorMap map[string]string) (string) {
	// We are passing the pointer of the Payload to modify it's value
	replaceRegex(r, &request.Payload, extractorMap)
	return request.Payload
}

func generateHeaders(request GoltRequest, extractorMap map[string]string) map[string]*string {
	for k := range request.Headers {
		// We are passing a pointer of the value in the map to replace it's value
		replaceRegex(r, request.Headers[k], extractorMap)
	}
	return request.Headers
}

func replaceRegex(regex *regexp.Regexp, value *string, extractorMap map[string]string) {
	/*
	Given a specific regular expression, a pointer to a string and a map of stored variable
	This method will have the side effect of changing the value pointer by the string if the regex is matching
	*/
	if regex.MatchString(*value) {
		for _, foundMatch := range regex.FindAllString(*value, -1) {
			mapKey := foundMatch[2:len(foundMatch)-1]
			extractedValue := extractorMap[mapKey]
			*value = strings.Replace(*value, foundMatch, extractedValue, -1)
		}
	}
}
