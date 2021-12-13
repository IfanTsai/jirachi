package interpreter

import (
	"fmt"

	"github.com/IfanTsai/jirachi/common"
)

type JValue interface {
	fmt.Stringer
	GetValue() interface{}
	SetJPos(startPos, endPos *common.JPosition) JValue
	SetJContext(context *common.JContext) JValue
	GetContext() *common.JContext
	GetStartPos() *common.JPosition
	GetEndPos() *common.JPosition
	Copy() JValue
	AddTo(other JValue) (JValue, error)
	SubBy(other JValue) (JValue, error)
	MulBy(other JValue) (JValue, error)
	DivBy(other JValue) (JValue, error)
	PowBy(other JValue) (JValue, error)
	EqualTo(other JValue) (JValue, error)
	NotEqualTo(other JValue) (JValue, error)
	LessThan(other JValue) (JValue, error)
	LessThanOrEqualTo(other JValue) (JValue, error)
	GreaterThan(other JValue) (JValue, error)
	GreaterThanOrEqualTo(other JValue) (JValue, error)
	AndBy(other JValue) (JValue, error)
	OrBy(other JValue) (JValue, error)
	Not() (JValue, error)
	IsTrue() bool
	Execute(args []JValue) (JValue, error)
}

type JBaseValue struct {
	Value    interface{} // used to number(only support int and float64) or function name
	StartPos *common.JPosition
	EndPos   *common.JPosition
	Context  *common.JContext
}

func (v *JBaseValue) String() string {
	return fmt.Sprintf("%v", v.Value)
}

func (v *JBaseValue) GetValue() interface{} {
	return v.Value
}

func (v *JBaseValue) GetStartPos() *common.JPosition {
	return v.StartPos
}

func (v *JBaseValue) GetEndPos() *common.JPosition {
	return v.EndPos
}

func (v *JBaseValue) GetContext() *common.JContext {
	return v.Context
}

func (v *JBaseValue) Copy() JValue {
	panic("implement me")
}

func (v *JBaseValue) AddTo(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) SubBy(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) MulBy(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) DivBy(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) PowBy(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) EqualTo(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) NotEqualTo(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) LessThan(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) LessThanOrEqualTo(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) GreaterThan(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) GreaterThanOrEqualTo(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) AndBy(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) OrBy(other JValue) (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) Not() (JValue, error) {
	panic("implement me")
}

func (v *JBaseValue) IsTrue() bool {
	panic("implement me")
}

func (v *JBaseValue) Execute(args []JValue) (JValue, error) {
	panic("implement me")
}
