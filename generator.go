package main

import (
	"bytes"
	"net/http"
	"regexp"
	"strings"
)

const regexForVariable = "\\$\\(.*?\\)"

var r, _ = regexp.Compile(regexForVariable)

type GoltGenerator struct {
	RegexMap map[string]string
}

func (g *GoltGenerator) BuildRequest(shouldRegex bool, request GoltRequest) *http.Request {
	if shouldRegex {
		return g.buildRegexRequest(request)
	} else {
		return g.buildRegularRequest(request)
	}
}

func (g *GoltGenerator) buildRegexRequest(request GoltRequest) *http.Request {
	payloadString := g.generatePayload(request)
	payload := []byte(payloadString)

	req, _ := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(payload))

	headers := g.generateHeaders(request)
	for k, v := range headers {
		req.Header.Set(k, *v)
	}
	return req
}

func (g *GoltGenerator) buildRegularRequest(request GoltRequest) *http.Request {
	payload := []byte(request.Payload)
	req, _ := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(payload))
	for k, v := range request.Headers {
		req.Header.Set(k, *v)
	}
	return req
}

func (g *GoltGenerator) generatePayload(request GoltRequest) string {
	// We are passing the pointer of the Payload to modify it's value
	g.replaceRegex(r, &request.Payload)
	return request.Payload
}

func (g *GoltGenerator) generateHeaders(request GoltRequest) map[string]*string {
	for k := range request.Headers {
		// We are passing a pointer of the value in the map to replace it's value
		g.replaceRegex(r, request.Headers[k])
	}
	return request.Headers
}

func (g *GoltGenerator) replaceRegex(regex *regexp.Regexp, value *string) {
	/*
		Given a specific regular expression, a pointer to a string and a map of stored variable
		This method will have the side effect of changing the value pointer by the string if the regex is matching
	*/
	if regex.MatchString(*value) {
		for _, foundMatch := range regex.FindAllString(*value, -1) {
			mapKey := foundMatch[2 : len(foundMatch)-1]
			extractedValue := g.RegexMap[mapKey]
			*value = strings.Replace(*value, foundMatch, extractedValue, -1)
		}
	}
}
