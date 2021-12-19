package object

import "github.com/IfanTsai/jirachi/common"

type ExecuteFunc func(function *JBuiltInFunction, args []JValue) (JValue, error)

type JBuiltInFunction struct {
	*JFunction
	ExecuteCb ExecuteFunc
}

func NewJBuiltInFunction(funcName interface{}, argNames []string, executeFunc ExecuteFunc) *JBuiltInFunction {
	return &JBuiltInFunction{
		JFunction: &JFunction{
			JBaseValue: &JBaseValue{
				Value: funcName,
			},
			ArgNames: argNames,
		},
		ExecuteCb: executeFunc,
	}
}

func (bif *JBuiltInFunction) SetJPos(startPos, endPos *common.JPosition) JValue {
	bif.StartPos = startPos
	bif.EndPos = endPos

	return bif
}

func (bif *JBuiltInFunction) SetJContext(context *common.JContext) JValue {
	bif.Context = context

	return bif
}

func (bif *JBuiltInFunction) Copy() JValue {
	return NewJBuiltInFunction(bif.Value, bif.ArgNames, bif.ExecuteCb)
}

func (bif *JBuiltInFunction) Execute(args []JValue) (JValue, error) {
	if err := bif.CheckArgs(args); err != nil {
		return nil, err
	}

	return bif.ExecuteCb(bif, args)
}
