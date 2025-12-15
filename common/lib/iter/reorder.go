package iter

type Tuple[A, B any] struct {
	First  A
	Second B
}

// Zip combines two iterators into a single iterator of tuples.
// It returns an iterator that returns tuples of the corresponding values from both iterators.
func Zip[A, B any](a Iter[A], b Iter[B]) Iter[Tuple[A, B]] {
	return Iter[Tuple[A, B]]{
		next: func() (Tuple[A, B], bool) {
			aItem, aOk := a.next()
			bItem, bOk := b.next()
			if !aOk || !bOk {
				var zero Tuple[A, B]
				return zero, false
			}
			return Tuple[A, B]{First: aItem, Second: bItem}, true
		},
	}
}

// ZipWith applies a function to corresponding elements from two slices and returns a new slice containing the results.
func ZipWith[A, B, R any](a []A, b []B, f func(A, B) R) []R {
	minLen := min(len(b), len(a))
	output := make([]R, minLen)
	for i := range minLen {
		output[i] = f(a[i], b[i])
	}
	return output
}

// Unzip takes a slice of tuples and returns two slices: the first slice contains the first elements of the tuples,
// and the second slice contains the second elements of the tuples.
func Unzip[A, B any](in []Tuple[A, B]) ([]A, []B) {
	a := make([]A, len(in))
	b := make([]B, len(in))
	for i, v := range in {
		a[i] = v.First
		b[i] = v.Second
	}
	return a, b
}

// Partition splits a slice into two slices based on a predicate function.
// The first returned slice contains elements that satisfy the predicate,
// while the second contains elements that do not.
func Partition[T any](in []T, pred func(T) bool) ([]T, []T) {
	truePart := make([]T, 0)
	falsePart := make([]T, 0)
	for _, v := range in {
		if pred(v) {
			truePart = append(truePart, v)
		} else {
			falsePart = append(falsePart, v)
		}
	}
	return truePart, falsePart
}
