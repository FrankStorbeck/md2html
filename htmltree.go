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
// Code licensed under the BSD License: see licence.txt
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//

package main

import (
	"errors"
	"fmt"
	"html"
	"strconv"
	"strings"

	"source.storbeck.nl/md2html/branch"
)

const (
	noAlign = iota
	left
	right
	center
)

// TableInfo holds table data.
type TableInfo []int

// ListIndents is a struct holding data for the indentiation if list items.
type ListIndents struct {
	indent int            // number of spaces for the indentiation
	parent *branch.Branch // the parent branch
}

// HTMLTree is a struct for holding the data for the construction of a HTML
// tree.
type HTMLTree struct {
	blockQuoteLevel int              // level for blockQuote
	br              *branch.Branch   // current branch
	isHighLighted   bool             // true when text is high ligted
	isQuoted        bool             // true is the lines are precoded quotes
	lstIndnts       []int            // positions for indents for lists items
	lstParents      []*branch.Branch // parents of start list
	nextIndnt       int              // position for new indent
	parCount        int              // number if lines in paragraph
	root            *branch.Branch   // root branch
	sCount          int              // string number
	tblInfo         TableInfo        // table information
}

// NewHTMLTree returns a pointer to a new HTMLTree struct.
func NewHTMLTree(s string) HTMLTree {
	ht := HTMLTree{
		root: branch.NewBranch(s),
	}
	ht.br = ht.root
	return ht
}

// BlockQuote adds string 's' as a block quote. If it isn't a continuation of
// a bock quote, it will be initialized.
func (ht *HTMLTree) BlockQuote(s string) {
	lvl := 0
	for strings.Index(s, ">") == 0 {
		lvl++
		s = strings.TrimSpace(s[1:])
	}
	if lvl > ht.blockQuoteLevel {
		for i := ht.blockQuoteLevel; i < lvl; i++ {
			ht.br, _ = ht.br.AddBranch(-1, cBlockQuote)
		}
	} else if lvl < ht.blockQuoteLevel {
		ht.TryParent(ht.blockQuoteLevel - lvl)
	}
	ht.blockQuoteLevel = lvl

	if len(s) > 0 {
		ht.br.Add(-1, s)
	}
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
	s = Inline(Break(strings.TrimRight(s, "\n\r")))

	if l := len(s); l > 0 {
		nEnd := strings.Index(s[indnt:], ".") // end of number for ordered list
		if nEnd > 0 {
			_, err = strconv.Atoi(s[indnt : indnt+nEnd])
			if err != nil {
				err = nil
				nEnd = 0
			}
		}

		newTblInfo := TableInfo{}
		if len(ht.tblInfo) <= 0 {
			// could be table separator ...
			newTblInfo = CountTblColls(s)
		}

		switch {
		case ht.blockQuoteLevel <= 0 && l >= 3 && s[:3] == "```":
			// Syntactic hightlighting starts or ends
			return ht.HighLight()

		case ht.isHighLighted:
			ht.br.Add(-1, html.EscapeString(raw))
			return nil

		case s[0] == '>':
			ht.BlockQuote(s)
			return nil

		case ht.blockQuoteLevel > 0:
			ht.br.Add(-1, s)
			return nil

		case l > indnt && (s[indnt] == '*' ||
			(s[indnt] == '-' && !OnlyRunes(s, '-')) ||
			s[indnt] == '+'):
			// new unordered list item
			fallthrough

		case nEnd > 0:
			// new ordered list item
			return ht.ListItem(s, indnt, nEnd)

		case len(ht.tblInfo) > 0:
			// table row?
			return ht.TblRow(s, false)

		case len(newTblInfo) > 0:
			// previous line was a table header row
			return ht.ChangePrevToTblHdr(s, &newTblInfo)

		case len(ht.lstIndnts) > 0:
			if indnt > ht.lstIndnts[0] {
				n := ht.IndentIndex(indnt)
				if n < 0 {
					return errors.New("Negative index found")
				}
				if n < len(ht.lstIndnts)-1 {
					ht.PopIndents(n + 1)
					for _, s := range ht.br.Siblings() {
						if b, ok := s.(*branch.Branch); ok {
							ht.br = b // use the last branch
						}
					}
				}
			} else {
				if ht.TestLeadingHash(s) {
					return nil
				}
				ht.PopIndents(0)
				if ht.br.ID != cP {
					ht.br, _ = ht.br.AddBranch(-1, cP)
					ht.parCount = 1
				}
			}
			return ht.br.Add(-1, s[indnt:])
		}

		switch ht.parCount {
		case 0:
			if ht.TestLeadingHash(s) {
				return nil
			}
			if ht.br.ID != cP {
				ht.br, _ = ht.br.AddBranch(-1, cP)
			}
			ht.br.Add(-1, s)
			ht.parCount++

		case 1:
			switch {
			case OnlyRunes(s, '='):
				// previous line was a <h1> line
				fallthrough

			case OnlyRunes(s, '-'):
				// previous line was a <h2> line
				ht.ChangePrevToHdr(s)
				return nil
			}
			fallthrough

		default:
			ht.br.Add(-1, s)
			ht.parCount++
		}

	} else {
		// empty line
		switch {
		case ht.isHighLighted:
			ht.br.Add(-1, raw)
			return nil

		case len(ht.tblInfo) > 0:
			// end of table
			if err =  ht.TryParent(1); err != nil {
				return ht.MakeErr("Build", err)
			}
			ht.tblInfo = TableInfo{}
			ht.br, _ = ht.br.AddBranch(-1, cP)
			ht.parCount = 0

		case ht.blockQuoteLevel > 0:
			if err = ht.TryParent(ht.blockQuoteLevel); err != nil {
				return ht.MakeErr("Build", err)
			}
			ht.blockQuoteLevel = 0

		case len(ht.lstIndnts) > 0:
			ht.br.Add(-1, "<br/>")
			ht.parCount = 0

		case ht.parCount > 0:
			ht.Reset()
		}
	}
	return err
}

