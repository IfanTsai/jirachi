package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/IfanTsai/jirachi/repl"

	"github.com/IfanTsai/jirachi/interpreter"
)

func main() {
	if len(os.Args) == 1 {
		repl.Run()
	} else {
		filename := os.Args[1]
		scriptFile, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		bytes, err := io.ReadAll(scriptFile)
		if err != nil {
			panic(err)
		}

		_, err = interpreter.Run(filename, string(bytes))
		if err != nil {
			if repl.Release == "true" {
				fmt.Printf("%v", errors.Cause(err))
			} else {
				fmt.Printf("%+v", err)
			}
		}

	}
}
