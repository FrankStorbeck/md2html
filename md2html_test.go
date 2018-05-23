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
		{s: []string{"aa", "# hdr1", "bb"}, want: "r{p{aa} h1{hdr1} p{bb}}"},
		{s: []string{"aa", "### hdr3", "bb"}, want: "r{p{aa} h3{hdr3} p{bb}}"},
		{s: []string{"###### hdr6"}, want: "r{h6{hdr6} p{}}"},
		{s: []string{"####### hdr7"}, want: "r{p{####### hdr7}}"},

		// Quoting
		{s: []string{"> quote"},
			want: "r{blockquote{quote }}"},
		{s: []string{"aa", "> quote1", "> quote2", "bb"},
			want: "r{p{aa} blockquote{quote1  quote2 } p{bb}}"},
		{s: []string{"aa", "> quote1", "", "> quote2", "bb"},
			want: "r{p{aa} blockquote{quote1  <br> quote2 } p{bb}}"},
		{s: []string{"aa`cc`bb"}, want: "r{p{aa<code>cc</code>bb}}"},
		{s: []string{"aa", "```", "a1", "a2", "```", "bb"}, want: "r{p{aa} pre{code{a1 a2}} p{bb}}"},

		// Lists
		{s: []string{"aa", "* 1", "* 2", "  + 2.1", "  + 2.2", "    - 2.2.1", "  + 2.3", "* 3", "cc"},
			want: "r{p{aa} ul{li{1} li{2} ul{li{2.1} li{2.2} ul{li{2.2.1}} li{2.3}} li{3}} p{cc}}"},
		{s: []string{"aa", "  * 1", "   l1", "  * 2", "   l2", "c"},
			want: "r{p{aa} ul{li{1 l1} li{2 l2}} p{c}}"},
		{s: []string{"* 1", "* 2", "a", "* p", "* q"},
			want: "r{ul{li{1} li{2}} p{a} ul{li{p} li{q}}}"},
		{s: []string{"a", "* 1", "* 2", "  + a", "* 3", "", "  b", " c", "* 4", "d"},
			want: "r{p{a} ul{li{1} li{2} ul{li{a}} li{p{3 b c}} li{4}} p{d}}"},
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

func TestImages(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "aa ![img](link) bb", want: "aa <img src=\"link\" alt=\"img\"> bb"},
		{s: "![im1](lin1)![im2](lin2)", want: "<img src=\"lin1\" alt=\"im1\"><img src=\"lin2\" alt=\"im2\">"},
	}

	for _, tst := range tests {
		got := Images(tst.s)
		if got != tst.want {
			t.Errorf("Images(%q) generates: %q, should be %q", tst.s, got, tst.want)
		}
	}
}

func TestInlineCodes(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "aa `code` bb", want: "aa <code>code</code> bb"},
	}

	for _, tst := range tests {
		got := InlineCodes(tst.s)
		if got != tst.want {
			t.Errorf("InlineCodes(%q) generates: %q, should be %q", tst.s, got, tst.want)
		}
	}
}

func TestLinks(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "aa [txt](link) bb", want: "aa <a href=\"link\">txt</a> bb"},
		{s: "[t1](l1)[t2](l2)", want: "<a href=\"l1\">t1</a><a href=\"l2\">t2</a>"},
	}

	for _, tst := range tests {
		got := Links(tst.s)
		if got != tst.want {
			t.Errorf("Links(%q) generates: %q, should be %q", tst.s, got, tst.want)
		}
	}
}
