package util

// Set represents a set of strings.
type Set map[any]struct{}

// Add adds an element to the set.
func (s Set) Add(item any) {
	s[item] = struct{}{}
}

// Contains checks if an element is in the set.
func (s Set) Contains(item any) bool {
	_, ok := s[item]
	return ok
}

// Remove removes an element from the set.
func (s Set) Remove(item any) {
	delete(s, item)
}

func (s Set) GetSlice() []any {
	var result = make([]any, len(s))
	for key := range s {
		result = append(result, key)
	}
	return result
}

// NewStringSet creates and returns a new empty Set.
func NewStringSet() Set {
	return make(Set)
}
