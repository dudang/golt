package parser

import (
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
)

type Golts struct {
	Golt []GoltItem
}

type GoltItem struct {
	URL string
	Method string
	Payload string
	Threads int
	Repetitions int
	Duration int
	Stage int
	Assert GoltAssert
}

type GoltAssert struct {
	Timeout int
	Status int
	Headers struct {}
	Body string
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