package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/IfanTsai/jirachi/lexer"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("jirachi > ")
		reader.Scan()
		node, err := lexer.Run("stdin", reader.Text())
		if err != nil {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Printf("%v\n", node)
		}
	}
}
