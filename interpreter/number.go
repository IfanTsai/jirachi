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
}

func NewJNumber(value interface{}) *JNumber {
	return &JNumber{
		Value: value,
	}
}

func (n *JNumber) SetPos(startPos, endPos *common.JPosition) *JNumber {
	n.StartPos = startPos
	n.EndPos = endPos

	return n
}

func (n *JNumber) AddTo(other *JNumber) (*JNumber, error) {
	switch otherValue := other.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			return NewJNumber(num + otherValue), nil
		} else {
			return NewJNumber(n.Value.(float64) + float64(otherValue)), nil
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			return NewJNumber(float64(num) + otherValue), nil
		} else {
			return NewJNumber(n.Value.(float64) + otherValue), nil
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
}

func (n *JNumber) SubBy(other *JNumber) (*JNumber, error) {
	switch otherValue := other.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			return NewJNumber(num - otherValue), nil
		} else {
			return NewJNumber(n.Value.(float64) - float64(otherValue)), nil
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			return NewJNumber(float64(num) - otherValue), nil
		} else {
			return NewJNumber(n.Value.(float64) - otherValue), nil
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
}

func (n *JNumber) MulBy(other *JNumber) (*JNumber, error) {
	switch otherValue := other.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			return NewJNumber(num * otherValue), nil
		} else {
			return NewJNumber(n.Value.(float64) * float64(otherValue)), nil
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			return NewJNumber(float64(num) * otherValue), nil
		} else {
			return NewJNumber(n.Value.(float64) * otherValue), nil
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
}

func (n *JNumber) DivBy(other *JNumber) (*JNumber, error) {
	switch otherValue := other.Value.(type) {
	case int:
		if otherValue == 0 {
			return nil, errors.Wrap(&common.JRunTimeError{
				JError: &common.JError{
					StartPos: other.StartPos,
					EndPos:   other.EndPos,
				},
				Details: "Division by zero",
			}, "failed to div")
		}

		if num, ok := n.Value.(int); ok {
			return NewJNumber(num / otherValue), nil
		} else {
			return NewJNumber(n.Value.(float64) / float64(otherValue)), nil
		}
	case float64:
		if otherValue == 0.0 {
			return nil, errors.Wrap(&common.JRunTimeError{
				JError: &common.JError{
					StartPos: other.StartPos,
					EndPos:   other.EndPos,
				},
				Details: "Division by zero",
			}, "failed to div")
		}

		if num, ok := n.Value.(int); ok {
			return NewJNumber(float64(num) / otherValue), nil
		} else {
			return NewJNumber(n.Value.(float64) / otherValue), nil
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
}

func (n *JNumber) PowBy(other *JNumber) (*JNumber, error) {
	switch otherValue := other.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			return NewJNumber(math.Pow(float64(num), float64(otherValue))), nil
		} else {
			return NewJNumber(math.Pow(n.Value.(float64), float64(otherValue))), nil
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			return NewJNumber(math.Pow(float64(num), otherValue)), nil
		} else {
			return NewJNumber(math.Pow(n.Value.(float64), otherValue)), nil
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
}
