package interpreter

import (
	"math"

	"github.com/IfanTsai/jirachi/common"
	"github.com/pkg/errors"
)

type JNumberType int

type JNumber struct {
	*JBaseValue
}

func NewJNumber(value interface{}) *JNumber {
	return &JNumber{
		JBaseValue: &JBaseValue{
			Value: value,
		},
	}
}

func (n *JNumber) SetJPos(startPos, endPos *common.JPosition) JValue {
	n.StartPos = startPos
	n.EndPos = endPos

	return n
}

func (n *JNumber) SetJContext(context *common.JContext) JValue {
	n.Context = context

	return n
}

func (n *JNumber) Copy() JValue {
	return NewJNumber(n.Value)
}

func (n *JNumber) AddTo(other JValue) (JValue, error) {
	var resNumber *JNumber
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "add")
	}

	switch otherValue := otherNumber.Value.(type) {
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
		return nil, createNumberTypeError(other, "add")
	}

	return resNumber.SetJContext(n.Context), nil
}

func (n *JNumber) SubBy(other JValue) (JValue, error) {
	var resNumber *JNumber
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "sub")
	}

	switch otherValue := otherNumber.Value.(type) {
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
		return nil, createNumberTypeError(other, "sub")
	}

	return resNumber.SetJContext(n.Context), nil
}

func (n *JNumber) MulBy(other JValue) (JValue, error) {
	var resNumber *JNumber
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "mul")
	}

	switch otherValue := otherNumber.Value.(type) {
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
		return nil, createNumberTypeError(other, "mul")
	}

	return resNumber.SetJContext(n.Context), nil
}

func (n *JNumber) DivBy(other JValue) (JValue, error) {
	var resNumber *JNumber
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "div")
	}

	switch otherValue := otherNumber.Value.(type) {
	case int:
		if otherValue == 0 {
			return nil, errors.Wrap(&common.JRunTimeError{
				JError: &common.JError{
					StartPos: otherNumber.StartPos,
					EndPos:   otherNumber.EndPos,
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
					StartPos: otherNumber.StartPos,
					EndPos:   otherNumber.EndPos,
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
		return nil, createNumberTypeError(other, "div")
	}

	return resNumber.SetJContext(n.Context), nil
}

func (n *JNumber) PowBy(other JValue) (JValue, error) {
	var resNumber *JNumber
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "pow")
	}

	switch otherValue := otherNumber.Value.(type) {
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
		return nil, createNumberTypeError(other, "pow")
	}

	return resNumber.SetJContext(n.Context), nil
}

func (n *JNumber) EqualTo(other JValue) (JValue, error) {
	var res bool
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "compare equal")
	}

	switch otherValue := otherNumber.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			res = num == otherValue
		} else {
			res = int(n.Value.(float64)) == otherValue
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			res = float64(num) == otherValue
		} else {
			res = n.Value.(float64) == otherValue
		}
	default:
		return nil, createNumberTypeError(other, "compare equal")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) NotEqualTo(other JValue) (JValue, error) {
	var res bool
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "compare not equal")
	}

	switch otherValue := otherNumber.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			res = num != otherValue
		} else {
			res = int(n.Value.(float64)) != otherValue
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			res = float64(num) != otherValue
		} else {
			res = n.Value.(float64) != otherValue
		}
	default:
		return nil, createNumberTypeError(other, "compare not equal")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) LessThan(other JValue) (JValue, error) {
	var res bool
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "compare less than")
	}

	switch otherValue := otherNumber.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			res = num < otherValue
		} else {
			res = int(n.Value.(float64)) < otherValue
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			res = float64(num) < otherValue
		} else {
			res = n.Value.(float64) < otherValue
		}
	default:
		return nil, createNumberTypeError(other, "compare less than")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) LessThanOrEqualTo(other JValue) (JValue, error) {
	var res bool
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "compare less than or equal")
	}

	switch otherValue := otherNumber.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			res = num <= otherValue
		} else {
			res = int(n.Value.(float64)) <= otherValue
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			res = float64(num) <= otherValue
		} else {
			res = n.Value.(float64) <= otherValue
		}
	default:
		return nil, createNumberTypeError(other, "compare less than or equal")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) GreaterThan(other JValue) (JValue, error) {
	var res bool
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "compare greater than")
	}

	switch otherValue := otherNumber.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			res = num > otherValue
		} else {
			res = int(n.Value.(float64)) > otherValue
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			res = float64(num) > otherValue
		} else {
			res = n.Value.(float64) > otherValue
		}
	default:
		return nil, createNumberTypeError(other, "compare greater than")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) GreaterThanOrEqualTo(other JValue) (JValue, error) {
	var res bool
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "compare greater than or equal")
	}

	switch otherValue := otherNumber.Value.(type) {
	case int:
		if num, ok := n.Value.(int); ok {
			res = num >= otherValue
		} else {
			res = int(n.Value.(float64)) >= otherValue
		}
	case float64:
		if num, ok := n.Value.(int); ok {
			res = float64(num) >= otherValue
		} else {
			res = n.Value.(float64) >= otherValue
		}
	default:
		return nil, createNumberTypeError(other, "compare greater than or equal")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context).(*JNumber), nil
}

func (n *JNumber) AndBy(other JValue) (JValue, error) {
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "and")
	}

	// short-circuit evaluation
	if numberToBool(n.Value) {
		return NewJNumber(otherNumber.Value).SetJContext(n.Context).(*JNumber), nil
	} else {
		return NewJNumber(n.Value).SetJContext(n.Context).(*JNumber), nil
	}
}

func (n *JNumber) OrBy(other JValue) (JValue, error) {
	otherNumber, ok := other.(*JNumber)
	if !ok {
		return nil, createNumberTypeError(other, "or")
	}

	// short-circuit evaluation
	if numberToBool(n.Value) {
		return NewJNumber(n.Value).SetJContext(n.Context).(*JNumber), nil
	} else {
		return NewJNumber(otherNumber.Value).SetJContext(n.Context).(*JNumber), nil
	}
}

func (n *JNumber) Not() (JValue, error) {
	res := !numberToBool(n.Value)

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context).(*JNumber), nil
}

func (n *JNumber) IsTrue() bool {
	return numberToBool(n.Value)
}

func boolToNumber(b bool) interface{} {
	if b {
		return 1
	}

	return 0
}

func numberToBool(n interface{}) bool {
	if value, ok := n.(int); ok {
		return value != 0
	}

	return n.(float64) != 0
}

func createNumberTypeError(number JValue, operation string) error {
	return errors.Wrap(&common.JNumberTypeError{
		JError: &common.JError{
			StartPos: number.GetStartPos(),
			EndPos:   number.GetEndPos(),
		},
		Number: nil,
	}, "failed to "+operation)
}
