package logger

import (
	"fmt"
	"os"
	"log"
)

var logFile *os.File

func init() {
	logFile, err := os.OpenFile("golt.log", os.O_TRUNC | os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.SetOutput(logFile)
}

func Log(event []byte) {
	log.Printf("%s\n", event)
}

func Finish() {
	logFile.Close()
}