package main

import (
	"io"
	"log"
	"os"
	"strings"
)

type Url struct {
	Link  string
	Title string
}

var (
	UrlList = []Url{}
)

/*
* Retreive list of URLs from file
 */
func init() {

	// open URL list file
	f, err := os.Open("urls/url_list.txt")

	if err != nil {
		log.Println("Error opening file:", err)
	}
	defer f.Close()

	// allocate storage buffer
	buf := make([]byte, 1024)

	var FileContents string
	FileContents = ""

	// read file until EOF
	for {

		n, err := f.Read(buf)

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}

		if n > 0 {
			FileContents = string(buf[:n])
		}
	}
	for _, v := range strings.Split(strings.ReplaceAll(FileContents, "\r\n", "\n"), "\n") {
		UrlList = append(UrlList, Url{
			Title: v,
			Link: v,
		})
	}
	

}
