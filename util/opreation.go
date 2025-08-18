package util

// Comparable defines an interface for types that can be compared, including string and int.
type Comparable interface {
	~string | ~int
}

// InSlice checks if a given value is present in a slice and returns true if found, otherwise false.
func InSlice[T Comparable](val T, slice []T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
