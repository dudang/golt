package runner

import (
	"time"
)

var requestSent int64
var start time.Time

// TODO: Calculate throughput during the test, not just at the end
func Watch(channel chan []byte) {
	requestSent = 0
	start = time.Now()
	for {
		event := <- channel
		if event != nil {
			requestSent++
		}
	}
}

func CalculateThroughput() float64 {
	spentTime := time.Since(start)
	throughput := float64(requestSent) / spentTime.Seconds()
	return throughput
}