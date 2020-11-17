package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"sync"
)

type Parser struct {
	Url      string
	PagesTemplate string
	PagesNums int
	PagesUrl    []string
	PostsUrl []string
	Articles []Article
}

type Article struct {
	Url         string
	Title       string
	Description string
}


func InitParser() {
	pagesNum := 10
	url := "https://habr.com/ru/"
	pageTemplate := url + "page"
	parser := Parser{Url: url, PagesTemplate: pageTemplate, PagesNums: pagesNum}
	parser.BuildPages()
	parser.GetPostUrls()
	parser.GetPosts()
	fmt.Println(parser.Articles)
}

func (p *Parser) GetPosts() {
	articles := make(chan Article)

	var wg sync.WaitGroup
	wg.Add(len(p.PostsUrl))
	for i := range p.PostsUrl {
		go func (postUrl string) {
			defer wg.Done()
			result := parseArticle(postUrl)
			articles <- result
		}(p.PostsUrl[i])
	}

	go func() {
		wg.Wait()
		close(articles)
	}()


	for article := range articles {
		p.Articles = append(p.Articles, article)
	}
}

func (p *Parser) BuildPages() {
	for i := 0; i < p.PagesNums; i++ {
		p.PagesUrl = append(p.PagesUrl, fmt.Sprintf("%s%d/", p.PagesTemplate, p.PagesNums+1))
	}
}

func (p *Parser) GetPostUrls() {
	postUrls := make(chan string)

	var wg sync.WaitGroup
	wg.Add(len(p.PagesUrl))

	for i := range p.PagesUrl {
		go func (pageUrl string) {

			defer wg.Done()

			results := parsePostUrls(pageUrl)

			for i := range results {
				postUrls <- results[i]
			}

		}(p.PagesUrl[i])
	}

	go func() {
		wg.Wait()
		close(postUrls)
	}()

	for url := range postUrls {
		p.PostsUrl = append(p.PostsUrl, url)
	}
}

func parseArticle(url string) Article {
	article := Article{}
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s, %s", res.StatusCode, res.Status, url)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	article.Url = url
	article.Title = doc.Find(".post__title-text").Text()
	article.Description = doc.Find(".post__body .post__text").Text()
	fmt.Println(article)
	return article
}

func parsePostUrls(url string) []string {

	var postUrls []string

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s, %s", res.StatusCode, res.Status, url)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".post").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band, ok := s.Find(".post__title a").Attr("href")
		if ok {
			postUrls = append(postUrls, band)
		}
	})

	return postUrls
}