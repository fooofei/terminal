package runes

import (
	"unicode/utf8"
)

// DecodeRune will decode p bytes to runes array of utf8 format
// return the runes array and decoded size
func DecodeRune(p []byte) ([]rune, int) {
	acceptSize := 0
	rn := make([]rune, 0)
	for len(p) > 0 {
		r, size := utf8.DecodeRune(p)
		if r == utf8.RuneError {
			break
		}
		acceptSize += size
		rn = append(rn, r)
		p = p[size:]
	}
	return rn, acceptSize
}

// DecodeRuneOnNewLine will decode bytes to runes
// call the newLine callback when got a new line runes
// return tail runes array and the decoded valid size
func DecodeRuneOnNewLine(p []byte, newLine func([]rune)) ([]rune, int) {
	acceptSize := 0
	rn := make([]rune, 0)
	for len(p) > 0 {
		r, size := utf8.DecodeRune(p)
		if r == utf8.RuneError {
			break
		}
		acceptSize += size
		p = p[size:]
		if r == '\n' {
			newLine(rn)
			rn = rn[:0]
			continue
		}
		rn = append(rn, r)
	}
	return rn, acceptSize
}
