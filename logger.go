package main

import (
	"fmt"
	"os"
	"log"
	"time"
)

type GoltLogger struct{
	LogFile *os.File
	Logger *log.Logger
}

type LogMessage struct {
	Url          string
	ErrorMessage string
	Status       int
	Success      bool
	Duration     time.Duration
}

func (l *GoltLogger) SetOutputFile(filename string) {
	var err error
	l.LogFile, err = os.OpenFile(filename, os.O_TRUNC | os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	l.Logger = log.New(l.LogFile, "", 0)
	l.Logger.Println("url,statusCode,success,duration,errorMessage")
}

func (l *GoltLogger) Log(message LogMessage) {
	milliseconds := message.Duration.Nanoseconds() / int64(time.Millisecond)
	msg := fmt.Sprintf("%s,%d,%t,%d,%v",
		message.Url,
		message.Status,
		message.Success,
		milliseconds,
		message.ErrorMessage)
	l.Logger.Printf("%s\n", msg)
}

func (l *GoltLogger) Finish() {
	l.LogFile.Close()
}