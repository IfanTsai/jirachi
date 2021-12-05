package repl

import (
	"errors"
	"fmt"

	"github.com/IfanTsai/jirachi/interpreter"
	"github.com/chzyer/readline"
)

func Run() {
	reader, err := readline.NewEx(&readline.Config{
		Prompt:            "\u001B[33mjirachi\u001B[0m \033[32m»\033[0m ",
		HistoryFile:       "/tmp/.jirachi_repl.tmp",
		HistorySearchFold: true,
	})
	if err != nil {
		panic(err)
	}

	defer reader.Close()

	for {
		line, err := reader.Readline()
		if err != nil {
			if !errors.Is(err, readline.ErrInterrupt) {
				fmt.Println(err)
			}

			break
		}

		if len(line) == 0 {
			continue
		}

		res, err := interpreter.Run("stdin", line)
		if err != nil {
			fmt.Printf("%+v\n", err)
		} else if res != nil {
			fmt.Printf("%v\n", res)
		}
	}
}
