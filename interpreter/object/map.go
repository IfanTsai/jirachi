package object

import (
	"fmt"
	"strings"

	"github.com/IfanTsai/jirachi/common"
	"github.com/IfanTsai/jirachi/pkg/safemap"
)

type JMap struct {
	*JBaseValue
	ElementMap *safemap.SafeMap[JValue]
}

func NewJMap(elementMap *safemap.SafeMap[JValue]) *JMap {
	return &JMap{
		JBaseValue: &JBaseValue{},
		ElementMap: elementMap,
	}
}

func (m *JMap) SetJPos(startPos, endPos *common.JPosition) JValue {
	m.StartPos = startPos
	m.EndPos = endPos

	return m
}

func (m *JMap) SetJContext(context *common.JContext) JValue {
	m.Context = context

	return m
}

func (m *JMap) Copy() JValue {
	copyMap := NewJMap(m.ElementMap)

	return copyMap
}

func (m *JMap) IndexAccess(arg JValue) (JValue, error) {
	if !CanHashed(arg) {
		return nil, &common.JRunTimeError{
			JError: &common.JError{
				StartPos: m.StartPos,
				EndPos:   m.EndPos,
			},
			Context: m.Context,
			Details: "Cannot hashed",
		}
	}

	if resValue, ok := m.ElementMap.Get(arg.GetValue()); ok {
		return resValue, nil
	}

	return NewJNull(), nil
}

func (m *JMap) IndexAssign(indexArg, indexValue JValue) (JValue, error) {
	if !CanHashed(indexArg) {
		return nil, &common.JRunTimeError{
			JError: &common.JError{
				StartPos: m.StartPos,
				EndPos:   m.EndPos,
			},
			Context: m.Context,
			Details: "Cannot hashed",
		}
	}

	if _, ok := indexValue.(*JNull); ok {
		m.ElementMap.Del(indexArg.GetValue())
	} else {
		m.ElementMap.Set(indexArg.GetValue(), indexValue)
	}

	return m, nil
}

func (m *JMap) String() string {

	strBuilder := strings.Builder{}
	strBuilder.WriteByte('{')
	firstKey := true

	m.ElementMap.Range(func(key any, value JValue) bool {
		if !firstKey {
			strBuilder.WriteString(", ")
		}

		firstKey = false
		strBuilder.WriteString(fmt.Sprintf("%v", key))
		strBuilder.WriteString(": ")
		strBuilder.WriteString(value.String())

		return true
	})

	strBuilder.WriteByte('}')

	return strBuilder.String()
}
