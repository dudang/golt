package runner
import (
	"fmt"
	"bytes"
	"time"
	"regexp"
	"strings"

	"io/ioutil"
	"net/http"

	"github.com/dudang/golt/logger"
	"github.com/dudang/golt/parser"
)

func executeHttpRequests(threadGroup parser.GoltThreadGroup, httpClient *http.Client) {
	for i := 1; i <= threadGroup.Repetitions; i++ {
		executeRequestsSequence(threadGroup.Requests, httpClient, threadGroup.Stage, i)
	}
}

// TODO: Refactor here, starting to have too many responsibilities for a single method
func executeRequestsSequence(httpRequests []parser.GoltRequest, httpClient *http.Client, stage int, repetition int) {
	extractorMap := make(map[string]string)

	for _, request := range httpRequests {

		payloadString := generatePayload(request, extractorMap)

		payload := []byte(payloadString)
		req, _ := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(payload))

		// TODO: need a generate header function to inject it inside the request

		sentRequest := []byte("sent")
		channel <- sentRequest

		start := time.Now()
		resp, err := httpClient.Do(req)
		elapsed := time.Since(start)

		if resp != nil {
			defer resp.Body.Close()
		}

		var msg logger.LogMessage
		if err != nil {
			errorMsg := fmt.Sprintf("%v", err)
			msg = logger.LogMessage{Stage: stage,
				Repetition: repetition,
				ErrorMessage: errorMsg,
				Status: 0,
				Success: false,
				Duration: elapsed}
		} else {
			isSuccess := isCallSuccessful(request.Assert, resp)
			msg = logger.LogMessage{Stage: stage,
				Repetition: repetition,
				ErrorMessage: "N/A",
				Status: resp.StatusCode,
				Success: isSuccess,
				Duration: elapsed}
		}
		logger.Log(msg)

		// Check the Regex and store in Map
		// FIXME: Weak check
		if request.Extract.Field != "" {
			value := executeExtraction(request.Extract, resp)
			if value != "" {
				extractorMap[request.Extract.Var] = value
			}
		}
	}
}

func generatePayload(request parser.GoltRequest, extractorMap map[string]string) (string) {
	regexForVariables := "\\$\\(.*?\\)"
	r, _ := regexp.Compile(regexForVariables)


	if r.MatchString(request.Payload) {
		for _, foundMatch := range r.FindAllString(request.Payload, -1) {
			mapKey := foundMatch[2:len(foundMatch)-1]
			value := extractorMap[mapKey]
			request.Payload = strings.Replace(request.Payload, foundMatch, value, -1)
		}
	}

	// TODO: handle headers to return a map

	return request.Payload
}

func executeExtraction(extractor parser.GoltExtractor, response *http.Response) string{
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
		if err != nil {
			return ""
		}
		if r.MatchString(string(body)) {
			return r.FindString(string(body))
		}
	}
	return ""
}

func isCallSuccessful(assert parser.GoltAssert, response *http.Response) bool {
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