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
	// arch0   = flag.Arg(0) //"md2html"
	kVersion = "0.1"
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

func Configure() *Config {
	cfg := new(Config)

	// parse de argumenten
	version := flag.Bool("version", false, "Show version number and exit")
	if *version {
		fmt.Printf("Version %s\n", kVersion)
		os.Exit(0)
	}

	input := flag.String("in", "stdin", "path to input file")
	output := flag.String("out", "stdout", "path to output file")
	flag.StringVar(&cfg.title, "title", "", "title for HTML document")
	flag.StringVar(&cfg.style, "style", "", "style sheet for HTML document")

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

// Header returns a branch holding HTML head data.
func (cfg *Config) Header() *branch.Branch {
	head := branch.NewBranch("head")

	if len(cfg.title) > 0 {
		title, _ := head.AddBranch(-1, "title")
		title.Add(-1, cfg.title)
	}

	meta, _ := head.AddBranch(-1, "meta")
	meta.Info = "charset=\"utf-8\""
	meta.Add(-1, "")

	meta, _ = head.AddBranch(-1, "meta")
	meta.Info = "name=\"generator\" content=\"md2html\""
	meta.Add(-1, "")

	meta, _ = head.AddBranch(-1, "meta")
	meta.Info = "http-equiv=\"Content-Style-Type\" content=\"text/css\""
	meta.Add(-1, "")

	if len(cfg.style) > 0 {
		style, _ := head.AddBranch(-1, "style")
		style.Info = fmt.Sprintf("rel=\"stylesheet\" href=\"%s\" type=\"text/css\"",
			cfg.style)
		style.Add(-1, "")
	}

	return head
}

func main() {

	cfg := Configure()

	html := NewHTMLTree("html")
	html.root.Add(-1, cfg.Header())

	body, err := BuildHTMLTree(cfg.fIn)
	if err != nil {
		fmt.Printf("building HTML tree: %s\n", err)
		os.Exit(1)
	}
	html.root.Add(-1, body)

	fmt.Fprintf(cfg.fOut, "%s", HTMLCode(html.root))
}
