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
	Init() error
	Log(message LogMessage)
	Finish()
}

type FileLogger struct{
	Filename string
	LogFile *os.File
	Logger *log.Logger
}

func (logger *FileLogger) Init() error {
	var err error
	logger.LogFile, err = os.OpenFile(logger.Filename, os.O_TRUNC | os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return err
	}
	logger.Logger = log.New(logger.LogFile, "", 0)
	logger.Logger.Println("url,statusCode,success,duration,errorMessage")
	return nil
}

func (logger *FileLogger) Log(message LogMessage) {
	milliseconds := message.Duration.Nanoseconds() / int64(time.Millisecond)
	msg := fmt.Sprintf("%s,%d,%t,%d,%v",
		message.Url,
		message.Status,
		message.Success,
		milliseconds,
		message.ErrorMessage)
	logger.Logger.Printf("%s\n", msg)
}

func (logger *FileLogger) Finish() {
	logger.LogFile.Close()
}