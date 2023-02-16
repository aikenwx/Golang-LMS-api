package helpers

func Map[T, V any](slice []T, function func(T) V) []V {
	result := make([]V, len(slice))
	for i, t := range slice {
		result[i] = function(t)
	}
	return result
}

func Filter[T any](slice []T, function func(T) bool) []T {
	result := make([]T, 0)
	for _, t := range slice {
		if function(t) {
			result = append(result, t)
		}
	}
	return result
}
