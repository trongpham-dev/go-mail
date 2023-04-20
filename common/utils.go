package common

func GetSubstring(s string, indices []int) string {
	return string(s[indices[0]:indices[1]])
}

func IsExist(o, s []uint32) bool {
	for _, v := range o {
		for _, v1 := range s {
			if v == v1 {
				return true
			}
		}
	}
	return false
}
