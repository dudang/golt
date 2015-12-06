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
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		defer resp.Body.Close()
		fmt.Printf("Stage: %d Repetitions: %d  Status Code: %d Success: %t\n", item.Stage, i, resp.StatusCode, resp.StatusCode == item.Assert.Status)
	}
	internalWaitGroup.Done()
}