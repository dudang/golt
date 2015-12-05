package main

import (
	"os"
	"github.com/codegangsta/cli"
	"github.com/dudang/golt/parser"
	"github.com/dudang/golt/runner"
	"fmt"
)

var filename string
var version = "0.1"

func main() {
	app := cli.NewApp()
	app.Name = "golt"
	// TODO: Find a good description for the cli
	app.Usage = "Go Load Test Framework!!"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "file, f",
			Value: "golt.json",
			Usage: "full path to the load test file",
			Destination: &filename,
		},
	}

	app.Action = func(c *cli.Context) {
		golt, err := parser.ParseInputFile(filename)
		if err != nil {
			fmt.Println("Error occured during parsing of the file:")
			fmt.Printf("%v\n",err)
			os.Exit(1)
		}
		runner.ExecuteGoltTest(golt)
	}

	app.Run(os.Args)
}
