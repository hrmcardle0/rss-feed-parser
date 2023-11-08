package main

import (
	"flag"
	//"fmt"
	"log"
	"sync"
	"os"
	"strings"
	//"time"
	"github.com/mmcdole/gofeed"
)

var (
	redirectFlag = flag.Int("c", 0, "True or false, check redirects")
	outputFile string = "output/output.html"
)

// waitgroup for concurrency management
var wg sync.WaitGroup

/*
* init() runs before main when package is loading, here the parsing of CLI options happens
 */
func init() {
	flag.Parse()
}

type Response struct {
	Title      string
	Link       string
	PubDate	string
	StatusCode int
}

type FeedItem struct {
	Title string
	Link string
	PubDate string
}

/*
* Main entry point function. Declare an HTTP client, retrive information
* from URLs from go-routines for parsing
 */
func main() {
	log.Println("Scraper Starting...")

	FeedItemList := []FeedItem{}

	// get client
	var Client *gofeed.Parser
	Client = GenerateClient()

	// init channels
	ResponseChannel := make(chan *Response)

	for _, url := range UrlList {
		log.Println("Getting:", url.Title)
		//wg.Add(1)
		
		go GetResponse(Client, url.Link, ResponseChannel)
	}

	// loop through channel to generate list of items containing title & link
	for res := range ResponseChannel {
		//log.Println(res.Title)
		FeedItemList = append(FeedItemList, FeedItem{
			Title: res.Title,
			Link: res.Link,
			PubDate: res.PubDate,
		})
	}
	
	var err error
	// ready output file
	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString("<html><title>RSS Feed Detailed</title><body>\n")
	if err != nil {
		panic(err)
	}

	// loop through feed list, adding each entry to output file
	for i, item := range FeedItemList {
		// if we are at the last entry
		var StrToWrite string
		if i == (len(FeedItemList) - 1) {
			StrToWrite = "<i>" + strings.TrimSuffix(item.PubDate, " +0000") + "</i>" + "<b>&nbsp;" + item.Title + "</b>\n\n<br><br>" + "<a href=\"" + item.Link + "\">" + item.Link + "</a><br><br><br>"
		} else {
			StrToWrite = "<i>" + strings.TrimSuffix(item.PubDate, " +0000") + "</i>" + "<b>&nbsp;" + item.Title + "</b>\n\n<br><br>" + "<a href=\"" + item.Link + "\">" + item.Link + "</a><br><br><br>" + "\n\n"
		}
		_, err := f.WriteString(StrToWrite)
		if err != nil {
			panic(err)
		}

	}

	_, err = f.WriteString("</body></html>")
	if err != nil {
		panic(err)
	}

	log.Println("Scrape Complete")

}

func GetResponse(client *gofeed.Parser, url string, rc chan *Response) {
	//time.Sleep(1000 * time.Millisecond)
	feed, err := client.ParseURL(url);
	if err != nil {
		panic(err)
	}

	//defer wg.Done()
	for _, item := range feed.Items {
		rc <- &Response{
			Title:      item.Title,
			Link:       item.Link,
			PubDate: 	item.Published,
			StatusCode: 200,
		}
	}

	// close channel
	close(rc)
}
