package logger

import (
	"fmt"
	"os"
	"log"
	"time"
)

var logFile *os.File

type LogMessage struct {
	Stage        int
	Repetition   int
	ErrorMessage string
	Status       int
	Success      bool
	Duration     time.Duration
}

var logger *log.Logger

func SetOutputFile(filename string) {
	logFile, err := os.OpenFile(filename, os.O_TRUNC | os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	logger = log.New(logFile, "", 0)
	logger.Println("Stage,Repetitions,Status Code,Success,Duration,Error Message")
}

func Log(message LogMessage) {
	milliseconds := message.Duration.Nanoseconds() / int64(time.Millisecond)
	msg := fmt.Sprintf("%d,%d,%d,%t,%d,%v",
		message.Stage,
		message.Repetition,
		message.Status,
		message.Success,
		milliseconds,
		message.ErrorMessage)
	logger.Printf("%s\n", msg)
}

func Finish() {
	logFile.Close()
}