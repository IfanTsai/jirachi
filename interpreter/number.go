package interpreter

import (
	"math"

	"github.com/IfanTsai/jirachi/common"
	"github.com/pkg/errors"
)

type JNumberType int

type JNumber struct {
	Value    interface{} // only support int and float64
	StartPos *common.JPosition
	EndPos   *common.JPosition
	Context  *common.JContext
}

func NewJNumber(value interface{}) *JNumber {
	return &JNumber{
		Value: value,
	}
}

func (n *JNumber) SetJPos(startPos, endPos *common.JPosition) *JNumber {
	n.StartPos = startPos
	n.EndPos = endPos

	return n
}

func (n *JNumber) SetJContext(context *common.JContext) *JNumber {
	n.Context = context

	return n
}

func (n *JNumber) AddTo(other *JNumber) (*JNumber, error) {
	var resNumber *JNumber

	switch otherValue := other.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			resNumber = NewJNumber(num + otherValue)
		} else {
			resNumber = NewJNumber(n.Value.(float64) + float64(otherValue))
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			resNumber = NewJNumber(float64(num) + otherValue)
		} else {
			resNumber = NewJNumber(n.Value.(float64) + otherValue)
		}
	default:
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to add")
	}

	return resNumber.SetJContext(n.Context), nil
}

func (n *JNumber) SubBy(other *JNumber) (*JNumber, error) {
	var resNumber *JNumber

	switch otherValue := other.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			resNumber = NewJNumber(num - otherValue)
		} else {
			resNumber = NewJNumber(n.Value.(float64) - float64(otherValue))
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			resNumber = NewJNumber(float64(num) - otherValue)
		} else {
			resNumber = NewJNumber(n.Value.(float64) - otherValue)
		}
	default:
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to sub")
	}

	return resNumber.SetJContext(n.Context), nil
}

func (n *JNumber) MulBy(other *JNumber) (*JNumber, error) {
	var resNumber *JNumber

	switch otherValue := other.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			resNumber = NewJNumber(num * otherValue)
		} else {
			resNumber = NewJNumber(n.Value.(float64) * float64(otherValue))
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			resNumber = NewJNumber(float64(num) * otherValue)
		} else {
			resNumber = NewJNumber(n.Value.(float64) * otherValue)
		}
	default:
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to mul")
	}

	return resNumber.SetJContext(n.Context), nil
}

func (n *JNumber) DivBy(other *JNumber) (*JNumber, error) {
	var resNumber *JNumber

	switch otherValue := other.Value.(type) {
	case int:
		if otherValue == 0 {
			return nil, errors.Wrap(&common.JRunTimeError{
				JError: &common.JError{
					StartPos: other.StartPos,
					EndPos:   other.EndPos,
				},
				Context: n.Context,
				Details: "Division by zero",
			}, "failed to div")
		}

		if num, ok := n.Value.(int); ok {
			resNumber = NewJNumber(num / otherValue)
		} else {
			resNumber = NewJNumber(n.Value.(float64) / float64(otherValue))
		}
	case float64:
		if otherValue == 0.0 {
			return nil, errors.Wrap(&common.JRunTimeError{
				JError: &common.JError{
					StartPos: other.StartPos,
					EndPos:   other.EndPos,
				},
				Context: n.Context,
				Details: "Division by zero",
			}, "failed to div")
		}

		if num, ok := n.Value.(int); ok {
			resNumber = NewJNumber(float64(num) / otherValue)
		} else {
			resNumber = NewJNumber(n.Value.(float64) / otherValue)
		}
	default:
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to div")
	}

	return resNumber.SetJContext(n.Context), nil
}

func (n *JNumber) PowBy(other *JNumber) (*JNumber, error) {
	var resNumber *JNumber

	switch otherValue := other.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			resNumber = NewJNumber(math.Pow(float64(num), float64(otherValue)))
		} else {
			resNumber = NewJNumber(math.Pow(n.Value.(float64), float64(otherValue)))
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			resNumber = NewJNumber(math.Pow(float64(num), otherValue))
		} else {
			resNumber = NewJNumber(math.Pow(n.Value.(float64), otherValue))
		}
	default:
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to pow")
	}

	return resNumber.SetJContext(n.Context), nil
}
