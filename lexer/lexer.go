package lexer

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const DIGITS = "0123456789"

func Run(filename, text string) (*JNode, error) {
	// generate tokens
	tokens, err := NewJLexer(filename, text).MakeTokens()
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", tokens)

	// generate AST
	return NewJParser(tokens, -1).Parse()
}

type JLexer struct {
	Text []byte
	Pos  *JPosition
}

func NewJLexer(filename, text string) *JLexer {
	return &JLexer{
		Text: []byte(text),
		Pos:  NewJPosition(-1, -1, 0, filename, text),
	}
}

func (l *JLexer) MakeTokens() ([]*JToken, error) {
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
			tokens = append(tokens, NewJToken(PLUS, "", l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '-':
			tokens = append(tokens, NewJToken(MINUS, "", l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '*':
			tokens = append(tokens, NewJToken(MUL, "", l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '/':
			tokens = append(tokens, NewJToken(DIV, "", l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '(':
			tokens = append(tokens, NewJToken(LPAREN, "", l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == ')':
			tokens = append(tokens, NewJToken(RPAREN, "", l.Pos, l.Pos))
			advanceAble = l.advance()
		default:
			startPos := l.Pos.Copy()
			l.advance()

			return nil, errors.Wrap(&JIllegalCharacterError{
				IllegalChar: char,
				JError: &JError{
					StartPos: startPos,
					EndPos:   l.Pos,
				},
			}, "failed to parse token")
		}
	}

	tokens = append(tokens, NewJToken(EOF, "", l.Pos, l.Pos))

	return tokens, nil
}

func (l *JLexer) makeNumberToken() (*JToken, bool) {
	advanceAble := true
	isFloat := false
	startPos := l.Pos.Copy()
	var numStrBuilder strings.Builder

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

		numStrBuilder.WriteByte(char)

		if !l.advance() {
			advanceAble = false

			break
		}
	}

	if isFloat {
		return NewJToken(FLOAT, numStrBuilder.String(), startPos, l.Pos), advanceAble
	}

	return NewJToken(INT, numStrBuilder.String(), startPos, l.Pos), advanceAble
}

func (l *JLexer) advance() bool {
	if l.Pos.Index+1 >= len(l.Text) {
		l.Pos.Advance(l.Text)

		return false
	}

	l.Pos.Advance(l.Text)

	return true
}
