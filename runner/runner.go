package runner
import (
	"fmt"
	"sync"
	"bytes"
	"sort"
	"net/http"
	"github.com/dudang/golt/parser"
	"github.com/dudang/golt/logger"
	"io/ioutil"
)

const parallelGroup = "parallel"

var stageWaitGroup sync.WaitGroup
var threadWaitGroup sync.WaitGroup
var requestsWaitGroup sync.WaitGroup

func ExecuteGoltTest(goltTest parser.Golts, logFile string) {
	m := generateGoltMap(goltTest)

	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	logger.SetOutputFile(logFile)
	for _, k := range keys {
		executeStage(m[k])
	}
	logger.Finish()
}
// FIXME: The three following functions are very repetitive. Find a way to clean it
func executeStage(stage []parser.GoltThreadGroup) {
	stageWaitGroup.Add(len(stage))
	for _, item := range stage {
		httpClient := generateHttpClient(item)
		go executeThreadGroup(item, httpClient)
	}
	stageWaitGroup.Wait()
}

func executeThreadGroup(threadGroup parser.GoltThreadGroup, httpClient *http.Client) {
	threadWaitGroup.Add(threadGroup.Threads)

	for i := 0; i < threadGroup.Threads; i++ {
		go executeRequests(threadGroup, httpClient)
	}

	threadWaitGroup.Wait()
	stageWaitGroup.Done()
}

func executeRequests(threadGroup parser.GoltThreadGroup, httpClient *http.Client) {
	if threadGroup.Type != parallelGroup {
		for _, request := range threadGroup.Requests {
			executeHttpRequests(request, threadGroup.Repetitions, httpClient, threadGroup.Stage)
		}
		threadWaitGroup.Done()
	} else {
		requestsWaitGroup.Add(len(threadGroup.Requests))
		for _, request := range threadGroup.Requests {
			go executeParallelRequests(request, threadGroup.Repetitions, httpClient, threadGroup.Stage)
		}
		requestsWaitGroup.Wait()
		threadWaitGroup.Done()
	}
}


func generateHttpClient(threadGroup parser.GoltThreadGroup) *http.Client {
	// TODO: Currently timeout is not supported with the new data model
	/*var httpClient *http.Client
	if item.Assert.Timeout > 0 {
		httpClient = &http.Client{
			Timeout: time.Duration(time.Millisecond * time.Duration(item.Assert.Timeout)),
		}
	} else {
		httpClient = &http.Client{}
	}*/
	return &http.Client{}
}

func executeParallelRequests(httpRequest parser.GoltRequest, repetitions int, httpClient *http.Client, stage int) {
	executeHttpRequests(httpRequest, repetitions, httpClient, stage)
	requestsWaitGroup.Done()
}

func executeHttpRequests(httpRequest parser.GoltRequest, repetitions int, httpClient *http.Client, stage int) {
	for i := 1; i <= repetitions; i++ {
		payload := []byte(httpRequest.Payload)
		req, _ := http.NewRequest(httpRequest.Method, httpRequest.URL, bytes.NewBuffer(payload))
		resp, err := httpClient.Do(req)

		if resp != nil {
			defer resp.Body.Close()
		}

		var msg string
		if err != nil {
			// TODO: Make a custom struct for log messages
			msg = fmt.Sprintf("Stage: %d Repetitions: %d Message: %v Success: %t", stage, i, err, false)
		} else {
			isSuccess := isCallSuccessful(httpRequest.Assert, resp)
			msg = fmt.Sprintf("Stage: %d Repetitions: %d  Status Code: %d Success: %t", stage, i, resp.StatusCode, isSuccess)
		}

		logger.Log([]byte(msg))
	}
}

func isCallSuccessful(assert parser.GoltAssert, response *http.Response) bool {
	var isCallSuccessful bool
	isContentTypeSuccessful := true
	isBodySuccessful := true
	isStatusCodeSuccessful := assert.Status == response.StatusCode

	if assert.Type != "" {
		isContentTypeSuccessful = assert.Type == response.Header.Get("content-type")
	}

	if assert.Body != "" {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			isBodySuccessful = false
		} else {
			isBodySuccessful = assert.Body == body
		}
	}

	isCallSuccessful = isStatusCodeSuccessful && isContentTypeSuccessful && isBodySuccessful
	// TODO: Finish this method to validate the whole assert
	return isCallSuccessful
}