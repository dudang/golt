package main

import (
	"testing"
)

func TestParseInputFile(t *testing.T) {
	_, err := ParseInputFile("test/golt-test.json")
	if err != nil {
		t.Fail()
	}

	_, err = ParseInputFile("wrong_file.json")
	if err == nil {
		t.Fail()
	}
}