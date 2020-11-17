package parser_2

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"net/http"
)

type pageInfo struct {
	StatusCode int
	Links      map[string]int
	// Added headings
	Headings   map[string]int
}

func Parse() string {
	GQuery()
	return ""
}

func GQuery() {
	// Request the HTML page.
	res, err := http.Get("https://habr.com/ru/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".post").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find(".post__title a").Text()
		title := s.Find(".post__body .post__text").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}

func Colly() string {
	fmt.Println("STARTS")
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("habr.com/ru"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://habr.com/ru/")


	return "string"
}