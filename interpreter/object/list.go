package object

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/IfanTsai/jirachi/common"
)

type JList struct {
	*JBaseValue
	ElementValues []JValue
}

func NewJList(elementValues []JValue) *JList {
	return &JList{
		JBaseValue:    &JBaseValue{},
		ElementValues: elementValues,
	}
}

func (l *JList) SetJPos(startPos, endPos *common.JPosition) JValue {
	l.StartPos = startPos
	l.EndPos = endPos

	return l
}

func (l *JList) SetJContext(context *common.JContext) JValue {
	l.Context = context

	return l
}

func (l *JList) Copy() JValue {
	return NewJList(l.ElementValues)
}

func (l *JList) AddTo(other JValue) (JValue, error) {
	resList := l.deepCopy()

	switch otherValue := other.(type) {
	case *JList:
		resList.ElementValues = append(resList.ElementValues, otherValue.ElementValues...)
	default:
		resList.ElementValues = append(resList.ElementValues, otherValue)
	}

	return resList.SetJContext(l.Context), nil
}

func (l *JList) MulBy(other JValue) (JValue, error) {
	resList := l.deepCopy()

	switch otherValue := other.(type) {
	case *JList:
		for index := range resList.ElementValues {
			mulValue, err := resList.ElementValues[index].MulBy(otherValue.ElementValues[index])
			if err != nil {
				return nil, errors.WithMessagef(err, "failed to mul")
			}
			resList.ElementValues[index] = mulValue
		}
	case *JNumber:
		if intValue, ok := otherValue.GetValue().(int); ok {
			for i := 0; i < intValue-1; i++ {
				resList.ElementValues = append(resList.ElementValues, l.ElementValues...)
			}
		} else {
			return nil, createNumberTypeError(other, "mul")
		}
	}

	return resList.SetJContext(l.Context), nil
}

func (l *JList) SubBy(other JValue) (JValue, error) {
	resList := l.deepCopy()

	if index, ok := other.GetValue().(int); ok {
		if err := l.checkIndex(index, other); err != nil {
			return nil, err
		}

		resList.ElementValues = append(resList.ElementValues[:index], resList.ElementValues[index+1:]...)
	} else {
		return nil, createNumberTypeError(other, "sub")
	}

	return resList, nil
}

func (l *JList) Index(arg JValue) (JValue, error) {
	if index, ok := arg.GetValue().(int); ok {
		if err := l.checkIndex(index, arg); err != nil {
			return nil, err
		}

		return l.ElementValues[index], nil
	} else {
		return nil, createNumberTypeError(arg, "index")
	}
}

func (l *JList) String() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteByte('[')
	for index, element := range l.ElementValues {
		if index != 0 {
			strBuilder.WriteString(", ")
		}

		strBuilder.WriteString(element.String())
	}
	strBuilder.WriteByte(']')

	return strBuilder.String()
}

func (l *JList) deepCopy() *JList {
	elementValues := make([]JValue, len(l.ElementValues))
	for index := range l.ElementValues {
		elementValue := l.ElementValues[index]
		if list, ok := elementValue.(*JList); ok {
			elementValues[index] = list.deepCopy()
		} else {
			elementValues[index] = l.ElementValues[index]
		}
	}

	return NewJList(elementValues)
}

func (l *JList) checkIndex(index int, arg JValue) error {
	if index < 0 || index >= len(l.ElementValues) {
		return errors.Wrap(&common.JRunTimeError{
			JError: &common.JError{
				StartPos: arg.GetStartPos(),
				EndPos:   arg.GetEndPos(),
			},
			Context: l.Context,
			Details: "index integer number must >= 0 and < length of list",
		}, "failed to index")
	}

	return nil
}
