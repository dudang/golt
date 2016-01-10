package main
import (
	"fmt"
	"time"
	"regexp"
	"io/ioutil"
	"net/http"
)


const regexForVariable = "\\$\\(.*?\\)"
var r, _ = regexp.Compile(regexForVariable)

func executeHttpRequests(threadGroup GoltThreadGroup, sender GoltSender) {
	for i := 1; i <= threadGroup.Repetitions; i++ {
		executeRequestsSequence(threadGroup.Requests, sender, threadGroup.Stage, i)
	}
}

func executeRequestsSequence(httpRequests []GoltRequest, sender GoltSender, stage int, repetition int) {
	// TODO: By defining the map here, it's local to the thread, maybe we want something else
	extractorMap := make(map[string]string)
	extractionWasDone := false

	for _, request := range httpRequests {
		req := BuildRequest(extractionWasDone, request, extractorMap)
		notifyWatcher()

		start := time.Now()
		resp, err := sender.Send(req)
		elapsed := time.Since(start)

		if resp != nil {
			defer resp.Body.Close()
		}

		// Log result
		logResult(request, resp, err, stage, repetition, elapsed)

		// Handle all the regular expression extraction
		extractionWasDone = handleExtraction(request, resp, extractorMap)
	}
}

func notifyWatcher() {
	sentRequest := []byte("sent")
	channel <- sentRequest
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

func handleExtraction(request GoltRequest, resp *http.Response, extractorMap map[string]string) bool{
	// Check if we are extracting anything and store it in a Map
	regexIsDefined := request.Extract.Field != "" && request.Extract.Regex != "" && request.Extract.Var != ""
	if regexIsDefined {
		value := executeExtraction(request.Extract, resp)
		if value != "" {
			extractorMap[request.Extract.Var] = value
			return true
		}
	}
	return false
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
