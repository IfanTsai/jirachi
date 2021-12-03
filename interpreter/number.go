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

func (n *JNumber) EqualTo(other *JNumber) (*JNumber, error) {
	var res bool

	switch otherValue := other.Value.(type) {
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
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to compare equal")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) NotEqualTo(other *JNumber) (*JNumber, error) {
	var res bool

	switch otherValue := other.Value.(type) {
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
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to compare not equal")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) LessThen(other *JNumber) (*JNumber, error) {
	var res bool

	switch otherValue := other.Value.(type) {
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
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to compare less then")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) LessThenOrEqualTo(other *JNumber) (*JNumber, error) {
	var res bool

	switch otherValue := other.Value.(type) {
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
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to compare less then or equal")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) GreaterThen(other *JNumber) (*JNumber, error) {
	var res bool

	switch otherValue := other.Value.(type) {
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
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to compare greater then")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) GreaterThenOrEqualTo(other *JNumber) (*JNumber, error) {
	var res bool

	switch otherValue := other.Value.(type) {
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
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to compare greater then or equal")
	}

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) AndBy(other *JNumber) (*JNumber, error) {
	if !isNumber(other.Value) {
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to perform and logic operation")
	}

	res := numberToBool(n.Value) && numberToBool(other.Value)

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) OrBy(other *JNumber) (*JNumber, error) {
	if !isNumber(other.Value) {
		return nil, errors.Wrap(&common.JNumberTypeError{
			JError: &common.JError{
				StartPos: other.StartPos,
				EndPos:   other.EndPos,
			},
			Number: nil,
		}, "failed to perform or logic operation")
	}

	res := numberToBool(n.Value) || numberToBool(other.Value)

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
}

func (n *JNumber) Not() (*JNumber, error) {
	res := !numberToBool(n.Value)

	return NewJNumber(boolToNumber(res)).SetJContext(n.Context), nil
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

func isNumber(n interface{}) bool {
	if _, ok := n.(int); ok {
		return true
	} else if _, ok := n.(float64); ok {
		return true
	}

	return false
}
