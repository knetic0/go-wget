package main

import (
	"flag"
	"path/filepath"

	golangcodes "github.com/knetic0/go-wget"
)

func main() {
	url := flag.String("u", "", "A URL for download")
	flag.Parse()

	if *url == "" {
		panic("URL flag must be valid!")
	}

	fileName := filepath.Base(*url)

	err := golangcodes.DownloadFile(*url, fileName)
	if err != nil {
		panic(err)
	}
}
