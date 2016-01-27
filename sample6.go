package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func GetPage(url string) {
	doc, _ := goquery.NewDocument(url)
	doc.Find(".s-item-container").Each(func(_ int, s *goquery.Selection) {
		title := s.Find(".s-access-title").Text()
		url, _ := s.Find(".a-link-normal").Attr("href")
		price := s.Find(".a-color-price").Text()
		img, _ := s.Find(".a-link-normal img").Attr("src")

		fmt.Println("title : ", title)
		fmt.Println("price : ", price)
		fmt.Println("url : ", url)
		fmt.Println("img", img)
	})
}

func main() {
	url := "http://www.amazon.co.jp/s/ref=sr_nr_n_1?fst=as%3Aoff&rh=n%3A2221112051%2Ck%3ANIKE&keywords=NIKE&ie=UTF8&qid=1438617594&rnid=2321267051"
	GetPage(url)
}
