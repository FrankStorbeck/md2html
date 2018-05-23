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
	"html"
	"strings"

	"source.storbeck.nl/md2html/branch"
)

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
}

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

// Build reconstructs the HTML tree based on the contents of 's'.
func (ht *HTMLTree) Build(s string) error {
	raw := s
	var err error
	ht.sCount++
	indnt := CountLeading(s, ' ', -1)
	s = strings.Repeat(" ", indnt) + strings.TrimSpace(s)
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

		case l >= 3 && s[:3] == "```":
			// Syntactic hightlighting starts or ends
			return ht.HighLite(s)

		case ht.isHighLited:
			// Pre coded text
			ht.br.Add(-1, html.EscapeString(raw))

		case l > 4 && indnt >= 4 && (ht.br.ID == "p" || ht.isQuoted):
			// pre coded quote
			err = ht.Quote(html.EscapeString(raw))

		case s[0] == '>':
			// block Quote
			err = ht.BlockQuote(Inline(s))

		case l > indnt && (s[indnt] == '*' || s[indnt] == '-' || s[indnt] == '+'):
			// new unordered list item
			err = ht.ListItem(s, indnt)

		case ht.inList && indnt > ht.indents[0]:
			l := len(ht.indents)
			n := IndentIndex(indnt, ht.indents)
			if n < l-1 {
				err = ht.ListParent(l - 1 - n)
				if err != nil {
					ht.inList = false
					ht.br = ht.root
				} else {
					ht.indents = ht.indents[:n+2]
				}
			}
			ht.br.Add(-1, Inline(s[indnt:]))

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
				ht.inList = false
				ht.br, _ = ht.root.AddBranch(-1, "p")
				s = s[indnt:]
			}

			ht.br.Add(-1, Inline(s))
		}
	} else {
		// empty line
		switch {
		case ht.inBlock:
			ht.br.Add(-1, "<br>")

		case ht.inList:
			// in a list when there is a blank line, siblings are put into a paragraph
			if ht.br.ID == "li" {
				sblngs := ht.br.Siblings()
				ht.br.RemoveAll()
				ht.br, _ = ht.br.AddBranch(-1, "p")
				for _, s := range sblngs {
					ht.br.Add(-1, s)
				}
			}

		case ht.br.ID == "p":
			ht.br, _ = ht.br.AddBranch(-1, "p")

		default:
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
func (ht *HTMLTree) ListItem(s string, indnt int) error {
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
		indntInc := CountLeading(s[indnt+1:], ' ', -1)
		ht.indents = append(ht.indents, indnt+indntInc+1)
		ht.br, _ = ht.br.AddBranch(-1, "ul")

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
	ht.br.Add(-1, strings.TrimSpace(s[indnt+1:]))

	return nil
}

// ListParent set the current branch to the 'n'-th parent that is a 'ul{..}'
func (ht *HTMLTree) ListParent(n int) error {
	for n > 0 {
		if ht.br.ID == "ul" {
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
		ht.br, _ = ht.br.AddBranch(-1, "pre")
		ht.br, _ = ht.br.AddBranch(-1, "code")
		ht.isQuoted = true
	}
	ht.br.Add(-1, html.EscapeString(line[4:]))
	return nil
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

// TryParent goes safely to the 'n'-th parent of the current branch.
func (ht *HTMLTree) TryParent(n int) error {
	var err error
	if ht.br == nil {
		ht.br = ht.root
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
