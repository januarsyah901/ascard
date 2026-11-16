package articles

import (
	"strings"
	"unicode/utf16"

	"github.com/januarsyah901/ascard/api/internal/httpx"
)

// Limits from pickArticle in server/routes/articles.js, in JavaScript string length units.
const (
	maxTitleLen    = 200
	maxContentLen  = 100_000
	maxCategoryLen = 50
)

// validate mirrors pickArticle: it returns ok=false whenever Express answered
// 400 "Title and content are required".
func validate(f httpx.Fields) (Input, bool) {
	title, ok := f.String("title")
	if !ok || jsTrim(title) == "" || jsLen(title) > maxTitleLen {
		return Input{}, false
	}
	content, ok := f.String("content")
	if !ok || jsTrim(content) == "" || jsLen(content) > maxContentLen {
		return Input{}, false
	}
	in := Input{Title: jsTrim(title), Content: content, Published: f.IsTrue("published")}
	if c, ok := f.String("category"); ok {
		if c = jsTrim(c); c != "" {
			c = jsCut(c, maxCategoryLen)
			in.Category = &c
		}
	}
	return in, true
}

// jsSpace reports whether r is removed by JavaScript's String.prototype.trim
// (WhiteSpace + LineTerminator). It differs from unicode.IsSpace on U+0085 and U+FEFF.
func jsSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x00A0, 0x1680, 0x2028, 0x2029, 0x202F, 0x205F, 0x3000, 0xFEFF:
		return true
	}
	return r >= 0x2000 && r <= 0x200A
}

func jsTrim(s string) string { return strings.TrimFunc(s, jsSpace) }

// jsLen is the JavaScript length of s (UTF-16 code units).
func jsLen(s string) int {
	n := 0
	for _, r := range s {
		n += utf16.RuneLen(r)
	}
	return n
}

// jsCut keeps at most max UTF-16 code units of s, like slice(0, max), but never splits
// a surrogate pair (a lone surrogate cannot be stored in PostgreSQL text anyway).
func jsCut(s string, max int) string {
	n := 0
	for i, r := range s {
		w := utf16.RuneLen(r)
		if n+w > max {
			return s[:i]
		}
		n += w
	}
	return s
}
