package util

// InSlice 是否存在于数组中
func InSlice(key string, array []string) bool {
	for _, v := range array {
		if key == v {
			return true
		}
	}
	return false
}

func InSliceInt(key int, array []int) bool {
	for _, v := range array {
		if key == v {
			return true
		}
	}
	return false
}
