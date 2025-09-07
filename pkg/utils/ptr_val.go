package utils

func Ptr[T any](v T) *T {
	return &v
}

func Val[T any](v *T) T {
	var zero T
	if v != nil {
		return *v
	}

	return zero
}
