package lexer

import (
	"strconv"
	"strings"

	"github.com/IfanTsai/jirachi/common"

	"github.com/IfanTsai/jirachi/token"

	"github.com/pkg/errors"
)

type JLexer struct {
	Text []byte
	Pos  *common.JPosition
}

func NewJLexer(filename, text string) *JLexer {
	return &JLexer{
		Text: []byte(text),
		Pos:  common.NewJPosition(-1, -1, 0, filename, text),
	}
}

func (l *JLexer) MakeTokens() ([]*token.JToken, error) {
	tokens := make([]*token.JToken, 0, len(l.Text))

	for advanceAble := l.advance(); advanceAble; {
		char := l.Text[l.Pos.Index]
		switch {
		case char == ' ' || char == '\t':
			advanceAble = l.advance()
		case isDigit(char):
			var numberToken *token.JToken
			var err error
			numberToken, advanceAble, err = l.makeNumberToken()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, numberToken)
		case isLetters(char):
			var numberToken *token.JToken
			numberToken, advanceAble = l.makeIdentifierToken()
			tokens = append(tokens, numberToken)
		case char == '+':
			tokens = append(tokens, token.NewJToken(token.PLUS, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '-':
			tokens = append(tokens, token.NewJToken(token.MINUS, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '*':
			tokens = append(tokens, token.NewJToken(token.MUL, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '/':
			tokens = append(tokens, token.NewJToken(token.DIV, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '^':
			tokens = append(tokens, token.NewJToken(token.POW, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '=':
			tokens = append(tokens, token.NewJToken(token.EQ, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '(':
			tokens = append(tokens, token.NewJToken(token.LPAREN, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == ')':
			tokens = append(tokens, token.NewJToken(token.RPAREN, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		default:
			startPos := l.Pos.Copy()
			l.advance()

			return nil, errors.Wrap(&common.JIllegalCharacterError{
				IllegalChar: char,
				JError: &common.JError{
					StartPos: startPos,
					EndPos:   l.Pos,
				},
			}, "failed to parse token")
		}
	}

	tokens = append(tokens, token.NewJToken(token.EOF, nil, l.Pos, l.Pos))

	return tokens, nil
}

func (l *JLexer) makeNumberToken() (*token.JToken, bool, error) {
	advanceAble := true
	isFloat := false
	startPos := l.Pos.Copy()
	var numStrBuilder strings.Builder

	for {
		char := l.Text[l.Pos.Index]

		if char != '.' && !isDigit(char) {
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
		floatNum, err := strconv.ParseFloat(numStrBuilder.String(), 64)
		if err != nil {
			return nil, false, errors.Wrapf(err, "failed to convert %s to float64", numStrBuilder.String())
		}
		return token.NewJToken(token.FLOAT, floatNum, startPos, l.Pos), advanceAble, nil
	}

	intNum, err := strconv.Atoi(numStrBuilder.String())
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to convert %s to int", numStrBuilder.String())
	}

	return token.NewJToken(token.INT, intNum, startPos, l.Pos), advanceAble, nil
}

func (l *JLexer) makeIdentifierToken() (*token.JToken, bool) {
	advanceAble := true
	startPos := l.Pos.Copy()
	var identifierStrBuilder strings.Builder

	for {
		char := l.Text[l.Pos.Index]

		if !isLetters(char) && !isDigit(char) {
			break
		}

		identifierStrBuilder.WriteByte(char)

		if !l.advance() {
			advanceAble = false

			break
		}
	}

	identifier := identifierStrBuilder.String()
	if token.IsKeyword(identifier) {
		return token.NewJToken(token.KEYWORD, identifier, startPos, l.Pos), advanceAble
	}

	return token.NewJToken(token.IDENTIFIER, identifier, startPos, l.Pos), advanceAble
}

func (l *JLexer) advance() bool {
	if l.Pos.Index+1 >= len(l.Text) {
		l.Pos.Advance(l.Text)

		return false
	}

	l.Pos.Advance(l.Text)

	return true
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func isLetters(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}
