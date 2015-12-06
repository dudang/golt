package runner
import (
	"net/http"
	"fmt"
	"sync"
	"bytes"
	"sort"
	"github.com/dudang/golt/parser"
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
}

func executeStage(stage []parser.GoltItem) {
	stageWaitGroup.Add(len(stage))
	for i:= 0; i < len(stage); i++ {
		go executeElement(stage[i])
	}
	stageWaitGroup.Wait()
}

func executeElement(element parser.GoltItem) {
	internalWaitGroup.Add(element.Threads)
	for i:= 0; i < element.Threads; i++ {
		go executeHttpRequest(element)
	}
	internalWaitGroup.Wait()
	stageWaitGroup.Done()
}

func executeHttpRequest(element parser.GoltItem) {
	for i := 1; i <= element.Repetitions; i++ {
		payload := []byte(element.Payload)
		req, err := http.NewRequest(element.Method, element.URL, bytes.NewBuffer(payload))

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		defer resp.Body.Close()
		fmt.Printf("Stage: %d Repetitions: %d  Status Code: %d Success: %t\n", element.Stage, i, resp.StatusCode, resp.StatusCode == element.Assert.Status)
	}
	internalWaitGroup.Done()
}