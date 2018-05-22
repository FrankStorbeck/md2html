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
	r := s
	l := len(r)
	if i := strings.Index(r, "!["); i >= 0 && l > i+5 {
		if j := strings.Index(r[i:], "]"); j > 0 && l > (i+j+2) {
			if r[i+j+1] == '(' {
				if k := strings.Index(r[i+j+1:], ")"); k > 0 {
					r = r[:i] + "<img src=\"" + r[i+j+2:i+j+k+1] + "\" alt=\"" +
						r[i+2:i+j] + "\">" + Images(r[i+j+k+2:])
				}
			}
		}
	}
	return r
}

// Inline translates all inline mark down definitions
// to their html equivalents
func Inline(s string) string {
	s = StrongEmDel(s)
	s = Images(s)
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
	r := UniCode(s, seps)

	for _, t := range tgs {
		for _, sp := range t.sep {
			var subs []string
			switch {
			// special case: "^\*[^\*]+.*$" will become an unordered list, no italics!
			case sp == "*" && len(r) > 0 && r[:1] == sp:
				subs = strings.Split(r[1:], sp)
				r = sp + subs[0]
			default:
				subs = strings.Split(r, sp)
				r = subs[0]
			}
			n := len(subs)
			for i := 1; i < n; i = i + 2 {
				if i+1 < n {
					r = r + "<" + t.tg + ">" + subs[i] + "</" + t.tg + ">" +
						subs[i+1]
				} else {
					r = r + sp + subs[i]
				}
			}
		}
	}

	return UnEscape(r, seps)
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
