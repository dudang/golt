package runner

import (
	"time"
	"fmt"
)

var requestSent int64
var requestDelta int64
var start time.Time
var delta time.Duration

// TODO: Parametrize this
var throughputInterval = 5.0

func Watch(channel chan []byte) {
	requestSent = 0
	requestDelta = 0
	timeCounter := 5.0
	start = time.Now()
	for {
		event := <- channel
		if event != nil {
			requestSent++
			requestDelta++
		}
		delta = time.Since(start)
		if delta.Seconds() > timeCounter {
			throughput := float64(requestDelta)/throughputInterval
			fmt.Printf("Throughput in the last %.2f seconds: %.2f R/S - # requests sent last .2f seconds %d\n",
				throughputInterval,
				throughput,
				throughputInterval,
				requestDelta)
			requestDelta = 0
			timeCounter += throughputInterval
		}
	}
	OutputAverageThroughput()
}

func OutputAverageThroughput() {
	spentTime := time.Since(start)
	throughput := float64(requestSent) / spentTime.Seconds()
	fmt.Printf("Average throughput: %f R/S - # requests sent %d\n",throughput, requestSent)
}