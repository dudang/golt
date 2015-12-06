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

var internalWaitGroup sync.WaitGroup
var stageWaitGroup sync.WaitGroup

var httpClient = &http.Client{}

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
		go executeItem(item)
	}
	stageWaitGroup.Wait()
}

func executeItem(item parser.GoltItem) {
	internalWaitGroup.Add(item.Threads)
	for i:= 0; i < item.Threads; i++ {
		go executeHttpRequest(item)
	}
	internalWaitGroup.Wait()
	stageWaitGroup.Done()
}

func executeHttpRequest(item parser.GoltItem) {
	for i := 1; i <= item.Repetitions; i++ {
		payload := []byte(item.Payload)
		req, err := http.NewRequest(item.Method, item.URL, bytes.NewBuffer(payload))

		resp, err := httpClient.Do(req)
		var msg string
		if err != nil {
			msg = fmt.Sprintf("%v\n", err)
		}
		defer resp.Body.Close()

		msg = fmt.Sprintf("[%s] Stage: %d Repetitions: %d  Status Code: %d Success: %t\n",
			time.Now().Format("2006-01-02 15:04:05"),
			item.Stage,
			i,
			resp.StatusCode,
			resp.StatusCode == item.Assert.Status)
		logger.Log([]byte(msg))
	}
	internalWaitGroup.Done()
}