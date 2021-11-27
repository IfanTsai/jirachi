package lexer

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const DIGITS = "0123456789"

type IllegalCharacterError struct {
	*JPosition
	IllegalChar byte
}

func (e *IllegalCharacterError) Error() string {
	return fmt.Sprintf("illegal character '%c'\nFile <%s>, line %d, col %d",
		e.IllegalChar, e.Filename, e.Ln, e.Col)
}

type JLexer struct {
	Text []byte
	Pos  *JPosition
}

func Run(filename, text string) (*JNode, error) {
	tokens, err := NewJLexer(filename, text).makeTokens()
	if err != nil {
		return nil, err
	}

	return NewJParser(tokens, -1).Parse(), nil
}

func NewJLexer(filename, text string) *JLexer {
	return &JLexer{
		Text: []byte(text),
		Pos:  NewJPosition(-1, 0, 0, filename),
	}
}

func (l *JLexer) makeTokens() ([]*JToken, error) {
	tokens := make([]*JToken, 0, len(l.Text))

	for advanceAble := l.advance(); advanceAble; {
		char := l.Text[l.Pos.Index]
		switch {
		case char == ' ' || char == '\t':
			advanceAble = l.advance()
		case strings.IndexByte(DIGITS, char) != -1:
			var token *JToken
			token, advanceAble = l.makeNumberToken()
			tokens = append(tokens, token)
		case char == '+':
			tokens = append(tokens, NewJToken(PLUS, ""))
			advanceAble = l.advance()
		case char == '-':
			tokens = append(tokens, NewJToken(MINUS, ""))
			advanceAble = l.advance()
		case char == '*':
			tokens = append(tokens, NewJToken(MUL, ""))
			advanceAble = l.advance()
		case char == '/':
			tokens = append(tokens, NewJToken(DIV, ""))
			advanceAble = l.advance()
		case char == '(':
			tokens = append(tokens, NewJToken(LPAREN, ""))
			advanceAble = l.advance()
		case char == ')':
			tokens = append(tokens, NewJToken(RPAREN, ""))
			advanceAble = l.advance()
		default:
			return nil, errors.Wrap(&IllegalCharacterError{
				IllegalChar: char,
				JPosition:   l.Pos,
			}, "failed to parse token")
		}
	}

	return tokens, nil
}

func (l *JLexer) makeNumberToken() (*JToken, bool) {
	advanceAble := true
	isFloat := false
	numStr := ""

	for {
		char := l.Text[l.Pos.Index]

		if strings.IndexByte(DIGITS+".", char) == -1 {

			break
		}

		if char == '.' {
			if isFloat {
				break
			}

			isFloat = true
		}

		numStr += string(char)

		if !l.advance() {
			advanceAble = false

			break
		}
	}

	if isFloat {
		return NewJToken(FLOAT, numStr), advanceAble
	}

	return NewJToken(INT, numStr), advanceAble
}

func (l *JLexer) advance() bool {
	if l.Pos.Index+1 >= len(l.Text) {
		return false
	}

	l.Pos.Advance(l.Text)

	return true
}
