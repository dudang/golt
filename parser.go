package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// A Golts contains all the GoltThreadGroup generated from a configuration file.
type Golts struct {
	Golt []GoltThreadGroup
}

// A GoltThreadGroup contains the configuration of a single thread generated
// from a configuration file.
type GoltThreadGroup struct {
	Threads     int
	Repetitions int
	Stage       int
	Type        string
	Requests    []GoltRequest
}

// A GoltRequest contains the configuration of a single HTTP request.
type GoltRequest struct {
	URL     string
	Method  string
	Payload string
	Headers map[string]*string
	Assert  GoltAssert
	// TODO: Have the possibility to extract multiple values
	Extract GoltExtractor
}

// A GoltAssert contains the configuration of the assertions to be made for a
// GoltRequest.
type GoltAssert struct {
	Timeout int
	Status  int
	Type    string
}

// A GoltExtractor contains the configuration to extract information of the
// response of a GoltRequest.
type GoltExtractor struct {
	Var   string
	Field string
	Regex string
	// TODO: Have the possibility to extract the value of a JSON field from the headers/body
}

type convert func([]byte, interface{}) error

// ParseInputFile converts a supported configuration file to a Golts object.
func ParseInputFile(filename string) (Golts, error) {
	convertFn := extractConvertFunction(filename)
	if convertFn == nil {
		return Golts{}, errors.New("Unsupported file format, exiting")
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return Golts{}, err
	}
	return convertToStruct(convertFn, file)
}

func extractConvertFunction(filename string) convert {
	extension := extractExtension(filename)
	switch extension {
	case ".json":
		return json.Unmarshal
	case ".yaml", ".yml":
		return yaml.Unmarshal
	default:
		return nil
	}
}

func extractExtension(filename string) string {
	extension := filepath.Ext(filename)
	return strings.ToLower(extension)
}

func convertToStruct(convertFunction convert, file []byte) (Golts, error) {
	var golt Golts
	err := convertFunction(file, &golt)
	return golt, err
}
