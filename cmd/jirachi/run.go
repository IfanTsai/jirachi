package main

import (
	"os"

	"github.com/IfanTsai/jirachi/repl"
)

func main() {
	if len(os.Args) == 1 {
		repl.Run()
	}
	// TODO: execute source file
}
