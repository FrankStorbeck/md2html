//
// htmltree.go
//
//  Questions and comments to:
//       <mailto:frank@foef.nl>
//
// definitions and functions for constructing an HTML tree.
//
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
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/FrankStorbeck/md2html/branch"
)

// TableInfo holds table data.
type TableInfo []int

// HTMLTree is a struct for holding the data for the construction of a HTML
// tree.
type HTMLTree struct {
	br          *branch.Branch // current branch
	inBlock     bool           // true while in blockQuote
	indents     []int          // positions for indents for lists items
	inList      bool           // true when in some (un)ordered list
	isHighLited bool           // true when text is high ligted
	isQuoted    bool           // true is the lines are precoded quotes
	sCount      int            // string number
	root        *branch.Branch // root branch
	tblInfo     TableInfo      // table information

}

const (
	noAlign = iota
	left
	right
	center
)

// BlockQuote adds string 's' as a block quote. If it isn't a continuation of
// a bock quote, it will be initialized.
func (ht *HTMLTree) BlockQuote(s string) error {
	var err error

	if !ht.inBlock {
		err = ht.TryParent(1)
		if err != nil {
			return err
		}
		ht.br, _ = ht.br.AddBranch(-1, "blockquote")
	}

	ht.inBlock = true

	if len(s) > 1 {
		ht.br.Add(-1, strings.TrimSpace(s[1:])+" ")
	}

	return nil
}

// Traverse traverses a slice of siblings. For each string found the function
// 'f' is called with this string as an argument. This sibling is then replaced
// by the result of 'f'. When the sibling is a pointer to a branch Traverse is
// called with a slice of all siblings in this branch as the argument.
func Traverse(intf []interface{}, f func(string) []interface{}) []interface{} {
	r := []interface{}{}
	for _, in := range intf {
		switch in.(type) {
		case string:
			r = append(r, f(in.(string))...)
		case *branch.Branch:
			br := in.(*branch.Branch)
			sblgs := br.Siblings()
			br.RemoveAll()
			br.Add(-1, Traverse(sblgs, f)...)
		}
	}
	return r
}

// Plain changes spaces in 's' to '-', puts everything in lower case and finally
// removes all tag info.
func Plain(s string) string {
	s = strings.ToLower(strings.Replace(s, " ", "-", -1))

	tags := []string{cCode, "strong", "em", "del"}
	for _, t := range tags {
		s = strings.Replace(s, "<"+t+">", "", -1)
		s = strings.Replace(s, "</"+t+">", "", -1)
	}

	return s
}

