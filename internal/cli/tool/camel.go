package tool

import (
	"strings"
	"unicode"
)

func ToCamelCaseWithFirstUpper(s string) string {
	return camelCase(s, true)
}

func ToCamelCaseWithFirstLower(s string) string {
	return camelCase(s, false)
}

func camelCase(s string, firstLetterUpper bool) string {
	s = strings.TrimFunc(s, func(r rune) bool {
		return isDelimiter(r) || unicode.IsSpace(r)
	})

	var builder strings.Builder

	builder.Grow(len(s))

	var prev rune

	for _, cur := range s {
		if !isDelimiter(cur) {
			if isDelimiter(prev) || (firstLetterUpper && prev == 0) {
				builder.WriteRune(unicode.ToUpper(cur))
			} else {
				builder.WriteRune(unicode.ToLower(cur))
			}
		}

		prev = cur
	}

	return builder.String()
}

func isDelimiter(symbol rune) bool {
	return symbol == '-' || symbol == '_' || unicode.IsSpace(symbol)
}
