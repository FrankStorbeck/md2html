//
// htmltree.go
//
//  Questions and comments to:
//       <mailto:frank@foef.nl>
//
// definitions and functions for constructing an HTML tree.
//
// © 2018 Frank Storbeck

package main

import (
	"fmt"
	"strings"

	"source.storbeck.nl/md2html2/branch"
)

// HTMLTree is a struct for holding the data for the construction of a HTML
// tree.
type HTMLTree struct {
	br   *branch.Branch // current branch
	root *branch.Branch // root branch
}

// Build reconstructs the HTML tree based on the contents of 's'.
func (ht *HTMLTree) Build(s string) error {
	var err error
	leadingHash := CountLeading(s, '#', 6)

	if l := len(s); l > 0 {
		switch {
		case OnlyRunes(s, '='):
			// previous line was a <h1> line
			fallthrough

		case OnlyRunes(s, '-'):
			// previous line was a <h2> line
			ht.ChangePrevToHdr(s)

		case leadingHash > 0:
			// <h'leadingHash'> line
			ht.Header(s[leadingHash:], leadingHash)

		default:
			ht.br.Add(-1, s)
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

// Header adds a header line with level 'n' to the HTML tree.
func (ht *HTMLTree) Header(line string, n int) {
	b := ht.br
	ht.br = ht.root
	ht.RmIfEmpty(b)
	ht.br, _ = ht.br.AddBranch(-1, fmt.Sprintf("h%d", n))
	ht.br.Add(-1, strings.TrimSpace(line))
	ht.br, _ = ht.root.AddBranch(-1, "p")
}

// NewHTMLTree returns a pointer to a new HTMLTree struct.
func NewHTMLTree(s string) HTMLTree {
	ht := HTMLTree{
		root: branch.NewBranch(s),
	}
	ht.br = ht.root
	return ht
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
