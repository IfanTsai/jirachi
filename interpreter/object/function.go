package object

import (
	"github.com/IfanTsai/jirachi/common"
	"github.com/IfanTsai/jirachi/parser"
)

type JFunction struct {
	*JBaseValue
	ArgNames []string
	BodyNode parser.JNode
}

func NewJFunction(funcName interface{}, argNames []string, bodyNode parser.JNode) *JFunction {
	if funcName == nil {
		funcName = "<anonymous>"
	}

	return &JFunction{
		JBaseValue: &JBaseValue{
			Value: funcName,
		},
		ArgNames: argNames,
		BodyNode: bodyNode,
	}
}

func (f *JFunction) SetJPos(startPos, endPos *common.JPosition) JValue {
	f.StartPos = startPos
	f.EndPos = endPos

	return f
}

func (f *JFunction) SetJContext(context *common.JContext) JValue {
	f.Context = context

	return f
}

func (f *JFunction) Copy() JValue {
	return NewJFunction(f.Value, f.ArgNames, f.BodyNode)
}

func (f *JFunction) String() string {
	return "<function " + f.JBaseValue.String() + ">"
}
