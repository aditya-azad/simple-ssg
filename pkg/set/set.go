package set

type stringSet struct {
	mp map[string]bool
}

func NewStringSet(vals ...string) *stringSet {
	st := &stringSet{map[string]bool{}}
	for _, val := range vals {
		st.mp[val] = true
	}
	return st
}

func (s *stringSet) Add(val string) {
	s.mp[val] = true
}

func (s *stringSet) Has(val string) bool {
	_, ok := s.mp[val]
	return ok
}
