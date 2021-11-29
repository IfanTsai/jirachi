package lexer

type JTokenType string

const (
	INT    JTokenType = "INT"
	FLOAT  JTokenType = "FLOAT"
	PLUS   JTokenType = "PLUS"
	MINUS  JTokenType = "MINUS"
	MUL    JTokenType = "MUL"
	DIV    JTokenType = "DIV"
	LPAREN JTokenType = "LPAREN"
	RPAREN JTokenType = "RPAREN"
	EOF    JTokenType = "EOF"
)

type JToken struct {
	Type     JTokenType
	Value    string
	StartPos *JPosition
	EndPos   *JPosition
}

func NewJToken(tokenType JTokenType, value string, startPos, endPos *JPosition) *JToken {
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

func (t *JToken) String() string {
	if len(t.Value) == 0 {
		return string(t.Type)
	}

	return string(t.Type) + ":" + t.Value
}
