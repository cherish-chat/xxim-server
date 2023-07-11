package utils

type xString struct {
}

var String = &xString{}

// FirstUtf8 获取字符串的第一个utf8字符，如果字符串为空，则返回”
// 如果字符串不是ascii编码，则返回第一个utf8字符，如果是ascii编码，则返回第一个字符
func (x *xString) FirstUtf8(s string) string {
	if len(s) == 0 {
		return ""
	}
	if s[0] < 128 {
		return s[0:1]
	} else if len(s) >= 3 {
		return s[0:3]
	} else {
		return s
	}
}

func (x *xString) Utf8Split(s string, length int) string {
	if len(s) == 0 {
		return ""
	}
	var result string
	for i := 0; i < length; i++ {
		first := x.FirstUtf8(s)
		if first == "" {
			break
		}
		result += first
		//去掉第一个字符
		s = s[len(first):]
	}
	return result
}
