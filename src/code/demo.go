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
	"sync"
)

var mutex = &sync.Mutex{}
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
	//count := 0
	//if error ==nil {
		for uri := range queue {
			mutex.Lock()
			if !visited[uri] {
				go enqueue(uri)
				//count ++
			}
			mutex.Unlock()
		}
	//}
	//fmt.Println(count)

}

func getHref(t html.Token) (ok bool, href string){
	for _,a := range t.Attr{
		if a.Key == "href"{
			href = a.Val
			ok = true
		}
	}
	return
}

func crawlHref(b io.Reader) []string{
	links := []string{}
	z:=html.NewTokenizer(b)
	for ;; {
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

			ok, url := getHref(t)
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

func enqueue(uri string) {
	countHref ++
        fmt.Println(countHref, " " + uri)
	mutex.Lock()
	visited[uri] = true
	mutex.Unlock()
	resp,err := http.Get(uri)
	if err != nil{
		//fmt.Println(err)
		return
	}
        defer resp.Body.Close()
	links := crawlHref(resp.Body)
	for _, link := range links{

		absolute := fixUrl(link, uri)
		if uri != "" {
			mutex.Lock()
			if !visited[absolute] {
				go func() { queue <- absolute }()
			}
			mutex.Unlock()
		}

	}
	return
}

func fixUrl(href, base string) (string) {
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