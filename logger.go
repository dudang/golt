package main

import (
	"fmt"
	"os"
	"log"
	"time"
)

type LogMessage struct {
	Url          string
	ErrorMessage string
	Status       int
	Success      bool
	Duration     time.Duration
}

type GoltLogger interface {
	SetOutputFile(filename string)
	Log(message LogMessage)
	Finish()
}

type FileLogger struct{
	LogFile *os.File
	Logger *log.Logger
}

func (l FileLogger) SetOutputFile(filename string) {
	var err error
	l.LogFile, err = os.OpenFile(filename, os.O_TRUNC | os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	l.Logger = log.New(l.LogFile, "", 0)
	l.Logger.Println("url,statusCode,success,duration,errorMessage")
}

func (l FileLogger) Log(message LogMessage) {
	milliseconds := message.Duration.Nanoseconds() / int64(time.Millisecond)
	msg := fmt.Sprintf("%s,%d,%t,%d,%v",
		message.Url,
		message.Status,
		message.Success,
		milliseconds,
		message.ErrorMessage)
	l.Logger.Printf("%s\n", msg)
}

func (l FileLogger) Finish() {
	l.LogFile.Close()
}