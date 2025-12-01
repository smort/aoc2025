package util

func Intersect[T comparable](slice1 []T, slice2 []T) []T {
	result := make([]T, 0)
	hash := make(map[T]bool)
	for _, e := range slice1 {
		hash[e] = false
	}
	for _, e := range slice2 {
		if isAdded, exists := hash[e]; exists && !isAdded {
			result = append(result, e)
			hash[e] = true
		}
	}

	return result
}

func HasIntersection[T comparable](slice1 []T, slice2 []T) bool {
	return len(Intersect(slice1, slice2)) > 0
}
