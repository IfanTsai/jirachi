package repl

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/IfanTsai/jirachi/interpreter"
	"github.com/chzyer/readline"
)

var Release string

var completer = readline.NewPrefixCompleter(
	readline.PcItem("AND"),
	readline.PcItem("OR"),
	readline.PcItem("NOT"),
	readline.PcItem("IF"),
	readline.PcItem("THEN"),
	readline.PcItem("ELIF"),
	readline.PcItem("ELSE"),
	readline.PcItem("FOR",
		readline.PcItem("i = 0", readline.PcItem("TO 10",
			readline.PcItem("STEP 1",
				readline.PcItem("THEN"),
			),
			readline.PcItem("THEN")),
		),
	),
	readline.PcItem("WHILE"),
)

func Run() {
	reader, err := readline.NewEx(&readline.Config{
		Prompt:            "\u001B[33mjirachi\u001B[0m \033[32m»\033[0m ",
		HistoryFile:       "/tmp/.jirachi_repl.tmp",
		HistorySearchFold: true,
		AutoComplete:      completer,
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
			if Release == "true" {
				fmt.Printf("%v\n", errors.Cause(err))
			} else {
				fmt.Printf("%+v\n", err)
			}
		} else if res != nil {
			fmt.Printf("%v\n", res)
		}
	}
}
