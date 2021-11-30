package parser_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/IfanTsai/jirachi/common"

	lexer2 "github.com/IfanTsai/jirachi/lexer"
	"github.com/IfanTsai/jirachi/parser"

	"github.com/stretchr/testify/require"
)

func TestJParser_Parse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		text        string
		checkResult func(t *testing.T, tokens *parser.JNode, err error)
	}{
		{
			name: "OK",
			text: "(-1 + 2) * 13 / 24 - 5.8",
			checkResult: func(t *testing.T, node *parser.JNode, err error) {
				t.Helper()
				require.NoError(t, err)
				require.NotEmpty(t, node, err)

				resStr := "(((((MINUS INT:1) PLUS INT:2) MUL INT:13) DIV INT:24) MINUS FLOAT:5.8)"
				require.Equal(t, resStr, node.String())
			},
		},
		{
			name: "Invalid Syntax: Expected '+', '-', '*' or '/'",
			text: "1 + ",
			checkResult: func(t *testing.T, node *parser.JNode, err error) {
				t.Helper()
				require.Error(t, err)
				require.IsType(t, &common.JInvalidSyntaxError{}, errors.Cause(err))
				require.Contains(t, err.Error(), "Invalid Syntax: Expected int, float, '+', '-' or '('")
				require.Nil(t, node)
			},
		},
		{
			name: "Invalid Syntax, End of token",
			text: "1 + *2",
			checkResult: func(t *testing.T, node *parser.JNode, err error) {
				t.Helper()
				require.Error(t, err)
				require.IsType(t, &common.JInvalidSyntaxError{}, errors.Cause(err))
				require.Contains(t, err.Error(), "Invalid Syntax: Expected int, float, '+', '-' or '('")
				require.Nil(t, node)
			},
		},
		{
			name: "Invalid Syntax: Expected ')'",
			text: "1 * (-2 * (3 / (2 * 5))",
			checkResult: func(t *testing.T, node *parser.JNode, err error) {
				t.Helper()
				require.Error(t, err)
				require.IsType(t, &common.JInvalidSyntaxError{}, errors.Cause(err))
				require.Contains(t, err.Error(), "Invalid Syntax: Expected ')'")
				require.Nil(t, node)
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

			node, err := parser.NewJParser(tokens, -1).Parse()
			testCase.checkResult(t, node, err)
		})
	}
}
