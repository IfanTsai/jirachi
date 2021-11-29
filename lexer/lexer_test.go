package lexer_test

import (
	"testing"

	"github.com/IfanTsai/jirachi/lexer"

	"github.com/stretchr/testify/require"
)

func TestJLexer_MakeTokens(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		text        string
		checkResult func(t *testing.T, tokens []*lexer.JToken, err error)
	}{
		{
			name: "OK",
			text: "(-1 + 2) * 13 / 24 - 5.8",
			checkResult: func(t *testing.T, tokens []*lexer.JToken, err error) {
				t.Helper()
				require.NoError(t, err)
				require.NotEmpty(t, tokens)

				resStr := []string{
					"LPAREN", "MINUS", "INT:1", "PLUS", "INT:2", "RPAREN", "MUL",
					"INT:13", "DIV", "INT:24", "MINUS", "FLOAT:5.8", "EOF",
				}
				require.Len(t, tokens, len(resStr))
				for index, token := range tokens {
					require.Equal(t, token.String(), resStr[index])
				}
			},
		},
		{
			name: "illegal character &",
			text: "1&",
			checkResult: func(t *testing.T, tokens []*lexer.JToken, err error) {
				t.Helper()
				require.Error(t, err)
				require.Contains(t, err.Error(), "Illegal Character")
				require.Empty(t, tokens)
			},
		},
		{
			name: "illegal character $",
			text: "1 + 3 * 5$",
			checkResult: func(t *testing.T, tokens []*lexer.JToken, err error) {
				t.Helper()
				require.Error(t, err)
				require.Contains(t, err.Error(), "Illegal Character")
				require.Empty(t, tokens)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			tokens, err := lexer.NewJLexer("stdin", testCase.text).MakeTokens()
			testCase.checkResult(t, tokens, err)
		})
	}
}
