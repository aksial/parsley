package parsley

import (
	"unicode"
)

func IsSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func IsNumber(r rune) bool {
	return unicode.IsNumber(r)
}

func IsEOF(item Item) bool {
	return item.Type == "eof"
}
