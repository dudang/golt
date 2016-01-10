package main
import (
	"sync"
	"sort"
	"time"
	"net/http"
)

const parallelGroup = "parallel"

var stageWaitGroup sync.WaitGroup
var threadWaitGroup sync.WaitGroup
var channel = make(chan []byte, 1024)
var httpClient *http.Client

type HttpSender struct {
	Client *http.Client
}

func (http HttpSender) Send(request *http.Request) (*http.Response, error) {
	return http.Client.Do(request)
}

func ExecuteGoltTest(goltTest Golts, logFile string) {
	m := generateGoltMap(goltTest)

	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	SetOutputFile(logFile)
	go Watch(channel)
	for _, k := range keys {
		executeStage(m[k])
	}

	// Output final throughput
	OutputAverageThroughput()
	Finish()
}

func generateGoltMap(goltTest Golts) map[int][]GoltThreadGroup {
	m := make(map[int][]GoltThreadGroup)
	for _, element := range goltTest.Golt {
		array := m[element.Stage]
		if len(array) == 0 {
			m[element.Stage] = []GoltThreadGroup{element}
		} else {
			m[element.Stage] = append(array, element)

		}
	}
	return m
}

// FIXME: The two following functions are very repetitive. Find a way to clean it
func executeStage(stage []GoltThreadGroup) {
	stageWaitGroup.Add(len(stage))

	for _, item := range stage {
		httpClient = generateHttpClient(item)
		go executeThreadGroup(item)
	}
	stageWaitGroup.Wait()
}

func executeThreadGroup(threadGroup GoltThreadGroup) {
	threadWaitGroup.Add(threadGroup.Threads)

	for i := 0; i < threadGroup.Threads; i++ {
		go func() {
			executeHttpRequests(threadGroup, HttpSender{httpClient})
			threadWaitGroup.Done()
		}()
	}

	threadWaitGroup.Wait()
	stageWaitGroup.Done()
}

func generateHttpClient(threadGroup GoltThreadGroup) *http.Client {
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