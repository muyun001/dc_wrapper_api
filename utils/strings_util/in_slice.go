package strings_util

import "strings"

// InSlice 判断字符串是否在切片中
func InSlice(arr []string, findMe string) bool {
	for i := range arr {
		if strings.Compare(arr[i], findMe) == 0 {
			return true
		}
	}

	return false
}
