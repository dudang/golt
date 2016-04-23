package main

import (
	"net/http"
	"sort"
	"sync"
	"time"
)

// Multithreading variables
var stageWaitGroup sync.WaitGroup
var threadWaitGroup sync.WaitGroup

// Communication channel to calculate throughput (request per second)
var channel = make(chan []byte, 1024)

var httpClient *http.Client
var logger GoltLogger
var watcher *GoltWatcher

func init() {
	// TODO: Make the interval parameterizable
	watcher = &GoltWatcher{
		Interval:        5.0,
		WatchingChannel: channel,
	}
}

func ExecuteGoltTest(goltTest Golts, logFile string) {
	m := generateStageMap(goltTest)

	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	logger = &FileLogger{
		Filename: logFile,
	}
	logger.Init()
	go watcher.Watch()
	for _, k := range keys {
		executeStage(m[k])
	}

	// Output final throughput
	watcher.OutputAverageThroughput()
	logger.Finish()
}

func generateStageMap(goltTest Golts) map[int][]GoltThreadGroup {
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

	executor := GoltExecutor{
		ThreadGroup:    threadGroup,
		Sender:         HttpSender{httpClient},
		Logger:         logger,
		SendingChannel: channel,
	}

	for i := 0; i < threadGroup.Threads; i++ {
		go func() {
			executor.ExecuteHttpRequests()
			threadWaitGroup.Done()
		}()
	}

	threadWaitGroup.Wait()
	stageWaitGroup.Done()
}

func generateHttpClient(threadGroup GoltThreadGroup) *http.Client {
	var httpClient *http.Client
	if threadGroup.Timeout > 0 {
		httpClient = &http.Client{
			Timeout: time.Duration(time.Millisecond * time.Duration(threadGroup.Timeout)),
		}
	} else {
		// Default timeout of 30 seconds for HTTP calls to avoid hung threads
		httpClient = &http.Client{
			Timeout: time.Duration(time.Second * 30),
		}
	}
	return httpClient
}
