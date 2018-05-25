//
// tools.go
//
//  Questions and comments to:
//       <mailto:frank@foef.nl>
//
// tools for md2html.go.
//
// © 2018 Frank Storbeck

package main

import (
	"bytes"
	"fmt"
	"html"
	"strings"
)

// OnlyRunes tests if a string of at least three runes 'rn' and no other runes.
func OnlyRunes(s string, rn rune) bool {
	if len(s) < 3 {
		return false
	}
	for _, r := range s {
		if r != rn {
			return false
		}
	}
	return true
}

// CountLeading returns the number of leading runes 'rn' in string 's'. 'm'
// is the maximum that is alowed. When 'm' is less than 0, unlimited runes are
// alowed.
func CountLeading(s string, rn rune, m int) int {
	n := 0
	for _, r := range s {
		switch r {
		case rn:
			n++
			if m >= 0 && n > m {
				return 0
			}
		default:
			return n
		}
	}
	return 0
}

// Images translates mark down image definitions to their html equivalents
func Images(s string) string {
	l := len(s)
	if i := strings.Index(s, "!["); i >= 0 && l > i+5 {
		if j := strings.Index(s[i:], "]"); j > 0 && l > (i+j+2) {
			if s[i+j+1] == '(' {
				if k := strings.Index(s[i+j+1:], ")"); k > 0 {
					s = s[:i] + "<img src=\"" + s[i+j+2:i+j+k+1] + "\" alt=\"" +
						s[i+2:i+j] + "\"/>" + Images(s[i+j+k+2:])
				}
			}
		}
	}
	return s
}

// Inline translates all inline mark down definitions
// to their html equivalents
func Inline(s string) string {
	s = StrongEmDel(s)
	s = Images(s)
	s = Links(s) // always after Images(s)!
	s = InlineCodes(s)
	return s
}

// InlineCodes translates mark down code definitions to their html equivalents
func InlineCodes(s string) string {
	l := len(s)
	if i := strings.Index(s, "`"); i >= 0 && l > i+2 {
		if j := strings.Index(s[i+1:], "`"); j > 0 {
			s = s[:i] + "<code>" + html.EscapeString(s[i+1:i+j+1]) +
				"</code>" + InlineCodes(s[i+j+2:])
		}
	}
	return s
}

// Links translates mark down link definitions to their html equivalents
func Links(s string) string {
	l := len(s)
	if i := strings.Index(s, "["); i >= 0 && l > i+4 {
		if j := strings.Index(s[i:], "]"); j > 0 && l > (i+j+1) {
			if s[i+j+1] == '(' {
				if k := strings.Index(s[i+j+1:], ")"); k > 0 {
					s = s[:i] + "<a href=\"" + s[i+j+2:i+j+k+1] + "\">" + s[i+1:i+j] +
						"</a>" + Links(s[i+j+k+2:])
				}
			}
		}
	}
	return s
}

// StrongEmDel translates mark down strong, emphasis and deleted definitions
// to their html equivalents
func StrongEmDel(s string) string {
	tgs := []struct {
		tg  string
		sep []string
	}{
		{"strong", []string{"**", "__"}},
		{"em", []string{"_", "*"}},
		{"del", []string{"~~", "~~"}},
	}

	seps := make([]byte, len(tgs))
	for i, tg := range tgs {
		seps[i] = tg.sep[0][0]
	}
	s = UniCode(s, seps)

	for _, t := range tgs {
		for _, sp := range t.sep {
			var subs []string
			switch {
			// special case: "^\*[^\*]+.*$" will become an unordered list, no italics!
			case sp == "*" && len(s) > 0 && s[:1] == sp:
				subs = strings.Split(s[1:], sp)
				s = sp + subs[0]
			default:
				subs = strings.Split(s, sp)
				s = subs[0]
			}
			n := len(subs)
			for i := 1; i < n; i = i + 2 {
				if i+1 < n {
					s = s + "<" + t.tg + ">" + subs[i] + "</" + t.tg + ">" +
						subs[i+1]
				} else {
					s = s + sp + subs[i]
				}
			}
		}
	}

	return UnEscape(s, seps)
}

// UnEscape replaces unicode by its (non escaped) character.
func UnEscape(s string, esc []byte) string {
	b := []byte(s)

	for _, sp := range esc {
		old := []byte(fmt.Sprintf("U+%04X", sp))
		new := []byte{sp}
		b = bytes.Replace(b, old, new, -1)
	}

	return string(b)
}

// UniCode replaces escaped characters by its uni code. A character is escaped
// by putting a '\' in front of it.
func UniCode(s string, esc []byte) string {
	b := []byte(s)
	for _, sp := range esc {
		old := []byte{'\\', sp}
		new := []byte(fmt.Sprintf("U+%04X", sp))
		b = bytes.Replace(b, old, new, -1)
	}
	return string(b)
}
