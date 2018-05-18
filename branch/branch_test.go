//
// branch_test.go
//
//  Questions and comments to:
//       <mailto:frank@foef.nl>
//
// branch_test does some testing for branch.go
//
// © 2018 Frank Storbeck

package branch

import (
	"fmt"
	"testing"
)

func TestBranch(t *testing.T) {
	root := &Branch{ID: "root"}
	want := "root{}"
	if s := root.String(); s != want {
		t.Errorf("In TestBranch root.String() returns: %q, should be: %q", s, want)
	}

	p, err := root.Parent(0)
	if err != nil {
		t.Errorf("In TestBranch root.Parent(0) reports an error: %q, should be nil",
			err)
	} else if s := p.String(); s != want {
		t.Errorf("In TestBranch: root.Parent(0).String() returns: %q, should be %q",
			s, want)
	}
}

func TestAdd(t *testing.T) {
	tests1 := []struct {
		n       int
		sibling string
		want    string
	}{
		{0, "s2", fmt.Sprintf("%s", ErrNoSiblings)},
		{-1, "s2", "root{s2}"},
		{0, "s0", "root{s0 s2}"},
		{1, "s1", "root{s0 s1 s2}"},
		{-1, "{s3}", "root{s0 s1 s2 \\{s3\\}}"},
		{4, "s3", "sibling index (4) out of range (0-3)"},
	}

	root := &Branch{ID: "root"}

	for _, tst := range tests1 {
		err := root.Add(tst.n, tst.sibling)
		if err != nil {
			if fmt.Sprintf("%s", err) != tst.want {
				t.Errorf("Add(%d, %s) returns error: %q, should be %s",
					tst.n, tst.sibling, err, tst.want)
			}
		} else if s := root.String(); s != tst.want {
			t.Errorf("Add(%d, %s) returns: %q, should be %q",
				tst.n, tst.sibling, s, tst.want)
		}
	}

	root = &Branch{ID: "root"}

	tests2 := []struct {
		n       int
		sibling []interface{}
		want    string
	}{
		{-1, []interface{}{"s2", "s3"}, "root{s2 s3}"},
		{0, []interface{}{"s0", "s1"}, "root{s0 s1 s2 s3}"},
		{-1, []interface{}{"s6", "s7"}, "root{s0 s1 s2 s3 s6 s7}"},
		{4, []interface{}{"s4", "s5"}, "root{s0 s1 s2 s3 s4 s5 s6 s7}"},
	}

	for _, tst := range tests2 {
		err := root.Add(tst.n, tst.sibling...)
		if err != nil {
			if fmt.Sprintf("%s", err) != tst.want {
				t.Errorf("Add(%d, %s) returns error: %q, should be %s",
					tst.n, tst.sibling, err, tst.want)
			}
		} else if s := root.String(); s != tst.want {
			t.Errorf("Add(%d, %s) returns: %q, should be %q",
				tst.n, tst.sibling, s, tst.want)
		}
	}

}

func TestAddBranch(t *testing.T) {
	tests := []struct {
		n    int
		id   string
		want string
	}{
		{0, "s2", fmt.Sprintf("%s", ErrNoSiblings)},
		{-1, "b2", "root{b2{s}}"},
		{0, "b0", "root{b0{s} b2{s}}"},
		{1, "b1", "root{b0{s} b1{s} b2{s}}"},
		{-1, "b3", "root{b0{s} b1{s} b2{s} b3{s}}"},
		{4, "b3", "sibling index (4) out of range (0-3)"},
	}

	root := &Branch{ID: "root"}

	for _, tst := range tests {
		br, err := root.AddBranch(tst.n, tst.id)
		if err != nil {
			if fmt.Sprintf("%s", err) != tst.want {
				t.Errorf("AddBranch(%d, %s) returns error: %q, should be: %q",
					tst.n, tst.id, err, tst.want)
			}
		} else {
			br.Add(-1, "s")
			if s := root.String(); s != tst.want {
				t.Errorf("AddBranch(%d, %s) returns: %q, should be %q",
					tst.n, tst.id, s, tst.want)
			}
		}
	}
}

func TestIdParent(t *testing.T) {
	brnchs := []struct {
		n  int
		id string
	}{
		{-1, "b0"},
		{-1, "b1"},
		{-1, "b2"},
		{-1, "b3"},
	}

	l := len(brnchs)
	root := &Branch{ID: "root"}

	b := root
	var err error
	for _, tst := range brnchs {
		b, err = b.AddBranch(tst.n, tst.id)
		if err != nil {
			t.Fatalf("AddBranch(%d, %s) returns error: %q, should be nil",
				tst.n, tst.id, err)
		}
	}

	for i := 0; i < l; i++ {
		p, err := b.ParentById(brnchs[i].id)
		if err != nil {
			t.Fatalf("ParentById(%q) generates an error: %q", brnchs[i].id, err)
		}
		if p.ID != brnchs[i].id {
			t.Fatalf("ParentById(%q).ID is %q, should be %q", brnchs[i].id, p.ID,
				brnchs[i].id)
		}
	}
}

