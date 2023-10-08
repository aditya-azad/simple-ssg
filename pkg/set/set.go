package set

type Set[K comparable] struct {
	mp map[K]bool
}

func New[K comparable](vals ...K) *Set[K] {
	st := &Set[K]{map[K]bool{}}
	for _, val := range vals {
		st.mp[val] = true
	}
	return st
}

func (s *Set[K]) Add(val K) {
	s.mp[val] = true
}

func (s *Set[K]) Has(val K) bool {
	_, ok := s.mp[val]
	return ok
}
