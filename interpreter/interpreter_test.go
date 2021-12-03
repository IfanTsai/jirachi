package interpreter_test

import (
	"testing"

	"github.com/IfanTsai/jirachi/common"
	"github.com/pkg/errors"

	"github.com/IfanTsai/jirachi/interpreter"
	"github.com/IfanTsai/jirachi/lexer"
	"github.com/IfanTsai/jirachi/parser"
	"github.com/stretchr/testify/require"
)

func TestJInterpreter_Visit(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		text        string
		preRun      func(t *testing.T)
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
			name: "OK4",
			text: "1 + (a = 2)",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.Value)
				require.Equal(t, 3, number.Value.(int))

				varNumber := interpreter.GlobalSymbolTable.Get("a")
				require.IsType(t, &interpreter.JNumber{}, varNumber)
				require.IsType(t, 0, varNumber.(*interpreter.JNumber).Value)
				require.Equal(t, 2, varNumber.(*interpreter.JNumber).Value)
			},
		},
		{
			name: "OK5",
			text: "b + 3",
			preRun: func(t *testing.T) {
				t.Helper()

				tokens, err := lexer.NewJLexer("stdin", "b = 5").MakeTokens()
				require.NoError(t, err)
				require.NotEmpty(t, tokens)

				ast, err := parser.NewJParser(tokens, -1).Parse()
				require.NoError(t, err)
				require.NotNil(t, ast)

				context := common.NewJContext("test", interpreter.GlobalSymbolTable, nil, nil)
				number, err := interpreter.NewJInterpreter(context).Visit(ast)
				require.NoError(t, err)
				require.NotNil(t, number)
			},
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.Value)
				require.Equal(t, 8, number.Value.(int))

				varNumber := interpreter.GlobalSymbolTable.Get("b")
				require.IsType(t, &interpreter.JNumber{}, varNumber)
				require.IsType(t, 0, varNumber.(*interpreter.JNumber).Value)
				require.Equal(t, 5, varNumber.(*interpreter.JNumber).Value)
			},
		},
		{
			name: "OK6",
			text: "5 - 5 OR 1 + 2 AND (NOT 0 AND 10) - 2 * 5",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.Value)
				require.Equal(t, 1, number.Value.(int))
			},
		},
		{
			name: "OK7",
			text: "5 - 5 OR 1 + 2 AND NOT 0 AND 10 - 2 * 5",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.Value)
				require.Equal(t, 0, number.Value.(int))
			},
		},
		{
			name: "OK8",
			text: " (3 > 2) OR 1 + 2 AND NOT 0 AND 10 - 2 * 5 AND (5 - 5 == 0)",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.Value)
				require.Equal(t, 0, number.Value.(int))
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
		{
			name: "Division by zero variable",
			text: "13 / (c - 2)",
			preRun: func(t *testing.T) {
				t.Helper()

				tokens, err := lexer.NewJLexer("stdin", "c = 2").MakeTokens()
				require.NoError(t, err)
				require.NotEmpty(t, tokens)

				ast, err := parser.NewJParser(tokens, -1).Parse()
				require.NoError(t, err)
				require.NotNil(t, ast)

				context := common.NewJContext("test", interpreter.GlobalSymbolTable, nil, nil)
				number, err := interpreter.NewJInterpreter(context).Visit(ast)
				require.NoError(t, err)
				require.NotNil(t, number)
			},
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.Error(t, err)
				require.IsType(t, &common.JRunTimeError{}, errors.Cause(err))
				require.Contains(t, err.Error(), "Runtime Error: Division by zero")
				require.Nil(t, number)
			},
		},
		{
			name: "Variable is not defined",
			text: "13 / (abc - 2)",
			checkResult: func(t *testing.T, number *interpreter.JNumber, err error) {
				t.Helper()
				require.Error(t, err)
				require.IsType(t, &common.JRunTimeError{}, errors.Cause(err))
				require.Contains(t, err.Error(), "Runtime Error: 'abc' is not defined")
				require.Nil(t, number)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			if testCase.preRun != nil {
				testCase.preRun(t)
			}

			tokens, err := lexer.NewJLexer("stdin", testCase.text).MakeTokens()
			require.NoError(t, err)
			require.NotEmpty(t, tokens)

			ast, err := parser.NewJParser(tokens, -1).Parse()
			require.NoError(t, err)
			require.NotNil(t, ast)

			context := common.NewJContext("test", interpreter.GlobalSymbolTable, nil, nil)
			number, err := interpreter.NewJInterpreter(context).Visit(ast)
			testCase.checkResult(t, number, err)
		})
	}
}
