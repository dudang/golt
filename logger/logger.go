package logger

import (
	"time"
	"fmt"
)

// TODO: Make this parameterizable
const workerCount = 4
var channel = make(chan []byte, 1024)
var workers = make([]*Worker, workerCount)

// TODO: We need to merge the log files after
func init() {
	for i := 0; i < workerCount; i++ {
		workers[i] = NewWorker(i)
		go workers[i].Work(channel)
	}
}

func Log(event []byte) {
	select {
		case channel <- event:
		case <- time.After(5 * time.Second):
			fmt.Println("Message hanged for more than 5 seconds, lost message")
			fmt.Printf("%s", event)
	}
}

func Flush() {
	for _, worker := range workers {
		worker.Save()
	}
}