func TestLen(t *testing.T) {
	brnchs := []struct {
		n  int
		id string
	}{
		{-1, "b2"},
		{0, "b0"},
		{1, "b1"},
		{-1, "b3"},
	}

	root := &Branch{ID: "root"}

	for _, tst := range brnchs {
		_, err := root.AddBranch(tst.n, tst.id)
		if err != nil {
			t.Fatalf("AddBranch(%d, %s) returns error: %q, should be nil",
				tst.n, tst.id, err)
		}
	}

	if root.Len() != len(brnchs) {
		t.Errorf("Len() returns: %d, should be: %d", root.Len(), len(brnchs))
	}
}

func TestParent(t *testing.T) {
	brnchs := []struct {
		n  int
		id string
	}{
		{-1, "b0"},
		{-1, "b1"},
		{-1, "b2"},
		{-1, "b3"},
	}

	l := len(brnchs)
	root := &Branch{ID: "root"}

	b := root
	var err error
	for _, tst := range brnchs {
		b, err = b.AddBranch(tst.n, tst.id)
		if err != nil {
			t.Fatalf("AddBranch(%d, %s) returns error: %q, should be nil",
				tst.n, tst.id, err)
		}
	}

	for i := 0; i < l; i++ {
		p, err := b.Parent(i)
		if err != nil {
			t.Errorf("AddBranch(%d) returns error: %q, should be nil", i, err)
		}
		if p.ID != brnchs[l-1-i].id {
			t.Errorf("Parent(%d).ID returns: %q, should be: %q",
				i, p.ID, brnchs[len(brnchs)-i].id)
		}
	}
}

func TestRemove(t *testing.T) {
	brnchs := []struct {
		n  int
		id string
	}{
		{-1, "b2"},
		{0, "b0"},
		{1, "b1"},
		{-1, "b3"},
	}

	root := &Branch{ID: "root"}
	for _, b := range brnchs {
		br, err := root.AddBranch(b.n, b.id)
		if err != nil {
			t.Fatalf("AddBranch(%d, %s) returns error: %q, should be nil",
				b.n, b.id, err)
		}
		br.Add(-1, "s")
	}

	tests := []struct {
		n    int
		want string
	}{
		{1, "root{b0{s} b2{s} b3{s}}"},
		{0, "root{b2{s} b3{s}}"},
		{-1, "root{b2{s}}"},
	}

	for i, tst := range tests {
		_, err := root.Remove(tst.n)
		if err != nil {
			t.Fatalf("Remove(%d) returns error: %q, should be nil", tst.n, err)
		}
		if s := root.String(); s != tst.want {
			t.Errorf("Remove(%d) returns: %q != %q", i, s, tst.want)
		}
	}
}

func TestSiblings(t *testing.T) {
	sibs := []struct {
		n       int
		sibling string
	}{
		{-1, "s0"},
		{-1, "s1"},
		{-1, "s2"},
		{-1, "s3"},
		{-1, "s4"},
	}

	root := &Branch{ID: "root"}
	for _, s := range sibs {
		if err := root.Add(s.n, s.sibling); err != nil {
			t.Fatalf("Add(%d, %s) returns error: %q, should be nil",
				s.n, s.sibling, err)
		}
	}

	siblings := root.Siblings()
	for i, s := range siblings {
		if s != sibs[i].sibling {
			t.Errorf("Siblings()[%d] is: %q, should be %q", i, s, sibs[i].sibling)
		}
	}
}

func TestSiblingsN(t *testing.T) {
	sibs := []struct {
		n       int
		sibling string
	}{
		{-1, "s1"},
		{-1, "s2"},
		{-1, "s3"},
		{-1, "s4"},
		{-1, "s5"},
	}

	root := &Branch{ID: "root"}
	for _, s := range sibs {
		if err := root.Add(s.n, s.sibling); err != nil {
			t.Fatalf("Add(%d, %s) returns error: %q, should be nil",
				s.n, s.sibling, err)
		}
	}

	for i := 0; i < 5; i++ {
		s, err := root.SiblingN(i)
		if err != nil {
			t.Fatalf("SiblingN(%d) returns error: %q, should be nil", i, err)
		}
		if s != sibs[i].sibling {
			t.Errorf("SiblingN(%d) returns: %q, should be %q", i, s, sibs[i].sibling)
		}
	}
}

func TestIndex(t *testing.T) {
	sibs := []struct {
		n       int
		sibling string
	}{
		{-1, "s0"},
		{-1, "s1"},
		{-1, "s2"},
	}

	root := &Branch{ID: "root"}
	for i, s := range sibs {
		if i == 1 {
			if _, err := root.AddBranch(-1, s.sibling); err != nil {
				t.Fatalf("AddBranch(%d, %s) returns error: %q, should be nil",
					s.n, s.sibling, err)
			}
		} else if err := root.Add(s.n, s.sibling); err != nil {
			t.Fatalf("Add(%d, %s) returns error: %q, should be nil",
				s.n, s.sibling, err)
		}
	}

	for i, s := range root.siblings {
		if j, err := root.Index(s); j != i {
			if err != nil {
				t.Errorf("Index(..) returns error %q, should be nil", err)
			}
			t.Errorf("Index(..) returns %d, should be %d", j, i)
		}
	}

}
