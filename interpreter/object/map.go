package object

import (
	"fmt"
	"strings"
	"sync"

	"github.com/IfanTsai/jirachi/common"
)

const maxDeletion = 10000

type JMap struct {
	*JBaseValue
	ElementMap map[interface{}]JValue
	Deletion   int
	lock       *sync.RWMutex
}

func NewJMap(elementMap map[interface{}]JValue) *JMap {
	return &JMap{
		JBaseValue: &JBaseValue{},
		ElementMap: elementMap,
		Deletion:   0,
		lock:       &sync.RWMutex{},
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
	m.lock.RLock()
	copyMap := NewJMap(m.ElementMap)
	m.lock.RUnlock()

	copyMap.Deletion = m.Deletion
	copyMap.lock = m.lock

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

	m.lock.RLock()
	defer m.lock.RUnlock()

	if resValue, ok := m.ElementMap[arg.GetValue()]; ok {
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

	m.lock.Lock()
	if _, ok := indexValue.(*JNull); ok {
		delete(m.ElementMap, indexArg.GetValue())
		m.Deletion++

		if m.Deletion >= maxDeletion {
			newElementMap := make(map[interface{}]JValue, len(m.ElementMap))
			for key, value := range m.ElementMap {
				newElementMap[key] = value
			}

			m.ElementMap = newElementMap
			m.Deletion = 0
		}
	} else {
		m.ElementMap[indexArg.GetValue()] = indexValue
	}
	m.lock.Unlock()

	return m, nil
}

func (m *JMap) String() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteByte('{')
	firstKey := true

	m.lock.RLock()
	for key, value := range m.ElementMap {
		if !firstKey {
			strBuilder.WriteString(", ")
		}

		firstKey = false
		strBuilder.WriteString(fmt.Sprintf("%v", key))
		strBuilder.WriteString(": ")
		strBuilder.WriteString(value.String())
	}
	m.lock.RUnlock()

	strBuilder.WriteByte('}')

	return strBuilder.String()
}
