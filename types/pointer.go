package types

// Pointer for inline usage (e.g. output of a function)
func Pointer[T any](value T) *T {
	return &value
}

// Safe pointer deference
func Derefer[T any](pointer *T) T {
	if pointer == nil {
		return *new(T)
	}
	return *pointer
}
