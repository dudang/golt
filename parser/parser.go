package parser

import (
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
)

type Golts struct {
	Golt []GoltThreadGroup
}

type GoltThreadGroup struct {
	Threads     int
	Repetitions int
	Stage       int
	Type        string
	Requests    []GoltRequest
}

type GoltRequest struct {
	URL     string
	Method  string
	Payload string
	Headers map[string] *string
	Assert  GoltAssert
	Extract GoltExtractor
}

type GoltAssert struct {
	Timeout int
	Status  int
	Type    string
}

type GoltExtractor struct {
	Var   string
	Field string
	Regex string
}

type convert func([]byte, interface{}) error

func ParseInputFile(filename string) (Golts, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return Golts{}, err
	}

	switch filepath.Ext(filename) {
	case ".json":
		return convertToStruct(json.Unmarshal, file)
	case ".yaml":
		return convertToStruct(yaml.Unmarshal, file)
	default:
		return Golts{}, errors.New("Unknown file type, exiting")
	}
}

func convertToStruct(convertFunction convert, file []byte) (Golts, error) {
	var golt Golts
	err := convertFunction(file, &golt)
	return golt, err
}