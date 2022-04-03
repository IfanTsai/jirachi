package safemap

import "sync"

const (
	copyThreshold = 1000
	maxDeletion   = 10000
)

// SafeMap provides a map alternative to avoid memory leak.
// This implementation is not needed until issue below fixed.
// https://github.com/golang/go/issues/20135
type SafeMap[V any] struct {
	lock        sync.RWMutex
	deletionOld int
	deletionNew int
	dirtyOld    map[any]V
	dirtyNew    map[any]V
}

// NewSafeMap returns a SafeMap.
func NewSafeMap[V any]() *SafeMap[V] {
	return &SafeMap[V]{
		dirtyOld: make(map[any]V),
		dirtyNew: make(map[any]V),
	}
}

// Del deletes the object with the given key from m.
func (m *SafeMap[V]) Del(key any) {
	m.lock.Lock()

	if _, ok := m.dirtyOld[key]; ok {
		delete(m.dirtyOld, key)
		m.deletionOld++
	} else if _, ok := m.dirtyNew[key]; ok {
		delete(m.dirtyNew, key)
		m.deletionNew++
	}

	if m.deletionOld >= maxDeletion && len(m.dirtyOld) < copyThreshold {
		for k, v := range m.dirtyOld {
			m.dirtyNew[k] = v
		}
		m.dirtyOld = m.dirtyNew
		m.deletionOld = m.deletionNew
		m.dirtyNew = make(map[any]V)
		m.deletionNew = 0
	}

	if m.deletionNew >= maxDeletion && len(m.dirtyNew) < copyThreshold {
		for k, v := range m.dirtyNew {
			m.dirtyOld[k] = v
		}

		m.dirtyNew = make(map[any]V)
		m.deletionNew = 0
	}

	m.lock.Unlock()
}

// Get gets the object with the given key from m.
func (m *SafeMap[V]) Get(key any) (V, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if val, ok := m.dirtyOld[key]; ok {
		return val, true
	}

	val, ok := m.dirtyNew[key]

	return val, ok
}

// Set sets the object into m with the given key.
func (m *SafeMap[V]) Set(key any, value V) {
	m.lock.Lock()

	if m.deletionOld <= maxDeletion {
		if _, ok := m.dirtyNew[key]; ok {
			delete(m.dirtyNew, key)
			m.deletionNew++
		}

		m.dirtyOld[key] = value
	} else {
		if _, ok := m.dirtyOld[key]; ok {
			delete(m.dirtyOld, key)
			m.deletionOld++
		}

		m.dirtyNew[key] = value
	}

	m.lock.Unlock()
}

// Size returns the size of m.
func (m *SafeMap[V]) Size() int {
	m.lock.RLock()

	size := len(m.dirtyOld) + len(m.dirtyNew)

	m.lock.RUnlock()

	return size
}

// Range calls f sequentially for each key and value present in the map.
func (m *SafeMap[V]) Range(f func(key any, value V) bool) {
	validMap := m.dirtyNew
	if len(m.dirtyOld) > 0 {
		validMap = m.dirtyOld
	}

	for k, v := range validMap {
		if !f(k, v) {
			break
		}
	}
}
