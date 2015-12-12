package runner
import (
	"time"
	"fmt"
	"sync"
	"bytes"
	"sort"
	"net/http"
	"github.com/dudang/golt/parser"
	"github.com/dudang/golt/logger"
)

const dateFormat = "2006-01-02 15:04:05"
const parallelGroup = "parallel"

var stageWaitGroup sync.WaitGroup
var threadWaitGroup sync.WaitGroup

func ExecuteGoltTest(goltTest parser.Golts) {
	m := generateGoltMap(goltTest)

	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		executeStage(m[k])
	}
	// We need to flush the remaining messages still in buffer after the test is over
	logger.Flush()
}
// FIXME: The three following functions are very repetitive. Find a way to clean it
func executeStage(stage []parser.GoltThreadGroup) {
	stageWaitGroup.Add(len(stage))
	for _, item := range stage{
		httpClient := generateHttpClient(item)
		go executeThreadGroup(item, httpClient)
	}
	stageWaitGroup.Wait()
}

func executeThreadGroup(threadGroup parser.GoltThreadGroup, httpClient *http.Client) {
	threadWaitGroup.Add(threadGroup.Threads)

	for i:= 0; i < threadGroup.Threads; i++ {
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
	}
	// TODO: Handle the parallel thread groups
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
			msg = fmt.Sprintf("[%s] Stage: %d Repetitions: %d Message: %v\n", time.Now().Format(dateFormat), stage, i, err)
		} else {
			isSuccess := isCallSuccessful(httpRequest.Assert, resp)
			msg = fmt.Sprintf("[%s] Stage: %d Repetitions: %d  Status Code: %d Success: %t\n", time.Now().Format(dateFormat), stage, i, resp.StatusCode, isSuccess)
		}

		logger.Log([]byte(msg))
	}
}

func isCallSuccessful(assert parser.GoltAssert, response *http.Response) bool{
	// TODO: Finish this method to validate more than just response code
	return assert.Status == response.StatusCode
}