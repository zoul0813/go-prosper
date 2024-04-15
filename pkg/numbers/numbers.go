package numbers

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func Parse(word string, base int) (int, bool) {
	var num int = 0
	var length int = len(word)
	// unicode len of first char
	_, i := utf8.DecodeRuneInString(word)
	if strings.HasPrefix(word, "$") {
		// hex
		num = i
		base = 16
	} else if strings.HasPrefix(word, "0x") {
		// unicode len of second char
		_, j := utf8.DecodeRuneInString(word[i:])
		num = i + j
		base = 16
	} else if strings.HasPrefix(word, "&") || strings.HasPrefix(word, "#") {
		// decimal
		num = i
		base = 10
	} else if strings.HasPrefix(word, "%") {
		// binary
		num = i
		base = 2
	} else if strings.HasPrefix(word, "'") {
		// character
		base = 0
		if len(word) > 1 && len(word) <= 3 {
			char := word[1]
			return int(char), true
		} else {
			return 0, false
		}
	}

	if base > 0 {
		if i, err := strconv.ParseUint(word[num:length], base, strconv.IntSize); err == nil {
			return int(i), true
		} else {
			return 0, false
		}
	}
	return 0, false
}
