package runner
import (
	"sync"
	"sort"
	"time"
	"net/http"

	parser "github.com/dudang/golt/parser"
	logger "github.com/dudang/golt/logger"
)

const parallelGroup = "parallel"

var stageWaitGroup sync.WaitGroup
var threadWaitGroup sync.WaitGroup
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

	// Output final throughput
	OutputAverageThroughput()
	logger.Finish()
}

// FIXME: The two following functions are very repetitive. Find a way to clean it
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
		go func() {
			executeHttpRequests(threadGroup, httpClient)
			threadWaitGroup.Done()
		}()
	}

	threadWaitGroup.Wait()
	stageWaitGroup.Done()
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