// ChangePrevToHdr changes the string that was added just before into a header
// line. If 's' contains '-' runes, it will be a level 2, otherwise a leve 1.
func (ht *HTMLTree) ChangePrevToHdr(s string) {
	prev, err := ht.br.Remove(-1)
	if err != nil {
		ht.br.Add(-1, s)
		return
	}
	if ps, ok := prev.(string); ok {
		lvl := 1
		if s[0] == '-' {
			lvl = 2
		}
		ht.Header(ps, lvl)
	} else {
		ht.br.Add(-1, prev) // restore 'prev'
	}
}

// ChangePrevToTblHdr changes the string that was added just before into a table
// header.
func (ht *HTMLTree) ChangePrevToTblHdr(s string, newTblInfo *TableInfo) error {
	ht.tblInfo = *newTblInfo
	ln, err := ht.br.Remove(-1)
	if err != nil {
		return ht.MakeErr("ChangePrevToTblHdr", err)
	}

	s, ok := ln.(string)
	if ok { // 'ln' must be a string
		ht.br, _ = ht.br.AddBranch(-1, cTable)
		ht.br.Info = "style=\"width: 100%\""
		b := TRow(s, true, &(ht.tblInfo))
		if b != nil {
			ht.br.Add(-1, b)
		} else {
			if err = ht.TryParent(1); err != nil {
				return ht.MakeErr("ChangePrevToTblHdr", err)
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

// HighLight highlights the text 's'.
func (ht *HTMLTree) HighLight() error {
	var err error
	// Syntactic highlighting starts or ends
	ht.isHighLighted = !ht.isHighLighted
	if ht.isHighLighted { // starts
		ht.br, _ = ht.br.AddBranch(-1, cPre)
		ht.br, _ = ht.br.AddBranch(-1, cCode)
	} else if err = ht.TryParent(2); err != nil {
			return ht.MakeErr("HighLight)", err)
	}
	return nil
}

// IndentIndex returns the index of the first element in the sorted list
// 'lstIndnts' for which all values are larger or equal to 'indnt'.
// If the 'ht.nextIndnt' is smaller or equal to 'indnt', the length of
// 'lstIndnts' will be returned. If all values are smaller than 'indnt', -1 will
// be returned.
func (ht *HTMLTree) IndentIndex(indnt int) int {
	l := len(ht.lstIndnts)
	if ht.nextIndnt <= indnt {
		return l
	}
	i := l - 1
	for i >= 0 {
		if ht.lstIndnts[i] <= indnt {
			return i
		}
		i--
	}
	return -1
}

// ListItem inserts a (un)ordered list item.
func (ht *HTMLTree) ListItem(s string, indnt, nEnd int) error {
	if nEnd < 0 {
		nEnd = 0
	}

	if len(ht.lstIndnts) <= 0 {
		ht.nextIndnt = indnt
	} else if err := ht.TryListParent(0); err != nil {
		return ht.MakeErr("ListItem", err)
	}

	nIndents := len(ht.lstIndnts)
	n := ht.IndentIndex(indnt)

	switch {
	case n >= nIndents:
		// start a new indent level
		ht.PushIndents(indnt + nEnd + CountLeading(s[indnt+nEnd+1:], ' ', -1) + 1)
		tg := cUl
		if nEnd > 0 {
			tg = cOl
		}
		ht.br, _ = ht.br.AddBranch(-1, tg)

	case n == nIndents-1:
		// continuation of current indent level

	default:
		// use older indent level
		ht.PopIndents(n + 1)
	}

	ht.br, _ = ht.br.AddBranch(-1, cLi)
	ht.br.Add(-1, strings.TrimSpace(s[indnt+nEnd+1:]))

	return nil
}

// MakeErr creates an error with some extra information about the code that
// generated it.
func (ht *HTMLTree) MakeErr(name string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("Error on line %d (detected by %s): %s", ht.sCount, name, err)
}

// PopIndents pops 'n' ListIndents from the ht.lstIndnts stack
func (ht *HTMLTree) PopIndents(n int) {
	if n < 0 {
		n = 0
	}
	if n >= len(ht.lstIndnts) {
		n = len(ht.lstIndnts) - 1
		if n < 0 {
			return
		}
	}
	ht.br = ht.lstParents[n]
	ht.nextIndnt = ht.lstIndnts[n]
	ht.lstIndnts = ht.lstIndnts[:n]
	ht.lstParents = ht.lstParents[:n]
}

// PushIndents pushes 'indnt' on the 'ht.lstIndnts' stack
func (ht *HTMLTree) PushIndents(indnt int) {
	ht.lstIndnts = append(ht.lstIndnts, ht.nextIndnt)
	ht.lstParents = append(ht.lstParents, ht.br)
	ht.nextIndnt = indnt
}

// Quote makes the lines show up as pre coded text 's'.
func (ht *HTMLTree) Quote(line string) error {
	if !ht.isQuoted {
		ht.br, _ = ht.br.AddBranch(-1, cPre)
		ht.br, _ = ht.br.AddBranch(-1, cCode)
		ht.isQuoted = true
	}
	ht.br.Add(-1, html.EscapeString(line[4:]))
	return nil
}

// Reset resets the current tree to a new paragraph in 'root'
func (ht *HTMLTree) Reset() {
	ht.br = ht.root
	ht.blockQuoteLevel = 0
	ht.isHighLighted = false
	ht.isQuoted = false
	ht.lstIndnts = nil
	ht.lstParents = nil
	ht.nextIndnt = 0
	ht.parCount = 0
}

// RmIfEmpty removes the branch 'brnch' if it is empty.
func (ht *HTMLTree) RmIfEmpty(brnch *branch.Branch) error {
	if brnch != ht.root && brnch.Len() <= 0 {
		if i, err := ht.br.Index(brnch); err == nil {
			if _, err = ht.br.Remove(i); err != nil {
				return ht.MakeErr("RmIfEmpty", err)
			}
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
		if err = ht.TryParent(1); err != nil {
			return ht.MakeErr("TblRow", err)
		}
		ht.tblInfo = TableInfo{}
		ht.br.Add(-1, s)
	}
	return nil
}

// TestLeadingHash tests for a <hn> string. If it is, a header line will be
// added.
func (ht *HTMLTree) TestLeadingHash(s string) bool {
	if leadingHash := CountLeading(s, '#', 6); leadingHash > 0 {
		// <h'leadingHash'> line
		ht.Header(s[leadingHash:], leadingHash)
		return true
	}
	return false
}

// Traverse traverses a slice of siblings. For each string found the function
// 'f' is called with this string as an argument. This sibling is then replaced
// by the result of 'f'. When the sibling is a pointer to a branch Traverse is
// called with a slice of all siblings in this branch as the argument.
func Traverse(sbs []interface{}, f func(string) []interface{}) []interface{} {
	r := []interface{}{}
	for _, in := range sbs {
		switch in.(type) {
		case string:
			r = append(r, f(in.(string))...)
		case *branch.Branch:
			b := in.(*branch.Branch)
			br := branch.NewBranch(b.ID)
			br.Info = b.Info
			br.Add(-1, Traverse(b.Siblings(), f)...)
			r = append(r, br)
		}
	}
	return r
}

// TryListParent tries to make the 'n'-th <ul> ot <ol> parent the active one.
func (ht *HTMLTree) TryListParent(n int) error {
	i := 0
	for n >= 0 {
		atList := ht.br.ID == cUl || ht.br.ID == cOl
		if atList {
			n--
		}
		if n >= 0 {
			var br *branch.Branch
			err := ht.TryParent(1)
			if atList {
				br = ht.br
				i++
			}
			if err != nil {
				return ht.MakeErr("TryListParent", err)
			}

			if ht.br == ht.root {
				// use the last parent of the "ul" or "ol" branch
				ht.br = br
				return nil
			}
		}
	}

	if i > 0 {
		l := len(ht.lstIndnts)
		ht.PopIndents(l - i)
	}
	return nil
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
		if ht.br, err = ht.br.Parent(1); err != nil {
			return ht.MakeErr("TryParent", err)
		}

		ht.RmIfEmpty(b)
		n--
	}
	return nil
}
