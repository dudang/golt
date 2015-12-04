package parser

import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
)

func ParseInputFile(filename string) {
	switch filepath.Ext(filename) {
		case ".json":
			fmt.Println("We're dealing with JSON!")
		case ".yaml":
			fmt.Println("We're dealing with YAML!")
		default:
			fmt.Println("Unknown file type, exiting")
			os.Exit(1)
	}

	file, e := ioutil.ReadFile(filename)

	if e != nil {
		fmt.Printf("File error, %v\n", e)
		os.Exit(1)
	}

	fmt.Printf("%s\n",file)
}