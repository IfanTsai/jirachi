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
)

type JToken struct {
	Type  JTokenType
	Value string
}

func NewJToken(tokenType JTokenType, value string) *JToken {
	return &JToken{
		Type:  tokenType,
		Value: value,
	}
}

func (t *JToken) String() string {
	if len(t.Value) == 0 {
		return string(t.Type)
	}

	return string(t.Type) + ":" + t.Value
}
