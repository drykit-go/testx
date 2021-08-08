package strcase

import (
	"unicode"
)

func Camel(in string) string {
	cc := []rune{}
	ok := false
	for i, r := range in {
		switch {
		case r == '_':
			continue
		case i == len(in)-1:
			doAppend(&cc, r, !ok)
		case ok:
			doAppend(&cc, r, false)
		case unicode.IsLower(rune(in[i+1])):
			doAppend(&cc, r, i == 0)
			ok = true
		default:
			doAppend(&cc, r, true)
		}
	}
	return string(cc)
}

func doAppend(dst *[]rune, r rune, lower bool) {
	if lower {
		r = unicode.ToLower(r)
	}
	*dst = append(*dst, r)
}
