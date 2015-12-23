package main
import (
	"fmt"
	"bytes"
	"time"
	"regexp"
	"strings"

	"io/ioutil"
	"net/http"
)


const regexForVariable = "\\$\\(.*?\\)"
var r, _ = regexp.Compile(regexForVariable)

func executeHttpRequests(threadGroup GoltThreadGroup, httpClient *http.Client) {
	for i := 1; i <= threadGroup.Repetitions; i++ {
		executeRequestsSequence(threadGroup.Requests, httpClient, threadGroup.Stage, i)
	}
}

// TODO: Refactor here, starting to have too many responsibilities for a single method
func executeRequestsSequence(httpRequests []GoltRequest, httpClient *http.Client, stage int, repetition int) {
	// TODO: By defining the map here, it's local to the thread, maybe we want something else
	extractorMap := make(map[string]string)
	extractionWasDone := false

	for _, request := range httpRequests {
		var req *http.Request
		if extractionWasDone {
			req = buildRegexRequest(request, extractorMap)
		} else {
			req = buildRequest(request)
		}

		// Notify the watcher that the request is sent for throughput duties
		func() {
			sentRequest := []byte("sent")
			channel <- sentRequest
		}()

		// Send request and calculate time
		start := time.Now()
		resp, err := sendRequest(req, httpClient)
		elapsed := time.Since(start)

		if resp != nil {
			defer resp.Body.Close()
		}

		// Log result
		logResult(request, resp, err, stage, repetition, elapsed)

		// Check if we are extracting anything and store it in a Map
		regexIsDefined := request.Extract.Field != "" && request.Extract.Regex != "" && request.Extract.Var != ""
		if regexIsDefined {
			value := executeExtraction(request.Extract, resp)
			if value != "" {
				extractorMap[request.Extract.Var] = value
				extractionWasDone = true
			}
		}
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

func buildRequest(request GoltRequest) *http.Request {
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

// TODO: Possibly make this more generic in the future for other protocols
func sendRequest(req *http.Request, client *http.Client) (*http.Response, error) {
	return client.Do(req)
}

// TODO: Too many parameters on this method, to refactor
func logResult(request GoltRequest, resp *http.Response, err error, stage int, repetition int, elapsed time.Duration) {
	var msg LogMessage
	if err != nil {
		errorMsg := fmt.Sprintf("%v", err)
		msg = LogMessage{Stage: stage,
			Repetition: repetition,
			ErrorMessage: errorMsg,
			Status: 0,
			Success: false,
			Duration: elapsed}
	} else {
		isSuccess := isCallSuccessful(request.Assert, resp)
		msg = LogMessage{Stage: stage,
			Repetition: repetition,
			ErrorMessage: "N/A",
			Status: resp.StatusCode,
			Success: isSuccess,
			Duration: elapsed}
	}
	Log(msg)
}

func isCallSuccessful(assert GoltAssert, response *http.Response) bool {
	var isCallSuccessful bool
	isContentTypeSuccessful := true
	isBodySuccessful := true
	isStatusCodeSuccessful := assert.Status == response.StatusCode

	if assert.Type != "" {
		isContentTypeSuccessful = assert.Type == response.Header.Get("content-type")
	}

	isCallSuccessful = isStatusCodeSuccessful && isContentTypeSuccessful && isBodySuccessful
	return isCallSuccessful
}

func executeExtraction(extractor GoltExtractor, response *http.Response) string{
	r, _ := regexp.Compile(extractor.Regex)
	switch extractor.Field {
	case "headers":
		// FIXME: Find a cleaner algorithm
		for k, v := range response.Header {
			for _, value := range v {
				value = fmt.Sprintf("%s: %s", k, value)
				if r.MatchString(value) {
					return r.FindString(value)
				}
			}
		}
	case "body":
		body, err := ioutil.ReadAll(response.Body)
		if r.MatchString(string(body)) && err == nil {
			return r.FindString(string(body))
		}
	}
	return ""
}
