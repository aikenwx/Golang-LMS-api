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

func RemoveDuplicatesInStringSlice(slice []string) []string {

	// create a map of strings
	keys := make(map[string]bool)
	result := []string{}

	// add all strings to map if not already in map
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}

func RemoveAllStringsInSlice(slice []string, stringsToRemove []string) []string {
	result := []string{}

	// create a map of strings to remove
	stringsToRemoveMap := make(map[string]bool)
	for _, str := range stringsToRemove {
		stringsToRemoveMap[str] = true
	}

	// add all that are not in the map
	for _, str := range slice {
		if _, ok := stringsToRemoveMap[str]; !ok {
			result = append(result, str)
		}
	}

	return result
}
