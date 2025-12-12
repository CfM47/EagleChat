package iter

// Fold applies a binary function to each element of the iterator,
// accumulating a single result starting from initial value as parameter.
func Fold[T, R any](in Iter[T], init R, f func(R, T) R) R {
	result := init
	for {
		item, ok := in.next()
		if !ok {
			break
		}
		result = f(result, item)
	}
	return result
}

// Reduce applies a binary function to each element of the iterator,
// reducing the iterator to a single value. It returns the reduced value
// and a boolean indicating whether the reduction was successful (i.e., the iterator was not empty).
func Reduce[T any](in Iter[T], f func(T, T) T) (T, bool) {
	first, ok := in.next()
	if !ok {
		var zero T
		return zero, false
	}
	result := first
	for {
		item, ok := in.next()
		if !ok {
			break
		}
		result = f(result, item)
	}
	return result, true
}

// Scan applies a binary function to each element of the iterator,
// producing a new iterator of the intermediate accumulated results.
// It starts from the initial value provided as a parameter.
func Scan[T, R any](in Iter[T], init R, f func(R, T) R) Iter[R] {
	result := init
	return Iter[R]{
		next: func() (R, bool) {
			item, ok := in.next()
			if !ok {
				var zero R
				return zero, false
			}
			result = f(result, item)
			return result, true
		},
	}
}
