//
// branch.go
//
//  Questions and comments to:
//       <mailto:frank@foef.nl>
//
// branch implements a tree data stucture
//
// Copyright © 2018 Frank Storbeck. All rights reserved.
// Code licensed under the BSD License: see licence.txt
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

// Package branch can be used storing abitrary data in tree data stucture.
// It starts with a base branch that can hold a set of leaves and other
// branches. Togeter they are denoted as siblings. Each branch knows about the
// branch that holds it as one of its siblings, exept of course the base branch.
package branch

import (
	"errors"
	"fmt"
	"strings"
)

// Error codes
var (
	ErrNoSiblings     = errors.New("No siblings")
	ErrNoSiblingFound = errors.New("No sibling found")
)

// Branch holds the data for a branch.
type Branch struct {
	ID       string        // identifier for this branch
	Info     string        // some optional description
	parent   *Branch       // parent for wich this branch is a sibling
	siblings []interface{} // its siblings
}

// NewBranch returns a pointer to a new Branch struct with identifier id.
func NewBranch(id string) *Branch {
	return &Branch{ID: id, siblings: make([]interface{}, 0)}
}

// Add adds a number of new siblings. If n is non negative, they will appended
// to the currently present siblings. Otherwise they will be added just before
// the n-th sibling.
func (br *Branch) Add(n int, new ...interface{}) error {
	var err error
	n, err = br.testIndx(n)
	if err != nil {
		return err
	}

	lCur := len(br.siblings)
	lNew := len(new)
	l := lCur + lNew
	if n < 0 {
		n = lCur
	}

	sblngs := make([]interface{}, l)
	for i := 0; i < lCur; i++ {
		if i < n {
			sblngs[i] = br.siblings[i]
		} else {
			sblngs[i+lNew] = br.siblings[i]
		}
	}
	for i, s := range new {
		sblngs[n+i] = s
		if b, ok := s.(*Branch); ok {
			b.parent = br
		}
	}

	br.siblings = sblngs
	return nil
}

// AddBranch adds a new branch with identifier id. If n is negative, it will
// appended to the currently present siblings. Otherwise it will be added just
// before the n-th sibling. If successfull, a pointer to the inserted branch
// will be returned.
func (br *Branch) AddBranch(n int, id string) (*Branch, error) {
	var err error
	n, err = br.testIndx(n)
	if err != nil {
		return nil, err
	}
	b := &Branch{ID: id, siblings: make([]interface{}, 0)}
	br.Add(n, b)
	return b, nil
}

// Len returns the number of siblings.
func (br *Branch) Len() int {
	return len(br.siblings)
}

// Parent returns a pointer to the n-th parent of the branch. If n is non
// positive a pointer to itself will be returned.
func (br *Branch) Parent(n int) (*Branch, error) {
	p := br
	for n > 0 {
		if p = p.parent; p == nil {
			return nil, errors.New("missing parent")
		}
		n--
	}
	return p, nil
}

// ParentID returns a pointer to itself or the first parent with an ID equal to
// id.
func (br *Branch) ParentID(id string) (*Branch, error) {
	b := br
	for b.ID != id {
		var err error
		b, err = b.Parent(1)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

// Index returns the index of sblng. If it cannot be found it returns -1 and
// ErrNoSiblingFound.
func (br *Branch) Index(sblng interface{}) (int, error) {
	for n, sb := range br.siblings {
		if sb == sblng {
			return n, nil
		}
	}
	return -1, ErrNoSiblingFound
}

// Remove removes the n-th sibling. If n is negative, the last sibling will be
// removed. It returns a pointer to the removed sibling.
func (br *Branch) Remove(n int) (interface{}, error) {
	l := len(br.siblings)
	if l <= 0 {
		return nil, ErrNoSiblings
	}

	if n < 0 {
		sb := br.siblings[l-1]
		br.siblings = br.siblings[:l-1]
		return sb, nil
	}

	sblngs := make([]interface{}, l-1)
	var sb interface{}
	j := 0
	for i := 0; i < l; i++ {
		if i != n {
			sblngs[j] = br.siblings[i]
			j++
		} else {
			sb = br.siblings[i]
		}
	}

	br.siblings = sblngs
	return sb, nil
}

// RemoveAll removes all siblings
func (br *Branch) RemoveAll() {
	br.siblings = make([]interface{}, 0)
}

// Siblings returns a slice with all the siblings for this branch.
func (br *Branch) Siblings() []interface{} {
	return br.siblings
}

// SiblingN returns the n-th sibling of the branch. If n is less than zero, the
// last sibling will be returned. When an error occured, nil and the error will
// be returned.
func (br *Branch) SiblingN(n int) (interface{}, error) {
	l := len(br.siblings)
	if l <= 0 {
		return nil, ErrNoSiblings
	}
	if n < 0 {
		n = l - 1
	}
	return br.siblings[n], nil
}

// String returns a string showing the contents of the branch.
func (br *Branch) String() string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s", escape(escape(br.ID, true), false))
	if len(br.Info) > 0 {
		fmt.Fprintf(&sb, "(%s)", strings.TrimSpace(escape(br.Info, true)))
	}

	fmt.Fprintf(&sb, "{")
	space := ""
	for _, c := range br.Siblings() {
		switch k := c.(type) {
		case *Branch:
			fmt.Fprintf(&sb, "%s%s", space, k.String())
		case string:
			fmt.Fprintf(&sb, "%s\"%s\"", space, escape(k, false))
		default:
			fmt.Fprintf(&sb, "%s", k)
		}
		space = " "
	}
	fmt.Fprintf(&sb, "}")

	return sb.String()
}

// TreePath returns a string holding the tree path ending at b. The path has
// the format "/root.ID/.../br.ID".
func (br *Branch) TreePath() string {
	ids := []string{}
	// s := ""
	for br != nil {
		// s = br.ID + "/" + s
		ids = append(ids, br.ID)
		br = br.parent
	}
	sbr := strings.Builder{}
	for i := len(ids) - 1; i > 0; i-- {
		fmt.Fprintf(&sbr, "/%s", ids[i])
	}
	// return "/" + s
	return sbr.String()
}

// escape returns a string in wich some characters in s will be escaped.
func escape(s string, brackets bool) string {
	if brackets {
		s = strings.Replace(s, "(", "\\(", -1)
		s = strings.Replace(s, ")", "\\)", -1)
	} else {
		s = strings.Replace(s, "{", "\\{", -1)
		s = strings.Replace(s, "}", "\\}", -1)
		s = strings.Replace(s, "\"", "\\\"", -1)
	}
	return s
}

// testIndx tests if n is a valid index for the slice holding the siblings.
func (br *Branch) testIndx(n int) (int, error) {
	l := len(br.siblings)
	if l <= 0 {
		e := ErrNoSiblings
		if n < 0 {
			e = nil
		}
		return -1, e
	}
	if n >= l {
		return n, fmt.Errorf("sibling index (%d) out of range (0-%d)", n, l-1)
	}
	if n < 0 {
		n = -1
	}
	return n, nil
}
