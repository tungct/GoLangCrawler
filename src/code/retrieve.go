package main    // Note: this is a new temporary file, not our crawl.go

import (
	"fmt"
	"net/http"    // this is the 'http' package we'll be using to retrieve a page
	   // we'll only use 'ioutil' to maek reading and printing
	"github.com/jackdanger/collectlinks"
	"net/url"
)               // the html page a little easier in this example.

var count int = 0

func main() {
	res, err := http.Get("http://genk.vn/iphone-8-se-co-4-mau-bao-gom-lua-chon-mat-kinh-guong-2017071008531728.chn")
	fmt.Println("Error in http request : ", err)
        base := "http://genk.vn/iphone-8-se-co-4-mau-bao-gom-lua-chon-mat-kinh-guong-2017071008531728.chn"
	links := collectlinks.All(res.Body)

	queue := make(chan string)
	go func() {queue<-base}()
	for _, link := range(links){
		absolute := fixUrl3(link, base)
		count ++
		fmt.Println(absolute)

	}
	fmt.Println(count)

}


func fixUrl3(href, base string) (string) {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}