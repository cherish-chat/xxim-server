package utils

func AnyPtr[T any](v T) *T {
	return &v
}
