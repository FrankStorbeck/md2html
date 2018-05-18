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
