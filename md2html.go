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
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"source.storbeck.nl/md2html/branch"
)

const (
	cBlockQuote = "blockquote"
	cBody       = "body"
	cCode       = "code"
	cHead       = "head"
	cHTML       = "html"
	cH1         = "h1"
	cH2         = "h2"
	cH3         = "h3"
	cH4         = "h4"
	cH5         = "h5"
	cH6         = "h6"
	cLi         = "li"
	cLink       = "link"
	cMeta       = "meta"
	cOl         = "ol"
	cP          = "p"
	cPre        = "pre"
	cQ          = "q"
	cCrLf       = "\r\n"
	cScript     = "script"
	cStyle      = "style"
	cTable      = "table"
	cTd         = "td"
	cTh         = "th"
	cTitle      = "title"
	cTr         = "tr"
	cUl         = "ul"
	cVersion    = "0.1"
)

// Config holds all configuration data
type Config struct {
	fIn   *os.File
	fOut  *os.File
	style string
	title string
}

// BuildHTMLTree returns a pointer to a branch struct with all HTML elements
// from a named mark down file. In case of an error 'nil' and the error will be
// returned.
func BuildHTMLTree(f *os.File) (*branch.Branch, error) {
	buf := bufio.NewReader(f)

	st := NewHTMLTree(cBody)
	st.br, _ = st.root.AddBranch(-1, cP)

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

// Configure sets the configuration for 'main' based on its flags.
func Configure() *Config {
	cfg := new(Config)

	// parse de argumenten
	version := flag.Bool("version", false, "Show version number and exit")
	if *version {
		fmt.Printf("Version %s\n", cVersion)
		os.Exit(0)
	}

	input := flag.String("in", "stdin", "path to input file")
	output := flag.String("out", "stdout", "path to output file")
	flag.StringVar(&cfg.title, cTitle, "", "title for HTML document")
	flag.StringVar(&cfg.style, cStyle, "", "style sheet for HTML document")

	flag.Parse()

	var err error
	cfg.fIn = os.Stdin
	if *input != "stdin" {
		cfg.fIn, err = os.Open(*input)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	cfg.fOut = os.Stdout
	if *output != "stdout" {
		cfg.fOut, err = os.Create(*output)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	return cfg
}

// HTMLCode returns a string holding the html code.
func HTMLCode(br *branch.Branch, lvl int) string {
	sbl := br.Siblings()
	if len(sbl) <= 0 {
		return ""
	}

	indnt := ""
	if lvl > 0 {
		indnt = strings.Repeat(" ", lvl)
	}

	s := ""
	switch br.ID {
	case cTable:
		s = cCrLf
	}
	s = s + indnt + "<" + br.ID
	if len(br.Info) > 0 {
		s = s + " " + strings.TrimSpace(br.Info)
	}

	switch br.ID {
	case cLink, cMeta:
		s = s + "/>\n"
	default:
		s = s + ">"
		switch br.ID {
		case cBlockQuote, cBody, cCode, cHead, cHTML, cOl, cPre, cTable, cTr, cUl:
			s = s + cCrLf
		}

		spc := ""
		for _, c := range sbl {
			switch k := c.(type) {
			case *branch.Branch:
				l := lvl
				switch {
				case br.ID == cLi:
					l = -1
				case l >= 0:
					l++
				default:
				}
				s = s + HTMLCode(k, l)
				spc = ""
			case string:
				s = s + spc + k
				if s[len(s)-1] == '\n' {
					spc = ""
				} else {
					spc = " "
				}
			default:
				// s = s + string(k)
			}
		}

		nl := ""
		switch br.ID {
		case cBlockQuote:
			nl = cCrLf + indnt
		case cBody, cCode, cHead, cOl, cPre, cTable, cTr, cUl:
			nl = indnt
		}

		s = s + nl + "</" + br.ID + ">"

		if lvl >= 0 {
			switch br.ID {
			case cTable:
				s = s + cCrLf + strings.Repeat(" ", lvl-1)
			case cBody, cBlockQuote, cCode, cHead, cHTML, cH1, cH2, cH3, cH4, cH5,
				cH6, cLi, cLink, cOl, cP, cPre, cQ, cTitle, cScript, cStyle, cTd, cTh, cTr, cUl:
				s = s + cCrLf
			}
		}
	}
	return s
}

// Header returns a branch holding HTML head data.
func (cfg *Config) Header() *branch.Branch {
	head := branch.NewBranch(cHead)

	if len(cfg.title) > 0 {
		title, _ := head.AddBranch(-1, cTitle)
		title.Add(-1, cfg.title)
	}

	meta, _ := head.AddBranch(-1, cMeta)
	meta.Info = "charset=\"utf-8\""
	meta.Add(-1, "")

	meta, _ = head.AddBranch(-1, cMeta)
	meta.Info = "name=\"generator\" content=\"md2html\""
	meta.Add(-1, "")

	meta, _ = head.AddBranch(-1, cMeta)
	meta.Info = "http-equiv=\"Content-Style-Type\" content=\"text/css\""
	meta.Add(-1, "")

	if len(cfg.style) > 0 {
		style, _ := head.AddBranch(-1, cLink)
		style.Info = fmt.Sprintf("rel=\"stylesheet\" href=\"%s\" type=\"text/css\"",
			cfg.style)
		style.Add(-1, "")
	}

	return head
}

func main() {

	cfg := Configure()

	html := NewHTMLTree(cHTML)
	html.root.Add(-1, cfg.Header())

	body, err := BuildHTMLTree(cfg.fIn)
	if err != nil {
		fmt.Printf("building HTML tree: %s\n", err)
		os.Exit(1)
	}
	html.root.Add(-1, body)

	fmt.Fprintf(cfg.fOut, "%s", HTMLCode(html.root, 0))
}
