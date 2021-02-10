package hsutil

// IntToBool converts 32-bit intager into bool
// if input is equal to 0, then returns false, else - true
func IntToBool(i int32) bool {
	if i >= 1 {
		return true
	}

	return false
}

// BoolToInt converts bool into 32-bit intager
// if b is true, then returns 1, else 0
func BoolToInt(b bool) int32 {
	if b {
		return 1
	}

	return 0
}
