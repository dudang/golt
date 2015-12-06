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

func ParseInputFile(filename string) (Golts, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return Golts{}, err
	}

	switch filepath.Ext(filename) {
		case ".json":
			return jsonToStruct(file)
		case ".yaml":
			return yamlToStruct(file)
		default:
			return Golts{}, errors.New("Unknown file type, exiting")
	}
}

func jsonToStruct(file []byte) (Golts, error) {
	var golt Golts
	err := json.Unmarshal(file, &golt)
	return golt, err
}

func yamlToStruct(file []byte) (Golts, error) {
	var golt Golts
	err := yaml.Unmarshal(file, &golt)
	return golt, err
}