package main

import "testing"

var inputTable = []struct {
	in           string
	expectError  bool
	errorMessage string
}{
	{"test/golt-test.json", false, "Failed to parse the JSON file"},
	{"test/golt-test.yaml", false, "Failed to parse the YAML file"},
	{"test/golt-test.yml", false, "Failed to parse the YML file"},
	{"test/golt-test.txt", true, "Failed to detect files of the supported formats"},
	{"non_existing_file.json", true, "Failed to detect non-existing files"},
}

func TestParseInputFile(t *testing.T) {
	for _, entry := range inputTable {
		_, error := ParseInputFile(entry.in)
		hasError := error != nil
		if hasError != entry.expectError {
			t.Error(entry.errorMessage)
		}
	}
}
