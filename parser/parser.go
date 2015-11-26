package parser

import (
	"fmt"
	"os"
)

func ParseInputFile(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", filename)
		return
	}
	fmt.Printf("file exists: %s", filename)
}
