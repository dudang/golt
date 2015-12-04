package parser

import (
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"errors"
)

type GoltJsons struct {
	Golt []GoltJson
}

type GoltJson struct {
	URL string `json:"url"`
	Method string `json:"method"`
	Payload string `json:"body"`
	Threads int `json:"threads"`
	Repetitions int `json:"repetitions"`
	Duration int `json:"duration"`
	Stage int `json:"stage"`
	Assert struct {
			Timeout int `json:"timeout"`
			Status int `json:"status"`
			Headers struct {
					} `json:"headers"`
			Body string `json:"body"`
		} `json:"assert"`
}

func ParseInputFile(filename string) (GoltJsons, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return GoltJsons{}, err
	}

	switch filepath.Ext(filename) {
		case ".json":
			golt, err := jsonToStruct(file)
			return golt, err
		case ".yaml":
			return GoltJsons{}, errors.New("We're dealing with YAML, but it's not yet implemented. Sorry !")
		default:
			return GoltJsons{}, errors.New("Unknown file type, exiting")
	}
}

func jsonToStruct(file []byte) (GoltJsons, error) {
	var golt GoltJsons
	err := json.Unmarshal(file, &golt)
	return golt, err
}