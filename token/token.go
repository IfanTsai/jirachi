package token

import (
	"encoding/json"
	"strconv"

	"github.com/IfanTsai/jirachi/pkg/set"

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
	AND      = "AND"
	OR       = "OR"
	NOT      = "NOT"
	IF       = "IF"
	THEN     = "THEN"
	ELIF     = "ELIF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	TO       = "TO"
	STEP     = "STEP"
	WHILE    = "WHILE"
	FUN      = "FUN"
	END      = "END"
	RETURN   = "RETURN"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
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
		newValue, _ := json.Marshal(value)
		return string(newValue)
	}
}

func IsKeyword(identifier string) bool {
	return KEYWORDS.Contains(identifier)
}
