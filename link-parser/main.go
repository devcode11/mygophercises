package main

import (
	"fmt"
	"main/parser"
	"os"
	"strings"
	"flag"
)

func main() {
	str := flag.String("str", "", "HTML `string` to parse")
	path := flag.String("path", "", "HTML `file` to parse")
	flag.Parse()

	if (*str == "" && *path == "") {
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
		return
	}

	var links []parser.Link
	var err error

	if (*str != "") {
		links, err = parser.Links(strings.NewReader(*str))
	} else {
		f, err := os.Open(*path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer f.Close()
		links, err = parser.Links(f)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
	if len(links) == 0 {
		fmt.Fprintln(os.Stderr, "No links found")
		return
	}
	for i, link := range links {
		fmt.Printf("%d) %v => %v\n", i+1, link.Text, link.Href)
	}
}
