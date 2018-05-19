//
// md2html.go
//
//  Questions and comments to:
//       <mailto:frank@foef.nl>
//
// stack implements a Last In First Out buffer of values of type
// interface{}.
//
// © 2018 Frank Storbeck

package main

import (
	"fmt"
	"strings"

	"source.storbeck.nl/md2html/branch"
)

// StrongEmDel translates mark down strong, emphasis and deleted definitions
// to their html equivalents
func StrongEmDel(s string) string {
	tgs := []struct {
		tg  string
		sep []string
	}{
		{"strong", []string{"**", "__"}},
		{"em", []string{"_", "*"}},
		{"del", []string{"~~", "~~"}},
	}

	seps := make([]byte, len(tgs))
	for i, tg := range tgs {
		seps[i] = tg.sep[0][0]
	}
	rslt := UniCode(s, seps)

	for _, t := range tgs {
		for _, sp := range t.sep {
			var subs []string
			switch {
			// special case: "^\*[^\*]+.*$" will become an unordered list, no italics!
			case sp == "*" && len(rslt) > 0 && rslt[:1] == sp:
				subs = strings.Split(rslt[1:], sp)
				rslt = sp + subs[0]
			default:
				subs = strings.Split(rslt, sp)
				rslt = subs[0]
			}
			n := len(subs)
			for i := 1; i < n; i = i + 2 {
				if i+1 < n {
					rslt = rslt + "<" + t.tg + ">" + subs[i] + "</" + t.tg + ">" + subs[i+1]
				} else {
					rslt = rslt + sp + subs[i]
				}
			}
		}
	}

	return UnEscape(rslt, seps)
}

func main() {
	root := &branch.Branch{ID: "html"}
	fmt.Printf("%s\n", root.ID)
}
