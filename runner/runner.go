package runner
import (
	"net/http"
	"fmt"
	"sync"
	"bytes"
	"github.com/dudang/golt/parser"
)

var wg sync.WaitGroup

type httpRequest func(string) (*http.Response, error)

func ExecuteGoltTest(goltTest parser.Golt) {
	for _, element := range goltTest.Golt {
		executeElement(element)
	}
}

func executeElement(element parser.GoltJson) {
	wg.Add(element.Threads)
	for i:= 0; i < element.Threads; i++ {
		go executeHttpRequest(element)
	}
	wg.Wait()
}

func executeHttpRequest(element parser.GoltJson) {
	for i := 1; i <= element.Repetitions; i++ {
		payload := []byte(element.Payload)
		req, err := http.NewRequest(element.Method, element.URL, bytes.NewBuffer(payload))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		defer resp.Body.Close()
		fmt.Printf("Repetitions: %d  Status Code: %d Success: %t\n", i, resp.StatusCode, resp.StatusCode == element.Assert.Status)
	}
	wg.Done()
}