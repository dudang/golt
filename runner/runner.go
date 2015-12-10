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

var internalWaitGroup sync.WaitGroup
var stageWaitGroup sync.WaitGroup

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

func executeStage(stage []parser.GoltItem) {
	stageWaitGroup.Add(len(stage))
	for _, item := range stage{
		httpClient := generateHttpClient(item)
		go executeItem(item, httpClient)
	}
	stageWaitGroup.Wait()
}

func executeItem(item parser.GoltItem, httpClient *http.Client) {
	internalWaitGroup.Add(item.Threads)

	for i:= 0; i < item.Threads; i++ {
		go executeHttpRequest(item, httpClient)
	}

	internalWaitGroup.Wait()
	stageWaitGroup.Done()
}

func generateHttpClient(item parser.GoltItem) *http.Client {
	var httpClient *http.Client
	if item.Assert.Timeout > 0 {
		httpClient = &http.Client{
			Timeout: time.Duration(time.Millisecond * time.Duration(item.Assert.Timeout)),
		}
	} else {
		httpClient = &http.Client{}
	}
	return httpClient
}

func executeHttpRequest(item parser.GoltItem, httpClient *http.Client) {
	for i := 1; i <= item.Repetitions; i++ {
		payload := []byte(item.Payload)

		req, err := http.NewRequest(item.Method, item.URL, bytes.NewBuffer(payload))
		resp, err := httpClient.Do(req)

		if resp != nil {
			defer resp.Body.Close()
		}

		var msg string
		if err != nil {
			// TODO: Make a custom struct for log messages
			msg = fmt.Sprintf("[%s] Stage: %d Repetitions: %d Message: %v\n", time.Now().Format(dateFormat), item.Stage, i, err)
		} else {
			isSuccess := isCallSuccessful(item.Assert, resp)
			msg = fmt.Sprintf("[%s] Stage: %d Repetitions: %d  Status Code: %d Success: %t\n", time.Now().Format(dateFormat), item.Stage, i, resp.StatusCode, isSuccess)
		}

		logger.Log([]byte(msg))
	}
	internalWaitGroup.Done()
}

func isCallSuccessful(assert parser.GoltAssert, response *http.Response) bool{
	// TODO: Finish this method to validate more than just response code
	return assert.Status == response.StatusCode
}