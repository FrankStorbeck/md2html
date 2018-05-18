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

// esc2uni replaces escaped characters by its uni code. A character is escaped
// by putting a '\' in front of it.
func esc2uni(s string, esc []byte) string {
	b := []byte(s)
	for _, sp := range esc {
		old := []byte{'\\', sp}
		new := []byte(fmt.Sprintf("U+%04X", sp))
		b = bytes.Replace(b, old, new, -1)
	}
	return string(b)
}

// uni2esc replaces unicode by its (non escaped) character.
func uni2esc(s string, esc []byte) string {
	b := []byte(s)

	for _, sp := range esc {
		old := []byte(fmt.Sprintf("U+%04X", sp))
		new := []byte{sp}
		b = bytes.Replace(b, old, new, -1)
	}

	return string(b)
}
