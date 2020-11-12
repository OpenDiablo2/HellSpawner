package hsinput

// for converting our key/mod/mouse types to integers
func toInt(s []interface{}) []int {
	result := make([]int, len(s))

	for idx := range s {
		if n, ok := s[idx].(int); ok {
			result = append(result, n)
		}
	}

	return result
}
