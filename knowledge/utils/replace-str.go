package utils

import (
	"regexp"
)

func ReplaceStr(re *regexp.Regexp, s string, repl string, n int) string {

	count := 0
	if n == 0 {
		return s
	}

	if n < 0 {
		return re.ReplaceAllString(s, repl)
	}
	return re.ReplaceAllStringFunc(s, func(s string) string {
		count++
		if count > n {
			return s
		}
		return repl
	})
}
