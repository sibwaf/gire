package util

type Set[T comparable] interface {
	Add(item T) bool
	Remove(item T) bool
	Contains(item T) bool
	ToSlice() []T
}

type set[T comparable] struct {
	data map[T]bool
}

func NewSet[T comparable]() Set[T] {
	result := set[T]{data: make(map[T]bool)}
	return &result
}

func (s set[T]) Add(item T) bool {
	_, exists := s.data[item]
	if exists {
		return false
	}

	s.data[item] = true
	return true
}

func (s set[T]) Remove(item T) bool {
	_, exists := s.data[item]
	if exists {
		delete(s.data, item)
	}

	return exists
}

func (s set[T]) Contains(item T) bool {
	_, exists := s.data[item]
	return exists
}

func (s set[T]) ToSlice() []T {
	result := make([]T, 0, len(s.data))
	for item := range s.data {
		result = append(result, item)
	}
	return result
}
