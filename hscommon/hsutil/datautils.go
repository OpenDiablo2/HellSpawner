package hsutil

// BoolToInt converts bool into 32-bit intager
// if b is true, then returns 1, else 0
func BoolToInt(b bool) int32 {
	if b {
		return 1
	}

	return 0
}
