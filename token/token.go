package token

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/IfanTsai/go-lib/set"
	"github.com/IfanTsai/jirachi/common"
)

type JTokenType string

const (
	INT        JTokenType = "INT"
	FLOAT      JTokenType = "FLOAT"
	STRING     JTokenType = "STRING"
	IDENTIFIER JTokenType = "IDENTIFIER"
	KEYWORD    JTokenType = "KEYWORD"
	PLUS       JTokenType = "PLUS"    // +
	MINUS      JTokenType = "MINUS"   // -
	MUL        JTokenType = "MUL"     // *
	DIV        JTokenType = "DIV"     // /
	POW        JTokenType = "POW"     // ^
	EQ         JTokenType = "EQ"      // =
	LPAREN     JTokenType = "LPAREN"  // (
	RPAREN     JTokenType = "RPAREN"  // )
	LSQUARE    JTokenType = "LSQUARE" // [
	RSQUARE    JTokenType = "RSQUARE" // ]
	LBRACE     JTokenType = "LBRACE"  // {
	RBRACE     JTokenType = "RBRACE"  // }
	COLON      JTokenType = "COLON"   // :
	EE         JTokenType = "EE"      // ==
	NE         JTokenType = "NE"      // !=
	LT         JTokenType = "LT"      // <
	GT         JTokenType = "GT"      // >
	LTE        JTokenType = "LTE"     // <=
	GTE        JTokenType = "GTE"     // >=
	COMMA      JTokenType = "COMMA"   // ,
	ARROW      JTokenType = "ARROW"   // ->
	NEWLINE    JTokenType = "NEWLINE"
	EOF        JTokenType = "EOF"
)

const (
	AND      = "and"
	OR       = "or"
	NOT      = "not"
	IF       = "if"
	THEN     = "then"
	ELIF     = "elif"
	ELSE     = "else"
	FOR      = "for"
	TO       = "to"
	STEP     = "step"
	WHILE    = "while"
	FUN      = "fun"
	END      = "end"
	RETURN   = "return"
	BREAK    = "break"
	CONTINUE = "continue"
)

var KEYWORDS = set.NewSet(
	AND,
	OR,
	NOT,
	IF,
	THEN,
	ELIF,
	ELSE,
	FOR,
	TO,
	STEP,
	WHILE,
	FUN,
	END,
	RETURN,
	BREAK,
	CONTINUE,
)

type JToken struct {
	Type     JTokenType
	Value    interface{}
	StartPos *common.JPosition
	EndPos   *common.JPosition
}

func NewJToken(tokenType JTokenType, value interface{}, startPos, endPos *common.JPosition) *JToken {
	token := &JToken{
		Type:     tokenType,
		Value:    value,
		StartPos: startPos.Copy(),
		EndPos:   endPos.Copy(),
	}

	if startPos == endPos {
		token.EndPos.Advance(nil)
	}

	return token
}

func (t *JToken) IsKeyWord() bool {
	return t.Type == KEYWORD
}

func (t *JToken) Match(typ JTokenType, value interface{}) bool {
	return t.Type == typ && t.Value == value
}

func (t *JToken) String() string {
	if t.Value == nil {
		return string(t.Type)
	}

	return string(t.Type) + ":" + t.ValueToString()
}

func (t *JToken) ValueToString() string {
	switch value := t.Value.(type) {
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case int:
		return strconv.Itoa(value)
	case string:
		return value
	default:
		newValue, err := json.Marshal(value)
		if err != nil {
			fmt.Println("cannot marshal value")
		}

		return string(newValue)
	}
}

func IsKeyword(identifier string) bool {
	return KEYWORDS.Contains(identifier)
}
