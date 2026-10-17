package utils_value

type Optional[T any] struct {
	value   T
	present bool
}

func (x Optional[T]) Get() (T, bool) {
	return x.value, x.present
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{
		value:   value,
		present: true,
	}
}
