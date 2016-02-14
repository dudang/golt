package main

import (
	"testing"
	"os"
	"time"
	"bufio"
	"strings"
)

const logFile = "init-test.log";
const testFile = "unit-test.log";
const wrongFile = "\\ ";

func TestInit(t *testing.T) {
	logger := &FileLogger{Filename: logFile, }
	logger.Init()
	if _, err := os.Stat(logFile); err != nil {
		t.Error("The logger was not initialized properly")
	} else {
		logger.Finish()
	}
	loggerError := &FileLogger{Filename: wrongFile, }
	loggerError.Init()
	if loggerError == nil {
		t.Error("Was expecting an error opening the wrong file")
	}
}

func TestLog(t *testing.T) {
	// Prepare the file to be logged in
	msg := LogMessage{
		Url: "http://test.com",
		ErrorMessage: "N/A",
		Status: 200,
		Success: true,
		Duration: time.Since(time.Now()),
	}
	logger := &FileLogger{
		Filename: testFile,
	}
	logger.Init()

	// Log the message
	logger.Log(msg)

	// Open the file to go read the message
	file, _ := os.Open(testFile)
	reader := bufio.NewReader(file)

	// Skip the header line
	reader.ReadString('\n')

	// Read the expected log message
	loggedMessage, err := reader.ReadString('\n')
	if err != nil {
		t.Error("Could not read the line expected")
	}
	if !strings.Contains(loggedMessage, "http://test.com") {
		t.Error("The logged message was not read properly")
	}

	// Close the file
	logger.Finish()
}