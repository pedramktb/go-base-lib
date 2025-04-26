package types

type Optional[T any] struct {
	Val T
	Set bool
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{Val: value, Set: true}
}
