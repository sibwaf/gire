package util

func Coalesce[T comparable](values ...T) T {
	var defaultValue T

	for _, value := range values {
		if value != defaultValue {
			return value
		}
	}

	return defaultValue
}
