package utils

//StrSub 字符串截取
func StrSub(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)
	if length == 0 {
		return str
	}

	if start < 0 {
		start = 0
	} else if start >= length {
		return ""
	}

	if start == end {
		return ""
	} else if end <= start || end > length {
		end = length
	}

	return string(rs[start:end])
}
