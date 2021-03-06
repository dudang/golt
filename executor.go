package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

// Interface to send Requests (only HTTP right now, should extend to more)
// FIXME: Send method should take generic arguments, not HTTP specific arguments
type GoltSender interface {
	Send(request *http.Request) (*http.Response, error)
}

// Implementation of a GoltSender with a HTTP Client
type HttpSender struct {
	Client *http.Client
}

func (http HttpSender) Send(request *http.Request) (*http.Response, error) {
	return http.Client.Do(request)
}

type GoltExecutor struct {
	ThreadGroup    GoltThreadGroup
	Sender         GoltSender
	Logger         GoltLogger
	SendingChannel chan []byte
}

func (e *GoltExecutor) ExecuteHttpRequests() {
	for i := 0; i < e.ThreadGroup.Repetitions; i++ {
		e.executeRequestsSequence(e.ThreadGroup.Requests)
	}
}

func (e *GoltExecutor) executeRequestsSequence(httpRequests []GoltRequest) {
	// TODO: By defining the map here, it's local to the thread, maybe we want something else
	regexMap := make(map[string]string)
	generator := &GoltGenerator{RegexMap: regexMap}
	extractionWasDone := false
	for _, request := range httpRequests {
		req := generator.BuildRequest(extractionWasDone, request)
		notifyWatcher(e.SendingChannel)

		start := time.Now()
		resp, err := e.Sender.Send(req)
		elapsed := time.Since(start)

		if resp != nil {
			defer resp.Body.Close()

			// Handle all the regular expression extraction
			extractionWasDone = handleExtraction(request.Extract, resp, regexMap)
		}

		// Log result
		e.logResult(request, resp, err, elapsed)
	}
}

func (e *GoltExecutor) logResult(request GoltRequest, resp *http.Response, err error, elapsed time.Duration) {
	var msg LogMessage
	if err != nil {
		errorMsg := fmt.Sprintf("%v", err)
		msg = LogMessage{
			Url:          request.URL,
			ErrorMessage: errorMsg,
			Status:       0,
			Success:      false,
			Duration:     elapsed}
	} else {
		isSuccess := isCallSuccessful(request.Assert, resp)
		msg = LogMessage{
			Url:          request.URL,
			ErrorMessage: "N/A",
			Status:       resp.StatusCode,
			Success:      isSuccess,
			Duration:     elapsed}
	}
	e.Logger.Log(msg)
}

func notifyWatcher(channel chan []byte) {
	sentRequest := []byte("sent")
	channel <- sentRequest
}

func isCallSuccessful(assert GoltAssert, response *http.Response) bool {
	var isCallSuccessful bool
	isContentTypeSuccessful := true

	// Verify here if the status expect in the assert is the same as the response or it's not defined (not checking)
	isStatusCodeSuccessful := assert.Status == response.StatusCode || assert.Status == 0

	if assert.Type != "" {
		isContentTypeSuccessful = assert.Type == response.Header.Get("Content-Type")
	}

	isCallSuccessful = isStatusCodeSuccessful && isContentTypeSuccessful
	return isCallSuccessful
}

func handleExtraction(extractor GoltExtractor, resp *http.Response, regexMap map[string]string) bool {
	// Check if we are extracting anything and store it in a Map
	regexIsDefined := extractor.Field != "" && extractor.Regex != "" && extractor.Var != ""
	if regexIsDefined {
		value := executeExtraction(extractor, resp)
		if value != "" {
			regexMap[extractor.Var] = value
			return true
		}
	}
	return false
}

func executeExtraction(extractor GoltExtractor, response *http.Response) string {
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
