package interpreter_test

import (
	"testing"

	"github.com/IfanTsai/jirachi/interpreter/object"

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
		checkResult func(t *testing.T, value object.JValue, err error)
	}{
		{
			name: "OK1",
			text: "(-1 + 2 ^ 3) ^ 2 * 13 / (24 - 5.8)",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0.0, number.GetValue())
				require.Equal(t, 35.0, number.GetValue().(float64))
			},
		},
		{
			name: "OK2",
			text: "---1",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, -1, number.GetValue().(int))
			},
		},
		{
			name: "OK3",
			text: "----1",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, 1, number.GetValue().(int))
			},
		},
		{
			name: "OK4: variable assign",
			text: "1 + (a = 2)",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, 3, number.GetValue().(int))

				varNumber := interpreter.GlobalSymbolTable.Get("a")
				require.IsType(t, &object.JNumber{}, varNumber)
				require.IsType(t, 0, varNumber.(*object.JNumber).Value)
				require.Equal(t, 2, varNumber.(*object.JNumber).Value)
			},
		},
		{
			name: "OK5: variable access",
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
				number, err := interpreter.NewJInterpreter(context).Interpreter(ast)
				require.NoError(t, err)
				require.NotNil(t, number)
			},
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, 8, number.GetValue().(int))

				varNumber := interpreter.GlobalSymbolTable.Get("b")
				require.IsType(t, &object.JNumber{}, varNumber)
				require.IsType(t, 0, varNumber.(*object.JNumber).Value)
				require.Equal(t, 5, varNumber.(*object.JNumber).Value)
			},
		},
		{
			name: "OK6: logical operation",
			text: "5 - 5 OR 1 + 2 AND (NOT 0 AND 10) - 2 * 5",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, 0, number.GetValue().(int))
			},
		},
		{
			name: "OK7: logical operation",
			text: "5 - 5 OR 1 + 2 AND NOT 0 AND 10 - 2 * 5",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, 0, number.GetValue().(int))
			},
		},
		{
			name: "OK8: comparison operation",
			text: " (3 > 2) OR 1 + 2 AND NOT 0 AND 10 - 2 * 5 AND (5 - 5 == 0)",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, 0, number.GetValue().(int))
			},
		},
		{
			name: "OK9: if expression",
			text: "IF 5 > 3 THEN 4 ELSE 5",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, 4, number.GetValue().(int))
			},
		},
		{
			name: "OK10: if expression",
			text: "IF 5 < 3 THEN 4 ELSE 5",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, 5, number.GetValue().(int))
			},
		},
		{
			name: "OK11: if expression",
			text: "IF 5 > 6 THEN 4 ELIF 5 > 4 THEN 6 ELSE 5",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, 6, number.GetValue().(int))
			},
		},
		{
			name: "OK12: if expression",
			text: "IF 5 == 6 THEN 4",
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, nil, number.GetValue())
				require.Equal(t, nil, number.GetValue())
			},
		},
		{
			name: "OK13: for expression, default step",
			text: "FOR i = 1 TO 6 THEN res13 = res13 * i",
			preRun: func(t *testing.T) {
				t.Helper()

				assignVariable(t, "res13 = 1")
			},
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJList(nil), value)

				elementValues := value.(*object.JList).ElementValues
				resValue := elementValues[len(elementValues)-1]
				require.IsType(t, 0, resValue.GetValue())
				require.Equal(t, 120, resValue.GetValue())
			},
		},
		{
			name: "OK14: for expression, step = -1",
			text: "FOR i = 5 TO 0 STEP -1 THEN res14 = res14 * i",
			preRun: func(t *testing.T) {
				t.Helper()

				assignVariable(t, "res14 = 1")
			},
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJList(nil), value)

				elementValues := value.(*object.JList).ElementValues
				resValue := elementValues[len(elementValues)-1]
				require.IsType(t, 0, resValue.GetValue())
				require.Equal(t, 120, resValue.GetValue())
			},
		},
		{
			name: "OK15: while expression",
			text: "WHILE res15 < 10000 THEN res15 = res15 + 1",
			preRun: func(t *testing.T) {
				t.Helper()

				assignVariable(t, "res15 = 0")
			},
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJList(nil), value)

				elementValues := value.(*object.JList).ElementValues
				resValue := elementValues[len(elementValues)-1]
				require.IsType(t, 0, resValue.GetValue())
				require.Equal(t, 10000, resValue.GetValue())
			},
		},
		{
			name: "OK16: function",
			text: "add(1, 2)",
			preRun: func(t *testing.T) {
				t.Helper()

				assignVariable(t, "add = FUN(a, b) -> a + b")
			},
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, number.GetValue())
				require.Equal(t, 3, number.GetValue())
			},
		},
		{
			name: "OK17: string mul",
			text: "str17 * 4",
			preRun: func(t *testing.T) {
				t.Helper()

				assignVariable(t, "str17 = 'hello'")
			},
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, "", number.GetValue())
				require.Equal(t, "hellohellohellohello", number.GetValue())
			},
		},
		{
			name: "OK18: string add",
			text: "str18 + \" world\"",
			preRun: func(t *testing.T) {
				t.Helper()

				assignVariable(t, "str18 = 'hello'")
			},
			checkResult: func(t *testing.T, number object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, "", number.GetValue())
				require.Equal(t, "hello world", number.GetValue())
			},
		},
		{
			name: "OK19: list add list",
			text: "[1, 2, 3] + [4, 5, 6]",
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJList(nil), value)
				require.Equal(t, "[1, 2, 3, 4, 5, 6]", value.String())
			},
		},
		{
			name: "OK20: list add number",
			text: "[1, 2, 3] + 6",
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJList(nil), value)
				require.Equal(t, "[1, 2, 3, 6]", value.String())
			},
		},
		{
			name: "OK21: list remove number",
			text: "[1, 2, 3, 4] - 1",
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJList(nil), value)
				require.Equal(t, "[1, 3, 4]", value.String())
			},
		},
		{
			name: "OK22: list mul number",
			text: "[1, 2, 3] * 3",
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJList(nil), value)
				require.Equal(t, "[1, 2, 3, 1, 2, 3, 1, 2, 3]", value.String())
			},
		},
		{
			name: "OK23: list mul list",
			text: "[1, 2, 3] * [4, 5, 6]",
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, object.NewJList(nil), value)
				require.Equal(t, "[4, 10, 18]", value.String())
			},
		},
		{
			name: "OK24: list index access",
			text: "list24[3]",
			preRun: func(t *testing.T) {
				t.Helper()

				assignVariable(t, "list24 = [2, 3, 4, 5, 6]")
			},
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, value.GetValue())
				require.Equal(t, 5, value.GetValue())
			},
		},
		{
			name: "OK25: len list",
			text: "len(list25)",
			preRun: func(t *testing.T) {
				t.Helper()

				assignVariable(t, "list25 = [2, 3, 4, 5, 6, 7]")
			},
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, value.GetValue())
				require.Equal(t, 6, value.GetValue())
			},
		},
		{
			name: "OK26: len string",
			text: "len(string26)",
			preRun: func(t *testing.T) {
				t.Helper()

				assignVariable(t, "string26 = 'hello world'")
			},
			checkResult: func(t *testing.T, value object.JValue, err error) {
				t.Helper()
				require.NoError(t, err)
				require.IsType(t, 0, value.GetValue())
				require.Equal(t, 11, value.GetValue())
			},
		},
		{
			name: "Division by zero integer number",
			text: "13 / 0",
			checkResult: func(t *testing.T, number object.JValue, err error) {
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
			checkResult: func(t *testing.T, number object.JValue, err error) {
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
			checkResult: func(t *testing.T, number object.JValue, err error) {
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
				number, err := interpreter.NewJInterpreter(context).Interpreter(ast)
				require.NoError(t, err)
				require.NotNil(t, number)
			},
			checkResult: func(t *testing.T, number object.JValue, err error) {
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
			checkResult: func(t *testing.T, number object.JValue, err error) {
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
			number, err := interpreter.NewJInterpreter(context).Interpreter(ast)
			testCase.checkResult(t, number, err)
		})
	}
}

func assignVariable(t *testing.T, variableAssign string) {
	t.Helper()

	tokens, err := lexer.NewJLexer("stdin", variableAssign).MakeTokens()
	require.NoError(t, err)
	require.NotEmpty(t, tokens)

	ast, err := parser.NewJParser(tokens, -1).Parse()
	require.NoError(t, err)
	require.NotNil(t, ast)

	context := common.NewJContext("test", interpreter.GlobalSymbolTable, nil, nil)
	number, err := interpreter.NewJInterpreter(context).Interpreter(ast)
	require.NoError(t, err)
	require.NotNil(t, number)
}
