package iter

// Number is a constraint that matches all numeric types.
type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr |
		float32 | float64
}

// Sum computes the sum of all elements in the input iterator.
func Sum[T Number](in Iter[T]) T {
	return Fold(in, 0, func(a T, b T) T { return a + b })
}

// Product computes the product of all elements in the input iterator.
func Product[T Number](in Iter[T]) T {
	return Fold(in, 1, func(a T, b T) T { return a * b })
}

// Min finds the minimum element in the input iterator.
func Min[T Number](in Iter[T]) (T, bool) {
	return Reduce(in, func(a T, b T) T { return min(a, b) })
}

// Max finds the maximum element in the input iterator.
func Max[T Number](in Iter[T]) (T, bool) {
	return Reduce(in, func(a T, b T) T { return max(a, b) })
}
