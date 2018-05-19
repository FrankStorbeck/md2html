//
// branch.go
//
//  Questions and comments to:
//       <mailto:frank@foef.nl>
//
// branch implements a branch data stucture
//
// © 2018 Frank Storbeck

package branch

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoSiblings     = errors.New("no siblings")
	ErrNoSiblingFound = errors.New("no sibling found")
)

// Branch holds the data for a branch originating from a parent branch.
type Branch struct {
	ID       string        // identifier for the branch
	Info     string        // some optional description
	parent   *Branch       // parent for wich this one is a sibling
	siblings []interface{} // its siblings
}

// Add adds 'sls' to the slice of siblings just before sibling 'n', or, if 'n'
// is less than zero, after the last sibling. When an error occured, it will be
// returned.
func (br *Branch) Add(n int, sls ...interface{}) error {
	var err error
	n, err = br.testIndx(n)
	if err != nil {
		return err
	}

	lbrs := len(br.siblings)
	lsls := len(sls)
	l := lbrs + lsls
	if n < 0 {
		n = lbrs
	}

	siblings := make([]interface{}, l)
	for i := 0; i < lbrs; i++ {
		if i < n {
			siblings[i] = br.siblings[i]
		} else {
			siblings[i+lsls] = br.siblings[i]
		}
	}
	for i, s := range sls {
		siblings[n+i] = s
		if b, ok := s.(*Branch); ok {
			b.parent = br
		}
	}

	br.siblings = siblings
	return nil
}

// AddBranch adds a new branch with identifier 'id' just before the 'n'-th
// sibling, or, if 'n' is less than zero, after the last sibling. A
// pointer to the inserted branch will be returned. When an error occured, nil
// and the error will be returned.
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

// Len returns the number of siblings for the branch.
func (br *Branch) Len() int {
	return len(br.siblings)
}

func escape(s string) string {
	r := strings.Replace(s, "{", "\\{", -1)
	r = strings.Replace(r, "}", "\\}", -1)
	return r
}

// NewBranch returns a pointer to a new Branch struct with identifier 's'.
func NewBranch(s string) *Branch {
	return &Branch{ID: s, siblings: make([]interface{}, 0)}
}

// Parent returns a pointer to the 'n'-th parent of branch. If 'n' is less than
// one a pointer to itself will be returned. When an error occured, nil and the
// error will be returned.
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

// ParentById returns a pointer to the first parent with an ID equal to 'id'.
// When an error occured, nil and the error will be returned.
func (br *Branch) ParentById(id string) (*Branch, error) {
	for br.ID != id {
		var err error
		br, err = br.Parent(1)
		if err != nil {
			return nil, err
		}
	}
	return br, nil
}

// Index returns the index of 'intf' in the slice with siblings. If it cannot
// be found, it returns -1 and an error.
func (br *Branch) Index(intf interface{}) (int, error) {
	for n, sblg := range br.siblings {
		if sblg == intf {
			return n, nil
		}
	}
	return -1, ErrNoSiblingFound
}

// Remove removes the 'n'-th sibling, or, if 'n' is less than zero, the last
// sibling. It returns a pointer to the removed sibling. When an error occured,
// nil and the error will be returned.
func (br *Branch) Remove(n int) (interface{}, error) {

	l := len(br.siblings)
	if l <= 0 {
		return nil, ErrNoSiblings
	}

	if n < 0 {
		sblg := br.siblings[l-1]
		br.siblings = br.siblings[:l-1]
		return sblg, nil
	}

	rslt := make([]interface{}, l-1)
	var sblg interface{}
	j := 0
	for i := 0; i < l; i++ {
		if i != n {
			rslt[j] = br.siblings[i]
			j++
		} else {
			sblg = br.siblings[i]
		}
	}

	br.siblings = rslt
	return sblg, nil
}

// RemoveAll removes all siblings
func (br *Branch) RemoveAll() {
	br.siblings = []interface{}{}
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

// String returns a string holding the tree in the following format:
// ID[:Info]{siblings...}
func (br *Branch) String() string {
	s := escape(br.ID)
	if len(br.Info) > 0 {
		s = s + ":" + strings.TrimSpace(escape(br.Info))
	}

	s = s + "{"
	space := ""
	for _, c := range br.Siblings() {
		switch k := c.(type) {
		case *Branch:
			s = s + space + k.String()
		case string:
			s = s + space + escape(k)
		default:
			s = s + fmt.Sprint(k)
		}
		space = " "
	}

	return s + "}"
}

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

// TreePath returns a string holding the tree path ending at 'b'. The path has
// the format "/root.ID/.../br.ID".
func TreePath(b *Branch) string {
	s := ""
	for b != nil {
		s = b.ID + "/" + s
		b = b.parent
	}
	return "/" + s
}
