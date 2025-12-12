package iter

// Map applies the given function to each element of the input slice
// and returns a new slice containing the results.
func Map[T, R any](input Iter[T], fn func(T) R) Iter[R] {
	return Iter[R]{
		next: func() (R, bool) {
			value, ok := input.next()
			if !ok {
				var zero R
				return zero, false
			}
			return fn(value), true
		},
	}
}

// Filter returns a new slice containing only the elements of the input slice
// for which the given function returns true.
func Filter[T any](input Iter[T], fn func(T) bool) Iter[T] {
	return Iter[T]{
		next: func() (T, bool) {
			for {
				value, ok := input.next()
				if !ok {
					var zero T
					return zero, false
				}
				if fn(value) {
					return value, true
				}
			}
		},
	}
}

// FilterMap applies the given function to each element of the input iter.
// If the function returns a second value of true, the first value is included
// in the output slice. Otherwise, it is excluded.
func FilterMap[T, R any](input Iter[T], fn func(T) (R, bool)) Iter[R] {
	return Iter[R]{
		next: func() (R, bool) {
			for {
				value, ok := input.next()
				if !ok {
					var zero R
					return zero, false
				}
				if mapped, valid := fn(value); ok && valid {
					return mapped, true
				}
			}
		},
	}
}

// FlatMap applies the given function to each element of the input iter,
// which returns an iterator. It then flattens the results into a single iterator.
func FlatMap[T, R any](input Iter[T], fn func(T) Iter[R]) Iter[R] {
	var currentInnerIter Iter[R]
	var innerIterActive bool

	return Iter[R]{
		next: func() (R, bool) {
			for {
				if innerIterActive {
					mappedValue, mappedOk := currentInnerIter.next()
					if mappedOk {
						return mappedValue, true
					}
					innerIterActive = false // Current inner iterator is exhausted
				}

				// Get the next value from the outer iterator
				value, ok := input.next()
				if !ok {
					var zero R
					return zero, false // Outer iterator is exhausted
				}

				// Create a new inner iterator
				currentInnerIter = fn(value)
				innerIterActive = true
			}
		},
	}
}
