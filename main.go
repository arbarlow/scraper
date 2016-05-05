package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	BaseURL = "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
)

func main() {
	doc, err := goquery.NewDocument(BaseURL)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".productLister h3 a").Each(func(i int, s *goquery.Selection) {
		fmt.Printf(s.Attr("href"))
	})

	// now := time.Now()
	// res := make(chan http.Response)
	// errs := make(chan error)

	// for _, url := range os.Args[1:] {
	// go fetch(BaseURL, res, errs) // start a goroutine
	// // }
	// for {
	// 	select {
	// 	case resp := <-res:
	// 		parseBasePage(resp.Body)
	// 	case errs := <-res:
	// 		fmt.Print(errs)
	// 	}
	// }

	// fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

// func parseBasePage(r io.ReadCloser) {
// 	z := html.Parse(r)

// 	for {
// 		tt := z.Next()

// 	}
// }

func fetch(url string, res chan<- http.Response, errc chan<- error) {
	resp, err := http.Get(url)
	if err != nil {
		errc <- err
		return
	}

	res <- *resp
}
