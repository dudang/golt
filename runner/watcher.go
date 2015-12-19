package runner

import (
	"time"
	"fmt"
)

var requestSent int64
var start time.Time

// TODO: Calculate throughput during the test, not just at the end
func Watch(channel chan []byte) {
	requestSent = 0
	timeCounter := 5.0
	start = time.Now()
	var delta time.Duration
	for {
		event := <- channel
		if event != nil {
			requestSent++
		}
		delta = time.Since(start)
		if delta.Seconds() > timeCounter {
			fmt.Printf("Throughput in the last 5 seconds: %f R/S\n", (float64(requestSent)/delta.Seconds()))
			timeCounter += 5.0
		}
	}
}

func CalculateThroughput() float64 {
	spentTime := time.Since(start)
	throughput := float64(requestSent) / spentTime.Seconds()
	return throughput
}