// Build reconstructs the HTML tree based on the contents of 's'.
func (ht *HTMLTree) Build(s string) error {
	raw := s
	var err error
	ht.sCount++
	indnt := CountLeading(s, ' ', -1)
	s = Inline(strings.Repeat(" ", indnt) + strings.TrimSpace(s))
	leadingHash := CountLeading(s, '#', 6)

	nEnd := strings.Index(s[indnt:], ".") // end of number for ordered list
	if nEnd > 0 {
		_, err = strconv.Atoi(s[indnt : indnt+nEnd])
		if err != nil {
			nEnd = 0
		}
	}

	if l := len(s); l > 0 {
		newTblInfo := TableInfo{}
		if len(ht.tblInfo) <= 0 {
			// can be a new table
			newTblInfo = CountTblColls(s)
		}

		switch {
		case l >= 3 && s[:3] == "```":
			// Syntactic hightlighting starts or ends
			return ht.HighLite(s)

		case ht.isHighLited:
			// Pre coded text
			ht.br.Add(-1, html.EscapeString(raw))

		case OnlyRunes(s, '='):
			// previous line was a <h1> line
			fallthrough

		case OnlyRunes(s, '-'):
			// previous line was a <h2> line
			ht.ChangePrevToHdr(s)

		case leadingHash > 0:
			// <h'leadingHash'> line
			ht.Header(s[leadingHash:], leadingHash)

		case l > 4 && indnt >= 4 && (ht.br.ID == "p" || ht.isQuoted):
			// pre coded quote
			err = ht.Quote(html.EscapeString(raw))

		case s[0] == '>':
			// block quote
			err = ht.BlockQuote(s)

		case l > indnt && (s[indnt] == '*' || s[indnt] == '-' || s[indnt] == '+'):
			// new unordered list item
			fallthrough

		case nEnd > 0:
			// new ordered list item
			err = ht.ListItem(s, indnt, nEnd)

		case ht.inList && indnt > ht.indents[0]:
			l := len(ht.indents)
			n := IndentIndex(indnt, ht.indents)
			if n < l-1 {
				err = ht.ListParent(l - 1 - n)
				if err != nil {
					ht.Reset()
				} else {
					ht.indents = ht.indents[:n+1]
				}
				ht.br, _ = ht.br.AddBranch(-1, "p")
			}
			ht.br.Add(-1, s[indnt:])

		case len(ht.tblInfo) > 0:
			// table row?
			err = ht.TblRow(s, false)

		case len(newTblInfo) > 0:
			// previous line was a table header row
			err = ht.ChangePrevToTblHdr(s, &newTblInfo)

		default:
			switch {
			case ht.inBlock:
				ht.inBlock = false
				err = ht.TryParent(1)
				if err != nil {
					return err
				}
				ht.br, _ = ht.br.AddBranch(-1, "p")

			case ht.isQuoted:
				err = ht.TryParent(2)
				if err != nil {
					return err
				}
				ht.isQuoted = false
				ht.br, _ = ht.br.AddBranch(-1, "p")

			case ht.inList:
				ht.Reset()
				s = s[indnt:]
			}

			ht.br.Add(-1, s)
		}
	} else {
		// empty line
		switch {
		case len(ht.tblInfo) > 0:
			// end of table
			err = ht.TryParent(1)
			ht.tblInfo = TableInfo{}

		case ht.isHighLited:
			ht.br.Add(-1, raw)

		case ht.inBlock:
			ht.inBlock = false
			err = ht.TryParent(1)
			if err != nil {
				return err
			}
			ht.br, _ = ht.br.AddBranch(-1, "p")

		case ht.inList:
			ht.br.Add(-1, "<p></p>")

		default:
			ht.TryParent(1)
			ht.br, _ = ht.br.AddBranch(-1, "p")
		}
	}
	return err
}

// ChangePrevToHdr changes the string that was added just before into a header
// line. If 's' contains '-' runes, it will be a level 2, otherwise a leve 1.
func (ht *HTMLTree) ChangePrevToHdr(s string) {
	prev, err := ht.br.Remove(-1)
	if err != nil {
		// Remove failed, add the string to a new <p>
		ht.br, _ = ht.root.AddBranch(-1, "p")
		ht.br.Add(-1, s)
	} else {
		ps, ok := prev.(string)
		if ok { // 'prev' must be a string
			lvl := 1
			if s[0] == '-' {
				lvl = 2
			}
			ht.Header(ps, lvl)
		} else {
			ht.br.Add(-1, prev) // restore 'prev'
		}
	}
}

// ChangePrevToTblHdr changes the string that was added just before into a table
// header.
func (ht *HTMLTree) ChangePrevToTblHdr(s string, newTblInfo *TableInfo) error {
	ht.tblInfo = *newTblInfo
	ln, err := ht.br.Remove(-1)
	if err != nil {
		return err
	}

	s, ok := ln.(string)
	if ok { // 'ln' must be a string
		// ht.br, _ = ht.br.Parent(1)
		ht.br, _ = ht.br.AddBranch(-1, "table")
		ht.br.Info = "style=\"width: 100%\""
		b := TRow(s, true, &(ht.tblInfo))
		if b != nil {
			ht.br.Add(-1, b)
		} else {
			err = ht.TryParent(1)
			if err != nil {
				return err
			}
			ht.br.Add(-1, s)
			ht.tblInfo = TableInfo{}
		}
	} else {
		// restore 'ln'
		ht.br.Add(-1, ln)
	}
	return nil
}

