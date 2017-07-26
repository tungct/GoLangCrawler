package main

import (

	"flag"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"io"
	"strings"

)


var countHref int = 0
var visited = make(map[string]bool)
var queue chan string

func main() {
	flag.Parse()

	args := flag.Args()
	fmt.Println(args)
	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	queue = make(chan string)

	go func() { queue <- args[0] }()
	//ns := args[1]
	//n,error := strconv.Atoi(ns)
	count := 0
	//if error ==nil {
	for uri := range queue {

		if !visited[uri] {
			enqueue1(uri)
			count ++
		}

	}
	//}
	fmt.Println(count)

}

func getHref1(t html.Token) (ok bool, href string){
	for _,a := range t.Attr{
		if a.Key == "href"{
			href = a.Val
			ok = true
		}
	}
	return
}

func crawlHref1(b io.Reader) []string{
	links := []string{}
	z:=html.NewTokenizer(b)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return links

		case tt == html.StartTagToken:
			t := z.Token()
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			ok, url := getHref1(t)
			if !ok {
				continue
			}

			hasProto := strings.Index(url, "/") == 0

			if hasProto {
				links = append(links, url)
			}
		}
	}
	return links
}

func enqueue1(uri string) {
	countHref ++
	fmt.Println(countHref, " " + uri)

	visited[uri] = true

	resp,err := http.Get(uri)
	if err != nil{
		//fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	links := crawlHref1(resp.Body)
	for _, link := range links{

		absolute := fixUrl1(link, uri)
		if uri != "" {

			if !visited[absolute] {
				go func() { queue <- absolute }()
			}

		}

	}
	return
}

func fixUrl1(href, base string) (string) {
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