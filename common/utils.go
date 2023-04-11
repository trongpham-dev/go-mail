package common

func GetSubstring(s string, indices []int) string {
	return string(s[indices[0]:indices[1]])
}