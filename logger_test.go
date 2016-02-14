package main

import (
	"testing"
	"os"
)

const logFile = "unit-test.log";

func TestInit(t *testing.T) {
	logger := &FileLogger{
		Filename: logFile,
	}
	logger.Init()
	if _, err := os.Stat(logFile); err != nil {
		t.Error("The logger was not initialized properly")
	} else {
		os.Remove(logFile)
	}

}