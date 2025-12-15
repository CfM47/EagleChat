package lib

func Values[K comparable, V any](m map[K]V) []V {
	answ := make([]V, 0, len(m))
	for _, v := range m {
		answ = append(answ, v)
	}

	return answ
}

func Keys[K comparable, V any](m map[K]V) []K {
	answ := make([]K, 0, len(m))
	for k := range m {
		answ = append(answ, k)
	}

	return answ
}

func MakeSet[K comparable](slice []K) map[K]struct{} {
	answ := make(map[K]struct{})
	for _, v := range slice {
		answ[v] = struct{}{}
	}
	return answ
}

func MapKeys[K comparable, V any, K2 comparable](m map[K]V, f func(K) K2) map[K2]V {
	answ := make(map[K2]V)

	for k, v := range m {
		answ[f(k)] = v
	}

	return answ
}
