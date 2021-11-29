package lexer

import (
	"fmt"
	"strings"
)

type JError struct {
	StartPos *JPosition
	EndPos   *JPosition
}

func (e *JError) ErrorString(name, details string) string {
	return fmt.Sprintf("%s: %s\nFile <%s>, line %d, col %d\n\n%s",
		name, details,
		e.StartPos.Filename, e.StartPos.Ln, e.StartPos.Col+1,
		stringWithArrows(e.StartPos.Text, e.StartPos, e.EndPos),
	)
}

type JIllegalCharacterError struct {
	*JError
	IllegalChar byte
}

func (e *JIllegalCharacterError) Error() string {
	return e.ErrorString("Illegal Character", "'"+string(e.IllegalChar)+"'")
}

type JInvalidSyntaxError struct {
	*JError
	Details string
}

func (e *JInvalidSyntaxError) Error() string {
	return e.ErrorString("Invalid Syntax", e.Details)
}

func stringWithArrows(text string, startPos, endPos *JPosition) string {
	var strBuilder strings.Builder

	// calculate indices
	indexStart := strings.LastIndexByte(text[0:startPos.Index], '\n')
	if indexStart == -1 {
		indexStart = 0
	}

	indexEnd := len(text)
	if startPos.Index+1 <= len(text) {
		if indexEnd = strings.IndexByte(text[startPos.Index+1:], '\n'); indexEnd == -1 {
			indexEnd = len(text)
		}
	}

	lineCount := endPos.Ln - startPos.Ln + 1
	var i int64
	for i = 0; i < lineCount; i++ {
		line := text[indexStart:indexEnd]
		colStart := 0
		if i == 0 {
			colStart = startPos.Col
		}

		colEnd := len(line) - 1
		if i == lineCount-1 {
			colEnd = endPos.Col
		}

		// append to result
		strBuilder.WriteString(line + "\n")

		for i := 0; i < colStart; i++ {
			strBuilder.WriteByte(' ')
		}

		for i := 0; i < colEnd-colStart; i++ {
			strBuilder.WriteByte('^')
		}

		// re-calculate indices
		indexStart = indexEnd
		if indexStart+1 > len(text) {
			indexEnd = len(text)
		} else {
			if indexEnd = strings.IndexByte(text[indexStart+1:], '\n'); indexEnd < 0 {
				indexEnd = len(text)
			}
		}
	}

	return strings.ReplaceAll(strBuilder.String(), "\t", "")
}