// CountTblColls test is a table separator line is found, and if true it
// analyses the table structure
func CountTblColls(s string) TableInfo {
	cols := strings.Split(strings.TrimSpace(s), "|")

	l := len(cols)
	if l < 3 || len(cols[0]) > 0 || len(cols[l-1]) > 0 {
		return TableInfo{}
	}

	tblInfo := make(TableInfo, len(cols)-2)
	cols = cols[1 : l-1] // trim first and last elements

	for i, c := range cols {
		c = strings.TrimSpace(c)
		l := len(c)
		if l > 0 {

			if c1 := strings.Index(c, ":"); c1 == 0 {
				if c2 := strings.Index(c[1:], ":"); c2 == l-2 {
					if !OnlyRunes(c[1:l-1], '-') {
						return TableInfo{}
					}
					tblInfo[i] = center
				} else {
					if !OnlyRunes(c[1:], '-') {
						return TableInfo{}
					}
					tblInfo[i] = left
				}
			} else if c1 == l-1 {
				if !OnlyRunes(c[:l-1], '-') {
					return TableInfo{}
				}
				tblInfo[i] = right
			} else {
				if !OnlyRunes(c, '-') {
					return TableInfo{}
				}
				// tblInfo[i] = noAlign
			}
		}
	}
	return tblInfo
}

// Header adds a header with level 'n' to the HTML tree.
func (ht *HTMLTree) Header(s string, n int) {
	b := ht.br
	ht.br = ht.root
	ht.RmIfEmpty(b)
	ht.br, _ = ht.br.AddBranch(-1, fmt.Sprintf("h%d", n))
	hdr := strings.TrimSpace(s)
	ht.br.Info = "id=\"" + Plain(hdr) + "\""
	ht.br.Add(-1, hdr)
	ht.Reset()
}

// NewHTMLTree returns a pointer to a new HTMLTree struct.
func NewHTMLTree(s string) HTMLTree {
	ht := HTMLTree{
		root: branch.NewBranch(s),
	}
	ht.br = ht.root
	return ht
}

// HighLite highlites the text 's'.
func (ht *HTMLTree) HighLite(s string) error {
	var err error
	// Syntactic hightlighting starts or ends
	ht.isHighLited = !ht.isHighLited
	if ht.isHighLited { // starts
		err = ht.TryParent(1)
		if err != nil {
			return err
		}
		ht.br, _ = ht.br.AddBranch(-1, "pre")
		ht.br, _ = ht.br.AddBranch(-1, "code")
	} else {
		err = ht.TryParent(2)
		if err != nil {
			return err
		}
		ht.br, _ = ht.br.AddBranch(-1, "p")
	}
	return nil
}

// IndentIndex returns the index of the first element in 'indents' for which
// the value is larger or equal to 'n'
func IndentIndex(n int, indents []int) int {
	l := len(indents) - 1
	i := 0
	for i <= l && n > indents[i] {
		i++
	}
	return i
}

// ListItem inserts a unordered list item.
func (ht *HTMLTree) ListItem(s string, indnt, nEnd int) error {
	if nEnd < 0 {
		nEnd = 0
	}
	err := ht.TryParent(1)
	if err != nil {
		return err
	}

	if ht.br.ID == "li" {
		// might be stil at '...ul{li{}}' 'ListItem' started at '...ul{li{p{...}}}'
		err = ht.TryParent(1)
		if err != nil {
			return err
		}
	}

	// nIndents holds the actual indents found so far; nIndents+1 hold the
	// position for an upcoming new indent.
	nIndents := len(ht.indents) - 1
	if !ht.inList || nIndents < 0 {
		ht.indents = []int{indnt}
		nIndents = 0
		ht.inList = true
	}

	n := IndentIndex(indnt, ht.indents)
	switch {
	case n >= nIndents:
		// start a new indent level
		inc := nEnd + CountLeading(s[indnt+nEnd+1:], ' ', -1)
		ht.indents = append(ht.indents, indnt+inc+1)
		tg := "ul"
		if nEnd > 0 {
			tg = "ol"
		}
		ht.br, _ = ht.br.AddBranch(-1, tg)

	case n == nIndents-1:
		// continuation of current indent level

	default:
		// use older indent level
		err = ht.ListParent(nIndents - n - 1)
		if err != nil {
			return err
		}
		ht.indents = ht.indents[:n+2]
	}

	ht.br, _ = ht.br.AddBranch(-1, "li")
	ht.br.Add(-1, strings.TrimSpace(s[indnt+nEnd+1:]))

	return nil
}

