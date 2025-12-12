package iter

// Contains checks if the iterator contains the specified value.
// It returns true if the value is found, false otherwise.
func Contains[T comparable](in Iter[T], v T) bool {
	for {
		item, ok := in.next()
		if !ok {
			break
		}
		if item == v {
			return true
		}
	}
	return false
}

// Find searches for the first element in the iterator that satisfies the given predicate.
// It returns the index of the found element, the element itself, and a boolean indicating whether such an element was found.
func Find[T any](in Iter[T], pred func(T) bool) (int, T, bool) {
	index := 0
	for {
		item, ok := in.next()
		if !ok {
			break
		}
		if pred(item) {
			return index, item, true
		}
		index++
	}
	var zero T
	return -1, zero, false
}

// FindMap searches for the first element in the iterator that satisfies the given predicate.
// It returns the transformed element and a boolean indicating whether such an element was found.
func FindMap[T, R any](in Iter[T], f func(T) (R, bool)) (R, bool) {
	for {
		item, ok := in.next()
		if !ok {
			break
		}
		if result, found := f(item); found {
			return result, true
		}
	}
	var zero R
	return zero, false
}

// Any checks if any element in the iterator satisfies the given predicate.
// It returns true if at least one element satisfies the predicate, false otherwise.
func Any[T any](in Iter[T], pred func(T) bool) bool {
	for {
		item, ok := in.next()
		if !ok {
			break
		}
		if pred(item) {
			return true
		}
	}
	return false
}

// All checks if all elements in the iterator satisfy the given predicate.
// It returns true if all elements satisfy the predicate, false otherwise.
func All[T any](in Iter[T], pred func(T) bool) bool {
	for {
		item, ok := in.next()
		if !ok {
			break
		}
		if !pred(item) {
			return false
		}
	}
	return true
}
