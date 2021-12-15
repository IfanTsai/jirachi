package lexer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IfanTsai/jirachi/common"

	"github.com/IfanTsai/jirachi/token"

	"github.com/pkg/errors"
)

var escapeChars = map[byte]byte{
	'n': '\n',
	't': '\t',
}

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

	var err error
	var tok *token.JToken

	for advanceAble := l.advance(); advanceAble; {
		char := l.getCurrentChar()
		switch {
		case char == ' ' || char == '\t':
			advanceAble = l.advance()
		case isDigit(char):
			tok, advanceAble, err = l.makeNumberToken()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
		case isLetters(char):
			tok, advanceAble = l.makeIdentifierToken()
			tokens = append(tokens, tok)
		case char == '"' || char == '\'':
			tok, advanceAble, err = l.makeString()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
		case char == '+':
			tokens = append(tokens, token.NewJToken(token.PLUS, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '-':
			tok, advanceAble = l.makeMinusOrArrowToken()
			tokens = append(tokens, tok)
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
			tok, advanceAble = l.makeEqualToken()
			tokens = append(tokens, tok)
		case char == '(':
			tokens = append(tokens, token.NewJToken(token.LPAREN, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == ')':
			tokens = append(tokens, token.NewJToken(token.RPAREN, nil, l.Pos, l.Pos))
			advanceAble = l.advance()
		case char == '!':
			tok, advanceAble, err = l.makeNotEqualToken()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
		case char == '<':
			tok, advanceAble = l.makeLessThanToken()
			tokens = append(tokens, tok)
		case char == '>':
			tok, advanceAble = l.makeGreaterThanToken()
			tokens = append(tokens, tok)
		case char == ',':
			tokens = append(tokens, token.NewJToken(token.COMMA, nil, l.Pos, l.Pos))
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
		char := l.getCurrentChar()

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
		char := l.getCurrentChar()

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

func (l *JLexer) makeString() (*token.JToken, bool, error) {
	quote := l.getCurrentChar() // ' or "
	strBuilder := strings.Builder{}
	startPos := l.Pos.Copy()
	isEscape := false
	advanceAble := l.advance()

	for advanceAble && (l.getCurrentChar() != quote || isEscape) {
		if isEscape {
			escapeChar, ok := escapeChars[l.getCurrentChar()]
			if !ok {
				escapeChar = l.getCurrentChar()
			}
			strBuilder.WriteByte(escapeChar)

			isEscape = false
		} else if l.getCurrentChar() == '\\' {
			isEscape = true
		} else {
			strBuilder.WriteByte(l.getCurrentChar())
		}

		if !l.advance() {
			advanceAble = false
			goto unexpectedQuote
		}
	}

unexpectedQuote:
	if !advanceAble || l.getCurrentChar() != quote {
		if !l.advance() {
			fmt.Println(startPos.Index)
			return nil, advanceAble, errors.Wrap(&common.JExpectedCharacterError{
				JError: &common.JError{
					StartPos: startPos,
					EndPos:   startPos.Copy().Advance(l.Text),
				},
				ExpectedChar: quote,
			}, "failed to make string token")
		}
	}

	advanceAble = l.advance()

	return token.NewJToken(token.STRING, strBuilder.String(), startPos, l.Pos), advanceAble, nil
}

func (l *JLexer) makeMinusOrArrowToken() (*token.JToken, bool) {
	startPos := l.Pos.Copy()
	advanceAble := l.advance()
	tokenType := token.MINUS

	if l.getCurrentChar() == '>' {
		advanceAble = l.advance()

		tokenType = token.ARROW
	}

	return token.NewJToken(tokenType, nil, startPos, l.Pos), advanceAble
}

func (l *JLexer) makeNotEqualToken() (*token.JToken, bool, error) {
	startPos := l.Pos.Copy()
	advanceAble := l.advance()

	if l.getCurrentChar() == '=' {
		advanceAble = l.advance()

		return token.NewJToken(token.NE, nil, startPos, l.Pos), advanceAble, nil
	}

	return nil, advanceAble, errors.Wrap(&common.JExpectedCharacterError{
		JError: &common.JError{
			StartPos: startPos,
			EndPos:   l.Pos,
		},
		ExpectedChar: '=',
	}, "failed to make not equal token")
}

func (l *JLexer) makeEqualToken() (*token.JToken, bool) {
	startPos := l.Pos.Copy()
	advanceAble := l.advance()
	tokenType := token.EQ

	if l.getCurrentChar() == '=' {
		advanceAble = l.advance()

		tokenType = token.EE
	}

	return token.NewJToken(tokenType, nil, startPos, l.Pos), advanceAble
}

func (l *JLexer) makeLessThanToken() (*token.JToken, bool) {
	startPos := l.Pos.Copy()
	advanceAble := l.advance()
	tokenType := token.LT

	if l.getCurrentChar() == '=' {
		advanceAble = l.advance()

		tokenType = token.LTE
	}

	return token.NewJToken(tokenType, nil, startPos, l.Pos), advanceAble
}

func (l *JLexer) makeGreaterThanToken() (*token.JToken, bool) {
	startPos := l.Pos.Copy()
	advanceAble := l.advance()
	tokenType := token.GT

	if l.getCurrentChar() == '=' {
		advanceAble = l.advance()

		tokenType = token.GTE
	}

	return token.NewJToken(tokenType, nil, startPos, l.Pos), advanceAble
}

func (l *JLexer) advance() bool {
	if l.Pos.Index+1 >= len(l.Text) {
		l.Pos.Advance(l.Text)

		return false
	}

	l.Pos.Advance(l.Text)

	return true
}

func (l *JLexer) getCurrentChar() byte {
	return l.Text[l.Pos.Index]
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func isLetters(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}
