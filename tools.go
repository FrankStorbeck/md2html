//
// tools.go
//
//  Questions and comments to:
//       <mailto:frank@foef.nl>
//
// tools for md2html.go.
//
// Copyright © 2018 Frank Storbeck. All rights reserved.
// Code licensed under the BSD License:
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"bytes"
	"fmt"
	"html"
	"strings"
	"unicode"
)

// Break adds a break when the line ends with two or more spaces. It also slices
// off all trailing spaces.
func Break(s string) string {
	l := len(s)
	brk := ""
	if i := strings.LastIndex(s, "  "); i >= 0 && i >= l-2 {
		brk = "<br/>"
		s = s[:l-2]
	}
	return strings.TrimRightFunc(s, unicode.IsSpace) + brk
}

// CodeUni replaces all runes given in 'runes' by its uni code. when 'esc' is
// true, only escaped runes are replaced. A rune is escaped by putting a '\' in
// front of it.
func CodeUni(s string, runes []byte, esc bool) string {
	b := []byte(s)
	for _, r := range runes {
		var old []byte
		if esc {
			old = []byte{'\\', r}
		} else {
			old = []byte{r}
		}
		new := []byte(fmt.Sprintf("U+%04X", r))
		b = bytes.Replace(b, old, new, -1)
	}
	return string(b)
}

// DecodeUni replaces unicoded runes a by the rune itself as given in 'runes'.
// When 'esc' is true, the rune will be escaped by placing a '\' char in front
// of it.
func DecodeUni(s string, runes []byte, esc bool) string {
	b := []byte(s)

	for _, r := range runes {
		old := []byte(fmt.Sprintf("U+%04X", r))
		var new []byte
		if esc {
			new = []byte{'\\', r}
		} else {
			new = []byte{r}
		}
		b = bytes.Replace(b, old, new, -1)
	}

	return string(b)
}

// OnlyRunes tests if a string consists of one or more runes 'rn' and no
// other runes.
func OnlyRunes(s string, rn rune) bool {
	if len(s) < 1 {
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
					uc := CodeUni(s[i+j+2:i+j+k+1], []byte{'*', '_', '~'}, false)
					s = s[:i] + "<img src=\"" + uc + "\" alt=\"" +
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
	// order is important here
	s = InlineCodes(s)
	s = Images(s)
	s = Links(s)
	s = StrongEmDel(s)
	return DecodeUni(s, []byte{'*', '_', '~'}, false)
}

// InlineCodes translates mark down code definitions to their html equivalents
func InlineCodes(s string) string {
	l := len(s)
	if i := strings.Index(s, "`"); i >= 0 && l > i+2 {
		if j := strings.Index(s[i+1:], "`"); j > 0 {
			uc := CodeUni(s[i+1:i+j+1], []byte{'*', '_', '~'}, false)
			s = s[:i] + "<code>" + html.EscapeString(uc) +
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
					uc := CodeUni(s[i+j+2:i+j+k+1], []byte{'*', '_', '~'}, false)
					s = s[:i] + "<a href=\"" + uc + "\">" + s[i+1:i+j] +
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
	s = CodeUni(s, seps, true)

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

	return DecodeUni(s, seps, false)
}
