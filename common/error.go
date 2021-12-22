package common

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

type JExpectedCharacterError struct {
	*JError
	ExpectedChar byte
}

func (e *JExpectedCharacterError) Error() string {
	return e.ErrorString("Expected Character", "'"+string(e.ExpectedChar)+"'")
}

type JInvalidSyntaxError struct {
	*JError
	Details string
}

func (e *JInvalidSyntaxError) Error() string {
	return e.ErrorString("Invalid Syntax", e.Details)
}

type JRunTimeError struct {
	*JError
	Context *JContext
	Details string
}

func (e *JRunTimeError) Error() string {
	return e.generateTraceBack() + e.ErrorString("Runtime Error", e.Details)
}

func (e *JRunTimeError) generateTraceBack() string {
	pos := e.StartPos
	context := e.Context
	result := ""

	for context != nil {
		result = fmt.Sprintf("  File %s, line %d, in %s\n", pos.Filename, pos.Ln+1, context.Name) + result
		pos = context.ParentEntryPos
		context = context.Parent
	}

	return "Traceback (most recent call last):\n" + result
}

type JNumberTypeError struct {
	*JError
	Number interface{}
}

func (e *JNumberTypeError) Error() string {
	return e.ErrorString("Illegal Number Type", fmt.Sprintf("'%v'", e.Number))
}

func stringWithArrows(text string, startPos, endPos *JPosition) string {
	var strBuilder strings.Builder

	// calculate indices
	indexStart := strings.LastIndexByte(text[0:startPos.Index], '\n')
	if indexStart == -1 {
		indexStart = 0
	}

	indexEnd := len(text)
	if indexStart+1 <= len(text) {
		if indexEnd = strings.IndexByte(text[indexStart+1:], '\n'); indexEnd == -1 {
			indexEnd = len(text)
		} else {
			indexEnd += indexStart + 1
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
