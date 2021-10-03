package morestring

import "strings"

func ReplaceLast(text string, toReplace string, replaceWith string) string {
	i := strings.LastIndex(text, toReplace)
	excludingLast := text[:i] + strings.Replace(text[i:], toReplace, replaceWith, 1)
	return excludingLast
}
