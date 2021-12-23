package object

import "github.com/IfanTsai/jirachi/common"

type JNull struct {
	*JBaseValue
}

func NewJNull() *JNull {
	return &JNull{
		JBaseValue: &JBaseValue{},
	}
}

func (n *JNull) SetJPos(startPos, endPos *common.JPosition) JValue {
	n.StartPos = startPos
	n.EndPos = endPos

	return n
}

func (n *JNull) SetJContext(context *common.JContext) JValue {
	n.Context = context

	return n
}

func (n *JNull) Copy() JValue {
	return n
}

func (n *JNull) String() string {
	return "<null>"
}

func (n *JNull) EqualTo(other JValue) (JValue, error) {
	res := false

	if _, ok := other.(*JNull); ok {
		res = true
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNull) NotEqualTo(other JValue) (JValue, error) {
	res := true

	if _, ok := other.(*JNull); ok {
		res = false
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}
