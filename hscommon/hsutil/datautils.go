package hsutil

// BoolToInt converts bool into 32-bit intager
// if b is true, then returns 1, else 0
func BoolToInt(b bool) int32 {
	if b {
		return 1
	}

	return 0
}

// Wrap integer to max: wrap(450, 360) == 90
func Wrap(x, max int) int {
	wrapped := x % max

	if wrapped < 0 {
		return max + wrapped
	}

	return wrapped
}
