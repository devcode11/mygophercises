package main

import (
	"flag"
	"fmt"
	"os"

	"main/sitemap"
)

func main() {
	startUrl := flag.String("url", "https://www.calhoun.io", "`url` to start building sitemap")
	flag.Parse()
	fmt.Println("URL:", *startUrl)
	if *startUrl == "" {
		flag.PrintDefaults()
		return
	}

	siteMap, err := sitemap.Build(*startUrl)

	if err != nil {
		fmt.Println(err)
		return
	}

	// for i, link := range siteMap {
	// 	fmt.Printf("%d) %s\n", i, link)
	// }

	sitemap.ToXml(siteMap, os.Stdout)
}
