package main

import (
	"testing"
)

func TestStyling(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "", want: ""},
		{s: "**bb**", want: "<strong>bb</strong>"},
		{s: "xx**bb**", want: "xx<strong>bb</strong>"},
		{s: "**bb**yy", want: "<strong>bb</strong>yy"},
		{s: "xx**bb**yy", want: "xx<strong>bb</strong>yy"},
		{s: "**pp", want: "**pp"},
		{s: "xx**pp", want: "xx<em></em>pp"},
		{s: "**ppyy", want: "**ppyy"},
		{s: "xx__ppyy", want: "xx<em></em>ppyy"},
		{s: "**bb**xx**bb**__bb__**bb**yy", want: "<strong>bb</strong>xx<strong>bb</strong><strong>bb</strong><strong>bb</strong>yy"},
		{s: "**bb***bb****", want: "<strong>bb</strong>*bb<strong></strong>"},
		{s: "~~bb~~", want: "<del>bb</del>"},
		{s: "xx~~bb~~", want: "xx<del>bb</del>"},
		{s: "~~bb~~yy", want: "<del>bb</del>yy"},
	}
	for _, tst := range tests {
		got := StrongEmDel(tst.s)
		if got != tst.want {
			t.Errorf("StrongEmDel(%q) generates: %q, should be: %q", tst.s, got, tst.want)
		}
	}
}

func TestUniCode(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "a \\* b", want: "a U+002A b"},
		{s: "\\* b", want: "U+002A b"},
		{s: "a \\*", want: "a U+002A"},
		{s: "a \\_ b", want: "a \\_ b"},
	}
	for _, tst := range tests {
		got := UniCode(tst.s, []byte{'*'})
		if got != tst.want {
			t.Errorf("UniCode(%q, []byte{'*'}) generates: %q, should be: %q", tst.s, got, tst.want)
		}
	}
}

func TestUnEscape(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "a U+002A b", want: "a * b"},
		{s: "U+002A b", want: "* b"},
		{s: "a U+002A", want: "a *"},
		{s: "a U+002B b", want: "a U+002B b"},
	}
	for _, tst := range tests {
		got := UnEscape(tst.s, []byte{'*'})
		if got != tst.want {
			t.Errorf("UnEscape(%q, []byte{'*'}) generates: %q, should be: %q", tst.s, got, tst.want)
		}
	}
}

func TestOlnlyRunes(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{s: "", want: false},
		{s: "==", want: false},
		{s: "===", want: true},
		{s: "=======", want: true},
		{s: "=======x", want: false},
		{s: "====x==", want: false},
	}

	for _, tst := range tests {
		got := OnlyRunes(tst.s, '=')
		if got != tst.want {
			t.Errorf("'OnlyRunes(%q, '=')' generates: %t, should be: %t",
				tst.s, got, tst.want)
		}
	}
}

func TestCountLeading(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{s: "", want: 0},
		{s: "x", want: 0},
		{s: "#x", want: 1},
		{s: "######x", want: 6},
		{s: "#######x", want: 0},
		{s: " #####", want: 0},
	}

	for _, tst := range tests {
		n := CountLeading(tst.s, '#', 6)
		if n != tst.want {
			t.Errorf("'CountLeading(%q)' generates: %d, should be: %d",
				tst.s, n, tst.want)
		}
	}
}

func TestBuild(t *testing.T) {
	tests := []struct {
		s    []string
		want string
	}{
		// Headers
		{s: []string{"hdr1", "==="}, want: "r{h1{hdr1} p{}}"},
		{s: []string{"hdr2", "---"}, want: "r{h2{hdr2} p{}}"},
	}
	for _, tst := range tests {
		ht := NewHTMLTree("r")
		ht.br, _ = ht.br.AddBranch(-1, "p")

		for _, s := range tst.s {
			err := ht.Build(s)
			if err != nil {
				t.Fatalf("Build(%q) returns error: %s, should be nil", s, err)
			}
		}

		got := ht.root.String()
		if got != tst.want {
			t.Errorf("Build(%q)... generates %q, should be %q",
				tst.s[0], got, tst.want)
		}
	}
}
