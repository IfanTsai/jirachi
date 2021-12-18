package object

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IfanTsai/jirachi/common"
)

type JString struct {
	*JBaseValue
}

func NewJString(value interface{}) *JString {
	return &JString{
		JBaseValue: &JBaseValue{
			Value: value,
		},
	}
}

func (s *JString) String() string {
	return fmt.Sprintf("\"%v\"", s.Value)
}

func (s *JString) SetJPos(startPos, endPos *common.JPosition) JValue {
	s.StartPos = startPos
	s.EndPos = endPos

	return s
}

func (s *JString) SetJContext(context *common.JContext) JValue {
	s.Context = context

	return s
}

func (s *JString) Copy() JValue {
	return NewJString(s.Value)
}

func (s *JString) AddTo(other JValue) (JValue, error) {
	var resString *JString

	switch otherValue := other.GetValue().(type) {
	case int:
		resString = NewJString(s.Value.(string) + strconv.Itoa(otherValue))
	case float64:
		resString = NewJString(s.Value.(string) + strconv.FormatFloat(otherValue, 'f', 2, 64))
	case string:
		resString = NewJString(s.Value.(string) + otherValue)
	default:
		return nil, s.createIllegalOperationError(other, "add")
	}

	return resString.SetJContext(s.Context), nil
}

func (s *JString) MulBy(other JValue) (JValue, error) {
	var resString *JString

	switch otherValue := other.GetValue().(type) {
	case int:
		strBuilder := strings.Builder{}
		for i := 0; i < otherValue; i++ {
			strBuilder.WriteString(s.Value.(string))
		}
		resString = NewJString(strBuilder.String())
	default:
		return nil, s.createIllegalOperationError(other, "mul")
	}

	return resString.SetJContext(s.Context), nil
}

func (s *JString) IsTrue() bool {
	return len(s.Value.(string)) > 0
}

func (s *JString) Not() (JValue, error) {
	res := !s.IsTrue()

	return NewJNumber(boolToNumber(res)).SetJContext(s.Context).(*JNumber), nil
}
