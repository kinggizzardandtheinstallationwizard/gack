package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"flag"
	"fmt"
	"os"
)

var docify = flag.Bool("p", false, "Make a complete web(p)age.")
var writef = flag.Bool("f", false, "Write results to a new (f)ile.")

var pf_flags = flag.Bool("pf", false, "")
var fp_flags = flag.Bool("fp", false, "")

func main() {
	var in string
	var out string

	// Uuuugly argparse stuff
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage:", os.Args[0], "[-pf] file.md") // print usage
		os.Exit(0)
	}

	switch len(os.Args) {
	case 2:
		in = os.Args[1]
	case 3:
		in = os.Args[2]
	}

	contents := string(fparse(in))
	if contents == "" {
		os.Exit(0)
	}

	flag.Parse()
	both := *pf_flags || *fp_flags
	if both {
		*docify = true
		*writef = true
	}

	// Do HTML header/footer stuff if the user wants it
	if *docify {
		out  = "<!DOCTYPE html>\n"
		out += "<html lang=\"en\">\n"

		out += "<head>\n"
		out += "<title>"
		out += in
		out += "</title>\n"
		out += "</head>\n"

		out += "<body>\n"
		out += contents
		out += "</body>\n"

		out += "</html>\n"
	} else {
		out = contents
	}

	if *writef {
		// file crap
		fwrite(in, out)
	} else {
		fmt.Print(out)
	}
}

func fparse(fname string) string {
	dat, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println("!", err)
		fmt.Println("Oh, the humanity!")
		return ""
	}

	// Parse file
	exts := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(exts)
	parsed := p.Parse(dat)

	htflags := html.CommonFlags
	htopts := html.RendererOptions{Flags: htflags}
	r := html.NewRenderer(htopts)

	rendered := markdown.Render(parsed, r)

	// Let main() take care of actually outputting
	return string(rendered)
}

func fwrite(raw_name string, contents string) {
	// Assume/change file extension ".md" --> ".html" (naive fix)
	fname := raw_name[:len(raw_name)-2] + "html"

	// Open
	fout, err := os.Create(fname)
	if err != nil {
		fmt.Println("!", err)
		fmt.Println("Oh, the humanity!")
		return
	}

	// Write
	_, err = fout.WriteString(contents)
	if err != nil {
		fmt.Println("!", err)
		fmt.Println("Oh, the humanity!")
		return
	}

	fout.Close()
}
