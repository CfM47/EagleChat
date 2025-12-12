package iter

// Iter represents an iterator over a sequence of values of type T.
// It is a function that returns the next value and a boolean indicating whether there are more values to iterate over.
type Iter[T any] struct {
	next func() (T, bool)
}

// FromSlice creates an iterator from a slice of type T.
func FromSlice[T any](s []T) Iter[T] {
	index := 0
	return Iter[T]{
		next: func() (T, bool) {
			if index >= len(s) {
				var zero T
				return zero, false
			}
			value := s[index]
			index++
			return value, true
		},
	}
}

// ToSlice collects all values from the iterator into a slice of type T.
func ToSlice[T any](it Iter[T]) []T {
	var result []T
	for {
		value, ok := it.next()
		if !ok {
			break
		}
		result = append(result, value)
	}
	return result
}
