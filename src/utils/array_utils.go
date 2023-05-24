package utils

func ArrayContains[k comparable](array []k, value k) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
