package util

import (
	"strings"
	"unicode"
)

func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}

		out = append(out, unicode.ToLower(runes[i]))
	}

	return strings.Replace(string(out), "-", "", -1)
}

func BoolToOnOff(on bool) string {
	if on {
		return "on"
	}

	return "off"
}
