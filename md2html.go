//
// md2html.go
//
//  Questions and comments to:
//       <mailto:frank@foef.nl>
//
// md2html translates the contents of a markdown file into a file holding the
// HTML equivalent.
//
// © 2018 Frank Storbeck

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"source.storbeck.nl/md2html/branch"
)

// BuildHTMLTree returns a pointer to a branch struct with all HTML elements
// from a named mark down file. In case of an error 'nil' and the error will be
// returned.
func BuildHTMLTree(path string) (*branch.Branch, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	st := NewHTMLTree("body")
	st.br, _ = st.root.AddBranch(-1, "p")

	for {
		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return st.root, err
		}
		st.Build(line)
		if err == io.EOF {
			break
		}
	}

	return st.root, nil
}

// HTMLCode returns a string holding the html code.
func HTMLCode(br *branch.Branch) string {
	sbl := br.Siblings()
	if len(sbl) <= 0 {
		return ""
	}

	s := "<" + br.ID
	if len(br.Info) > 0 {
		s = s + " " + strings.TrimSpace(br.Info)
	}
	s = s + ">"

	switch br.ID {
	case "body", "ol", "table", "tr", "ul":
		s = s + "\n"
	}

	spc := ""
	for _, c := range sbl {
		switch k := c.(type) {
		case *branch.Branch:
			s = s + HTMLCode(k)
			spc = ""
		case string:
			s = s + spc + k
		default:
			// s = s + string(k)
		}
	}

	nl := ""
	switch br.ID {
	case "blockquote":
		nl = "\n"
	}

	switch br.ID {
	case "link", "meta":

	default:
		s = s + nl + "</" + br.ID + ">"
	}

	switch br.ID {
	case "blockquote", "h1", "h2", "h3", "h4", "h5", "h6", "li", "link", "meta",
		"ol", "p", "pre", "q", "script", "table", "td", "th", "tr", "ul":
		s = s + "\n"
	}
	return s
}

func main() {
	html := NewHTMLTree("html")
	head, _ := html.br.AddBranch(-1, "head")
	title, _ := head.AddBranch(-1, "title")
	title.Add(-1, "some title")
	meta, _ := head.AddBranch(-1, "meta")
	meta.Info = "charset=\"utf-8\""
	meta.Add(-1, "")

	body, err := BuildHTMLTree(filepath.Join(".", "syntaxMD.md"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	html.root.Add(-1, body)

	fmt.Printf("%s", HTMLCode(html.root))
}
