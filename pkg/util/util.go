package util

func FilterSlice[T any](slices []T, f func(T) bool) []T {
	var out []T
	for _, slice := range slices {
		if f(slice) {
			out = append(out, slice)
		}
	}
	return out
}

func FindSlice[T any](slices []T, f func(T) bool) (T, bool) {
	for _, slice := range slices {
		if f(slice) {
			return slice, true
		}
	}
	var result T
	return result, false
}
