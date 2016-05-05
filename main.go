package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	BaseURL = "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
)

type Product struct {
	Title string
	Size  int64
	Price float64
	Desc  string
}

func main() {
	urls, err := fetchBaseDocument(BaseURL)
	if err != nil {
		log.Fatal(err)
	}

	res := make(chan Product)
	errs := make(chan error)

	for _, url := range urls {
		go fetchProduct(url, res, errs) // start a goroutine
	}

	products := []Product{}

	for i := 0; i < len(urls); i++ {
		select {
		case resp := <-res:
			products = append(products, resp)
		case errs := <-errs:
			fmt.Printf("errs = %+v\n", errs)
		}
	}

	total := 0.0
	for _, p := range products {
		total = total + p.Price
	}

	result := map[string]interface{}{
		"results": products,
		"total":   fmt.Sprintf("%.2f", total),
	}

	json, _ := json.MarshalIndent(result, "", "    ")
	fmt.Printf("%s\n", json)
}

func fetchBaseDocument(url string) (urls []string, err error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return urls, err
	}

	doc.Find(".productLister h3 a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			urls = append(urls, href)
		}
	})

	return urls, err
}

func fetchProduct(url string, res chan<- Product, errc chan<- error) {
	resp, err := http.Get(url)
	if err != nil {
		errc <- err
	}

	product := Product{}
	product.Size = resp.ContentLength

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		errc <- err
	}

	product.Title = doc.Find(".pdp h1").First().Text()
	product.Desc = strings.Trim(doc.Find(".mainProductInfo .productText").First().Text(), " \n")

	price := doc.Find(".addToTrolleytabBox .pricing .pricePerUnit").First().Text()
	price = strings.Replace(price, "\n", "", -1)
	price = strings.Replace(price, "£", "", -1)
	price = strings.Replace(price, "/unit", "", -1)
	f, err := strconv.ParseFloat(price, 64)
	if err != nil {
		errc <- err
	}

	product.Price = f

	res <- product
}
