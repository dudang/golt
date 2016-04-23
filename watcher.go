package main

import (
	"fmt"
	"time"
)

var start time.Time
var delta time.Duration

type GoltWatcher struct {
	Interval        float64
	WatchingChannel chan []byte
	requestSent     int64
	requestDelta    int64
}

func (w *GoltWatcher) Watch() {
	timeCounter := w.Interval
	start = time.Now()
	for {
		event := <-w.WatchingChannel
		if event != nil {
			w.requestSent++
			w.requestDelta++
		}
		delta = time.Since(start)
		if delta.Seconds() > timeCounter {
			throughput := float64(w.requestDelta) / w.Interval
			fmt.Printf("Throughput in the last %.2f seconds: %.2f R/S - # requests sent last %.2f seconds %d\n",
				w.Interval,
				throughput,
				w.Interval,
				w.requestDelta)
			w.requestDelta = 0
			timeCounter += w.Interval
		}
	}
	w.OutputAverageThroughput()
}

func (w *GoltWatcher) OutputAverageThroughput() {
	spentTime := time.Since(start)
	throughput := float64(w.requestSent) / spentTime.Seconds()
	fmt.Printf("Average throughput: %f R/S - # requests sent %d\n", throughput, w.requestSent)
}
