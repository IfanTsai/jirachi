package repl

import (
	"fmt"
	"strings"

	"github.com/IfanTsai/jirachi/token"

	"github.com/pkg/errors"

	"github.com/IfanTsai/jirachi/interpreter"
	"github.com/chzyer/readline"
)

var Release string

var completer = readline.NewPrefixCompleter(
	readline.PcItem(token.AND),
	readline.PcItem(token.OR),
	readline.PcItem(token.NOT),
	readline.PcItem(token.IF),
	readline.PcItem(token.THEN),
	readline.PcItem(token.ELIF),
	readline.PcItem(token.ELSE),
	readline.PcItem(token.END),
	readline.PcItem(token.FOR,
		readline.PcItem("i = 0", readline.PcItem(token.TO+" 10",
			readline.PcItem(token.STEP+" 1",
				readline.PcItem(token.THEN),
			),
			readline.PcItem(token.THEN)),
		),
	),
	readline.PcItem(token.WHILE),
	readline.PcItem(token.FUN),
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

		line = strings.TrimSpace(line)

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
