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

	if l := len(s); l > 0 {
		switch {
		case OnlyRunes(s, '='):
			// previous line was a <h1> line
			ht.ChangePrevToHdr(s)

		default:
			ht.br.Add(-1, s)
		}
	}
	return err
}

// ChangePrevToHdr changes the string that was added just before into a header
// line.
func (ht *HTMLTree) ChangePrevToHdr(line string) {
	prev, err := ht.br.Remove(-1)
	if err != nil {
		// Remove failed, add the line to a new <p>
		ht.br, _ = ht.root.AddBranch(-1, "p")
		ht.br.Add(-1, line)
	} else {
		s, ok := prev.(string)
		if ok { // 'prev' must be a string
			ht.Header(s, 1)
		} else {
			ht.br.Add(-1, prev) // restore 'prev'
		}
	}
}

// Header adds a header line with level 'n' to the HTML tree.
func (ht *HTMLTree) Header(line string, n int) {
	ht.br = ht.root
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
