package set

var exist = struct{}{}

type Set struct {
	m map[interface{}]struct{}
}

func NewSet(items ...interface{}) *Set {
	_set := &Set{
		m: make(map[interface{}]struct{}),
	}

	for _, item := range items {
		_set.Add(item)
	}

	return _set
}

func (s *Set) Add(item interface{}) {
	s.m[item] = exist
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]

	return ok
}

func (s *Set) Remove(item interface{}) {
	delete(s.m, item)
}

func (s *Set) Clear() {
	if !s.IsEmpty() {
		s.m = make(map[interface{}]struct{})
	}
}

func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}