// ListParent set the current branch to the 'n'-th parent that is a 'ul{..}'
// or 'ol{...}'
func (ht *HTMLTree) ListParent(n int) error {
	for n > 0 {
		if ht.br.ID == "ul" || ht.br.ID == "ol" {
			n--
		}
		err := ht.TryParent(1)
		if err != nil {
			return err
		}
	}
	return nil
}

// Quote makes the lines show up as pre coded text 's'.
func (ht *HTMLTree) Quote(line string) error {
	if ht.br.ID == "p" {
		if err := ht.TryParent(1); err != nil {
			return err
		}
		ht.br, _ = ht.br.AddBranch(-1, "pre")
		ht.br, _ = ht.br.AddBranch(-1, "code")
		ht.isQuoted = true
	}
	ht.br.Add(-1, html.EscapeString(line[4:]))
	return nil
}

// Reset resets the current tree to a new paragraph in 'root'
func (ht *HTMLTree) Reset() {
	ht.br = ht.root
	ht.inBlock = false
	ht.indents = []int{}
	ht.inList = false
	ht.isHighLited = false
	ht.isQuoted = false
	ht.br, _ = ht.root.AddBranch(-1, "p")
}

// RmIfEmpty removes the branch 'brnch' if it is empty.
func (ht *HTMLTree) RmIfEmpty(brnch *branch.Branch) error {
	if brnch != ht.root && brnch.Len() <= 0 {
		if i, err := ht.br.Index(brnch); err == nil {
			_, err = ht.br.Remove(i)
			return err
		}
	}
	return nil
}

// TblRow adds a table row.
func (ht *HTMLTree) TblRow(s string, hdr bool) error {
	var err error
	b := TRow(s, false, &(ht.tblInfo))
	if b != nil {
		ht.br.Add(-1, b)
	} else {
		// end of table
		err = ht.TryParent(1)
		if err != nil {
			return err
		}
		ht.tblInfo = TableInfo{}
		ht.br.Add(-1, s)
	}
	return nil
}

// TRow returns a branch holding a table row. If hdr is true, a table header
// is asumed. tblInfo holds the TableInfo for the table collumns ,
func TRow(s string, hdr bool, tblInfo *TableInfo) *branch.Branch {
	cols := strings.Split(strings.TrimSpace(CodeUni(s, []byte{'|'}, true)), "|")

	l := len(cols)
	if l < 3 || len(cols[0]) > 0 || len(cols[l-1]) > 0 {
		return nil
	}

	cols = cols[1 : l-1]
	l = l - 2
	if l > len(*tblInfo) {
		l = len(*tblInfo)
		cols = cols[:l]
	}
	rslt := &branch.Branch{ID: "tr"}
	br := rslt

	tag := "td"
	if hdr {
		tag = "th"
	}

	for i, col := range cols {
		br, _ = br.AddBranch(-1, tag)

		if (*tblInfo)[i] > 0 {
			a := "left"
			switch (*tblInfo)[i] {
			case right:
				a = "right"
			case center:
				a = "center"
			}
			br.Info = fmt.Sprintf("style=\"text-align: %s\"", a)
		}

		br.Add(-1, strings.TrimSpace(CodeUni(col, []byte{'|'}, true)))
		br, _ = br.Parent(1)
	}
	return rslt
}

// TryParent goes safely to the 'n'-th parent of the current branch.
func (ht *HTMLTree) TryParent(n int) error {
	var err error
	if ht.br == nil {
		ht.Reset()
		return nil
	}
	for n > 0 {
		if ht.br == ht.root {
			return nil
		}
		b := ht.br
		ht.br, err = ht.br.Parent(1)
		if err != nil {
			err = fmt.Errorf("TryParent (%d): %s", ht.sCount, err)
			return err
		}

		ht.RmIfEmpty(b)
		n--
	}
	return nil
}
