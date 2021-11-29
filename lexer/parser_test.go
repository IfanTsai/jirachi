package lexer_test

import (
	"testing"

	"github.com/IfanTsai/jirachi/lexer"
	"github.com/stretchr/testify/require"
)

func TestJParser_Parse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		text        string
		checkResult func(t *testing.T, tokens *lexer.JNode, err error)
	}{
		{
			name: "OK",
			text: "(-1 + 2) * 13 / 24 - 5.8",
			checkResult: func(t *testing.T, node *lexer.JNode, err error) {
				t.Helper()
				require.NoError(t, err)
				require.NotEmpty(t, node, err)

				resStr := "(((((MINUS INT:1) PLUS INT:2) MUL INT:13) DIV INT:24) MINUS FLOAT:5.8)"
				require.Equal(t, node.String(), resStr)
			},
		},
		{
			name: "Invalid Syntax: Expected '+', '-', '*' or '/'",
			text: "1 + ",
			checkResult: func(t *testing.T, node *lexer.JNode, err error) {
				t.Helper()
				require.Error(t, err)
				require.Contains(t, err.Error(), "Invalid Syntax: Expected '+', '-', '*' or '/'")
				require.Nil(t, node)
			},
		},
		{
			name: "Invalid Syntax, End of token",
			text: "1 + *2",
			checkResult: func(t *testing.T, node *lexer.JNode, err error) {
				t.Helper()
				require.Error(t, err)
				require.Contains(t, err.Error(), "Invalid Syntax: Expected '+', '-', '*' or '/'")
				require.Nil(t, node)
			},
		},
		{
			name: "Invalid Syntax: Expected ')'",
			text: "1 * (-2 * (3 / (2 * 5))",
			checkResult: func(t *testing.T, node *lexer.JNode, err error) {
				t.Helper()
				require.Error(t, err)
				require.Contains(t, err.Error(), "Invalid Syntax: Expected ')'")
				require.Nil(t, node)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			tokens, err := lexer.NewJLexer("stdin", testCase.text).MakeTokens()
			require.NoError(t, err)
			require.NotEmpty(t, tokens)

			node, err := lexer.NewJParser(tokens, -1).Parse()
			testCase.checkResult(t, node, err)
		})
	}
}
