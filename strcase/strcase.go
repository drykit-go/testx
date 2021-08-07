package strcase

import (
	"unicode"
)

func Camel(in string) string {
	cc := []rune{}
	ok := false
	for i, r := range in {
		if r == '_' {
			continue
		}

		if isLast := i == len(in)-1; isLast {
			if ok {
				cc = append(cc, r)
			} else {
				cc = append(cc, unicode.ToLower(r))
			}
			break
		}

		if ok {
			cc = append(cc, r)
			continue
		}

		if next := rune(in[i+1]); unicode.IsLower(next) {
			if i == 0 {
				cc = append(cc, unicode.ToLower(r))
			} else {
				cc = append(cc, r)
			}
			ok = true
			continue
		}

		cc = append(cc, unicode.ToLower(r))
	}
	return string(cc)
}

// func Camel2(in string) string {
// 	for i, r := range in {
// 	}
// }
