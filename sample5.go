package main
import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func GetPage(url string) {
	doc, _ := goquery.NewDocument(url)
	doc.Find(".listLink a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if url != "#listtop" {
			keyword := s.Text()
			if len(keyword) > 0 {
				fmt.Println(keyword)
			}
		}
	})
}

func main() {
	tags := []string{"ア", "カ", "サ", "タ", "ナ", "ハ", "マ", "ヤ", "ラ", "ワ", "英字", "数字"}
	for _, tag := range tags {
		url := "http://artscape.jp/artword/index.php/wordlist_gyo?tag=" + tag
		fmt.Println(url)
		GetPage(url)
		break
	}
}
