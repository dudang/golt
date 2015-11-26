package main

import (
	"flag"
	"github.com/dudang/golt/parser"
)

var filename string

func init() {
	flag.StringVar(&filename, "file", "golt.yaml", "full path to the load test file")
}

func main() {
	flag.Parse()
	parser.ParseInputFile(filename)
}
