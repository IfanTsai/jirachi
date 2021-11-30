package interpreter_test

import (
	"testing"

	"github.com/IfanTsai/jirachi/common"
	"github.com/pkg/errors"

	"github.com/IfanTsai/jirachi/interpreter"
	lexer2 "github.com/IfanTsai/jirachi/lexer"
	"github.com/IfanTsai/jirachi/parser"
	"github.com/stretchr/testify/require"
)

func TestJInterpreter_Visit(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		text        string
		checkResult func(t *testing.T, number *interpreter.JNumber, err error)
	}{
		{
			name: "OK1",
			text: "(-1 + 2 ^ 3) ^ 2 * 13 / (24 - 5.8)",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0.0, number.Value)
				require.Equal(t, 35.0, number.Value.(float64))
			},
		},
		{
			name: "OK2",
			text: "---1",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.Value)
				require.Equal(t, -1, number.Value.(int))
			},
		},
		{
			name: "OK3",
			text: "----1",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.Value)
				require.Equal(t, 1, number.Value.(int))
			},
		},
		{
			name: "Division by zero integer number",
			text: "13 / 0",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.Error(t, err)
				require.IsType(t, &common.JRunTimeError{}, errors.Cause(err))
				require.Contains(t, err.Error(), "Runtime Error: Division by zero")
				require.Nil(t, number)
			},
		},
		{
			name: "Division by zero float number",
			text: "13 / 0.0",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.Error(t, err)
				require.IsType(t, &common.JRunTimeError{}, errors.Cause(err))
				require.Contains(t, err.Error(), "Runtime Error: Division by zero")
				require.Nil(t, number)
			},
		},
		{
			name: "Division by zero expression",
			text: "13 / (5 - 5)",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.Error(t, err)
				require.IsType(t, &common.JRunTimeError{}, errors.Cause(err))
				require.Contains(t, err.Error(), "Runtime Error: Division by zero")
				require.Nil(t, number)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			tokens, err := lexer2.NewJLexer("stdin", testCase.text).MakeTokens()
			require.NoError(t, err)
			require.NotEmpty(t, tokens)

			ast, err := parser.NewJParser(tokens, -1).Parse()
			require.NoError(t, err)
			require.NotNil(t, ast)

			number, err := interpreter.NewJInterpreter().Visit(ast)
			testCase.checkResult(t, number, err)
		})
	}
}
