package main

import (
	"os"
	"fmt"
	"github.com/codegangsta/cli"
)

var testFilename string
var logFilename string
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
			Destination: &testFilename,
		},
		cli.StringFlag{
			Name: "log, l",
			Value: "golt.log",
			Usage: "full path the the log file",
			Destination: &logFilename,
		},
	}

	app.Action = func(c *cli.Context) {
		fmt.Println("Started Golt")
		fmt.Println("Parsing input file...")
		golt, err := ParseInputFile(testFilename)
		if err != nil {
			fmt.Println("Error occured during parsing of the file:")
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Parsing completed")
		fmt.Println("Executing test...")
		ExecuteGoltTest(golt, logFilename)
		fmt.Println("Test completed, see results in the log file")
	}

	app.Run(os.Args)
}
