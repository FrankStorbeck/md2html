package main

// //
// func SelectInlineParts(s /*intf []interface{}*/, sep string) []interface{} {
// 	r := []interface{}{}
// 	// for _, in := range intf {
// 	// 	if s, ok := in.(string); ok {
// 	subs := strings.Split(s, sep)
// 	n := len(subs)
//
// 	switch {
// 	case n == 0:
// 		return r
// 	case n == 1:
// 		return append(r, subs[0])
// 	case n%2 == 0:
// 		n--
// 		subs[n-1] = subs[n-1] + sep + subs[n]
// 		subs = subs[:n]
// 	}
//
// 	for i := 0; i < n; i = i + 2 {
// 		if len(subs[i]) > 0 {
// 			r = append(r, subs[i])
// 		}
// 		if i+1 < n {
// 			if len(subs[i+1]) > 0 {
// 				if i+1 < n-1 {
// 					strong := branch.NewBranch("strong")
// 					r = append(r, strong)
// 					strong.Add(-1, subs[i+1])
// 				} else {
// 					r = append(r, subs[i+1])
// 				}
// 			} else {
// 				r = append(r, subs[i+1])
// 			}
// 			// 		}
// 			// 	}
// 			// } else {
// 			// 	r = append(r, in)
// 			// }
// 		}
// 	}
//
// 	return r
// }
//
// // Images2 splits a string in a series of siblings consisting of strings and
// // <img> branches.
// func Images2(s string) []interface{} {
// 	l := len(s)
// 	if i := strings.Index(s, "!["); i >= 0 && l > i+5 {
// 		if j := strings.Index(s[i:], "]"); j > 0 && l > (i+j+2) {
// 			if s[i+j+1] == '(' {
// 				if k := strings.Index(s[i+j+1:], ")"); k > 0 {
// 					r := []interface{}{}
// 					if i > 0 {
// 						r = append(r, s[:i])
// 					}
// 					im := branch.NewBranch("img")
// 					r = append(r, im)
// 					im.Info = "src=\"" + s[i+j+2:i+j+k+1] + "\" alt=\"" +
// 						s[i+2:i+j] + "\""
// 					if i+j+k+2 < l {
// 						r = append(r, Images2(s[i+j+k+2:])...)
// 					}
// 					return r
// 				}
// 			}
// 		}
// 	}
// 	return []interface{}{s}
// }
//
// // Links2 splits a string in a series of siblings consisting of strings and
// // <a> branches.
// func Links2(s string) []interface{} {
// 	l := len(s)
// 	if i := strings.Index(s, "["); i >= 0 && l > i+4 {
// 		if j := strings.Index(s[i:], "]"); j > 0 && l > (i+j+1) {
// 			if s[i+j+1] == '(' {
// 				if k := strings.Index(s[i+j+1:], ")"); k > 0 {
// 					r := []interface{}{}
// 					if i > 0 {
// 						r = append(r, s[:i])
// 					}
// 					a := branch.NewBranch("a")
// 					r = append(r, a)
// 					a.Info = "href=\"" + s[i+j+2:i+j+k+1] + "\""
// 					a.Add(-1, s[i+1:i+j])
// 					if i+j+k+2 < l {
// 						r = append(r, Links2(s[i+j+k+2:])...)
// 					}
// 					return r
// 				}
// 			}
// 		}
// 	}
// 	return []interface{}{s}
// }
//
// /// // InlineCodes2 splits a string in a series of siblings consisting of strings
// import (
// 	"html"
// 	"strings"
//
// 	"source.storbeck.nl/md2html/branch"
// )
//
// // and <code> branches.
// func InlineCodes2(s string) []interface{} {
// 	l := len(s)
// 	if i := strings.Index(s, "`"); i >= 0 && l > i+2 {
// 		if j := strings.Index(s[i+1:], "`"); j > 0 {
// 			r := []interface{}{}
// 			if i > 0 {
// 				r = append(r, s[:i])
// 			}
// 			code := branch.NewBranch("code")
// 			r = append(r, code)
// 			code.Add(-1, html.EscapeString(s[i+1:i+j+1]))
// 			if i+j+2 < l {
// 				r = append(r, InlineCodes2(s[i+j+2:])...)
// 			}
// 			return r
// 		}
// 	}
// 	return []interface{}{s}
// }
//
// // Strong2
// func Strong2(s string) []interface{} {
// 	return SelectInlineParts(s, "**")
// }
//
// func Em2(s string) []interface{} {
// 	return SelectInlineParts(s, "_")
// }
//
// func Del2(s string) []interface{} {
// 	return SelectInlineParts(s, "~~")
// }
//
// func TestImages2(t *testing.T) {
// 	tests := []struct {
// 		s    string
// 		want string
// 	}{
// 		{s: "aa [img](link) bb",
// 			want: "r{p{aa [img](link) bb}}"},
// 		{s: "aa ![img](link) bb",
// 			want: "r{p{aa  img:src=\"link\" alt=\"img\"{}  bb}}"},
// 		{s: "![i1](l1)![i2](l2)",
// 			want: "r{p{img:src=\"l1\" alt=\"i1\"{} img:src=\"l2\" alt=\"i2\"{}}}"},
// 	}
//
// 	for _, tst := range tests {
// 		ht := NewHTMLTree("r")
// 		ht.br, _ = ht.br.AddBranch(-1, "p")
// 		for _, im := range traverse([]interface{}{tst.s}, Images2) {
// 			ht.br.Add(-1, im)
// 		}
// 		got := ht.root.String()
// 		if got != tst.want {
// 			t.Errorf("Images2(%q) generates: %q, should be %q", tst.s, got, tst.want)
// 		}
// 	}
// }
//
// func TestLinks2(t *testing.T) {
// 	tests := []struct {
// 		s    string
// 		want string
// 	}{
// 		{s: "aa[txt](link)bb", want: "r{p{aa a:href=\"link\"{txt} bb}}"},
// 		{s: "[t1](l1)[t2](l2)", want: "r{p{a:href=\"l1\"{t1} a:href=\"l2\"{t2}}}"},
// 	}
//
// 	for _, tst := range tests {
// 		ht := NewHTMLTree("r")
// 		ht.br, _ = ht.br.AddBranch(-1, "p")
// 		for _, im := range traverse([]interface{}{tst.s}, Links2) {
// 			ht.br.Add(-1, im)
// 		}
// 		got := ht.root.String()
// 		if got != tst.want {
// 			t.Errorf("Link2(%q) generates:\n %q\n, should be:\n %q\n", tst.s, got, tst.want)
// 		}
// 	}
// }
//
// func TestInlineCode2(t *testing.T) {
// 	tests := []struct {
// 		s    string
// 		want string
// 	}{
// 		{s: "aa`code`bb", want: "r{p{aa code{code} bb}}"},
// 	}
// 	for _, tst := range tests {
// 		ht := NewHTMLTree("r")
// 		ht.br, _ = ht.br.AddBranch(-1, "p")
// 		for _, im := range traverse([]interface{}{tst.s}, InlineCodes2) {
// 			ht.br.Add(-1, im)
// 		}
// 		got := ht.root.String()
// 		if got != tst.want {
// 			t.Errorf("InlineCodes2(%q) generates:\n %q\n, should be:\n %q\n", tst.s, got, tst.want)
// 		}
// 	}
// }
//
// func TestStrong2(t *testing.T) {
// 	tests := []struct {
// 		s    string
// 		want string
// 	}{
// 		{s: "", want: "r{p{}}"},
// 		{s: "aa", want: "r{p{aa}}"},
// 		{s: "**aa", want: "r{p{**aa}}"},
// 		{s: "**aa**", want: "r{p{strong{aa}}}"},
// 		{s: "aa**", want: "r{p{aa**}}"},
// 		{s: "aa**bb", want: "r{p{aa**bb}}"},
// 		{s: "aa**bb**cc", want: "r{p{aa strong{bb} cc}}"},
// 		{s: "aa**bb**cc**", want: "r{p{aa strong{bb} cc**}}"},
// 		{s: "aa**bb**cc**dd", want: "r{p{aa strong{bb} cc**dd}}"},
// 		{s: "aa**bb**cc**dd**", want: "r{p{aa strong{bb} cc strong{dd}}}"},
// 		{s: "aa**bb**cc**dd**ee", want: "r{p{aa strong{bb} cc strong{dd} ee}}"},
// 	}
//
// 	for _, tst := range tests {
// 		ht := NewHTMLTree("r")
// 		ht.br, _ = ht.br.AddBranch(-1, "p")
// 		for _, im := range traverse([]interface{}{tst.s}, Strong2) {
// 			ht.br.Add(-1, im)
// 		}
// 		got := ht.root.String()
// 		if got != tst.want {
// 			t.Errorf("Strong2(%q) generates:\n %q\n, should be:\n %q\n", tst.s, got, tst.want)
// 		}
// 	}
// }
//
