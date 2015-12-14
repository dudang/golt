package runner
import (
	"sync"
	"bytes"
	"sort"
	"net/http"
	"github.com/dudang/golt/parser"
	logger "github.com/dudang/golt/logger"
	"io/ioutil"
	"time"
	"fmt"
)

const parallelGroup = "parallel"

var stageWaitGroup sync.WaitGroup
var threadWaitGroup sync.WaitGroup
var requestsWaitGroup sync.WaitGroup
var channel = make(chan []byte, 1024)

func ExecuteGoltTest(goltTest parser.Golts, logFile string) {
	m := generateGoltMap(goltTest)

	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	logger.SetOutputFile(logFile)
	go Watch(channel)
	for _, k := range keys {
		executeStage(m[k])
	}
	throughput := CalculateThroughput()
	fmt.Printf("Average Throughput: %f R/S\n", throughput)
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
	// Default timeout of 30 seconds for HTTP calls to avoid hung threads
	return &http.Client{
		Timeout: time.Duration(time.Second * 30),
	}
}

func executeParallelRequests(httpRequest parser.GoltRequest, repetitions int, httpClient *http.Client, stage int) {
	executeHttpRequests(httpRequest, repetitions, httpClient, stage)
	requestsWaitGroup.Done()
}

func executeHttpRequests(httpRequest parser.GoltRequest, repetitions int, httpClient *http.Client, stage int) {
	for i := 1; i <= repetitions; i++ {
		payload := []byte(httpRequest.Payload)
		req, _ := http.NewRequest(httpRequest.Method, httpRequest.URL, bytes.NewBuffer(payload))

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
				Repetition: i,
				ErrorMessage: errorMsg,
				Status: 0,
				Success: false,
				Duration: elapsed}
		} else {
			isSuccess := isCallSuccessful(httpRequest.Assert, resp)
			msg = logger.LogMessage{Stage: stage,
				Repetition: i,
				ErrorMessage: "N/A",
				Status: resp.StatusCode,
				Success: isSuccess,
				Duration: elapsed}
		}

		logger.Log(msg)
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
			isBodySuccessful = assert.Body == string(body)
		}
	}

	isCallSuccessful = isStatusCodeSuccessful && isContentTypeSuccessful && isBodySuccessful
	// TODO: Finish this method to validate the whole assert
	return isCallSuccessful
}