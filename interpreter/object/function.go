package object

import (
	"fmt"

	"github.com/IfanTsai/jirachi/common"
	"github.com/IfanTsai/jirachi/parser"
	"github.com/pkg/errors"
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

func (f *JFunction) CheckArgs(argValues []JValue) error {
	if len(argValues) > len(f.ArgNames) {
		return errors.Wrap(&common.JRunTimeError{
			JError: &common.JError{
				StartPos: f.GetStartPos(),
				EndPos:   f.GetEndPos(),
			},
			Context: f.GetContext(),
			Details: fmt.Sprintf("%d too many args passed into %v", len(argValues)-len(f.ArgNames), f.GetValue()),
		}, "failed to execute")
	}

	if len(argValues) < len(f.ArgNames) {
		return errors.Wrap(&common.JRunTimeError{
			JError: &common.JError{
				StartPos: f.GetStartPos(),
				EndPos:   f.GetEndPos(),
			},
			Context: f.GetContext(),
			Details: fmt.Sprintf("%d too few passed into %v", len(f.ArgNames)-len(argValues), f.GetValue()),
		}, "failed to execute")
	}

	return nil
}
