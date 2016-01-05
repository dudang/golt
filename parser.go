package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

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
