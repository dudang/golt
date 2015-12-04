package runner
import (
	"net/http"
	"fmt"
	"sync"
	"github.com/dudang/golt/parser"
)

var wg sync.WaitGroup

func ExecuteJsonGolt(testPlan parser.GoltJsons) {
	for _, element := range testPlan.Golt {
		executeElement(element)
	}
}

func executeElement(testElement parser.GoltJson) {
	wg.Add(testElement.Threads)
	for i:= 0; i < testElement.Threads; i++ {
		go spawnRoutine(testElement)
	}
	wg.Wait()
}

func spawnRoutine(testElement parser.GoltJson) {
	switch testElement.Method {
		case "GET":
			getRequest(testElement.URL, testElement.Repetitions)
		default:
			return
	}
	wg.Done()
}

func getRequest(url string, repetitions int) {
	for i := 1; i <= repetitions; i++ {
		resp, err := http.Get(url)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("Repetitions: %d  Status Code: %d\n", i, resp.StatusCode)
	}
}