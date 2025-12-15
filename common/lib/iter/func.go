package iter

// Compose composes two functions f and g into a single function.
// The resulting function takes an argument of type A, applies f to it to get a value of type B,
// and then applies g to that value to get a final result of type C.
func Compose[A, B, C any](f func(A) B, g func(B) C) func(A) C {
	return func(a A) C {
		return g(f(a))
	}
}

// Pipe applies a series of functions to an initial value in sequence.
// It takes an initial value of type T and a variadic list of functions that each take and return a value of type T.
// The result is the final value after all functions have been applied.
func Pipe[T any](v T, fns ...func(T) T) T {
	for _, f := range fns {
		v = f(v)
	}
	return v
}
