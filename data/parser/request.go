package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"sync"
)

type Parser struct {
	Rule          Rule
	Url           string
	PagesTemplate string
	PagesNums     int
	PagesUrl      []string
	PostsUrl      []string
	Articles      []Article
}

type Rule struct {
	Url                    string
	PageStruct             string
	PostContainerRule      string
	PostHrefRule           string
	ArticleTitleRule       string
	ArticleDescriptionRule string
	PagesNum               int
	HrefTemplate           string
}

type Article struct {
	Url         string
	Title       string
	Description string
}

func InitParser(rule Rule) {
	pagesNum := rule.PagesNum
	url := rule.Url
	pageTemplate := url + rule.PageStruct
	parser := Parser{Url: url, PagesTemplate: pageTemplate, PagesNums: pagesNum, Rule: rule}
	parser.BuildPages()
	parser.GetPostUrls()
	parser.GetPosts()
}

func (p *Parser) GetPosts() {
	articles := make(chan Article)

	var wg sync.WaitGroup
	wg.Add(len(p.PostsUrl))
	for i := range p.PostsUrl {
		go func(postUrl string) {
			defer wg.Done()
			result := p.parseArticle(postUrl)
			articles <- result
		}(p.PostsUrl[i])
	}

	go func() {
		wg.Wait()
		close(articles)
	}()

	for article := range articles {
		p.Articles = append(p.Articles, article)
		fmt.Println(article.Title)
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
		go func(pageUrl string) {

			defer wg.Done()

			results := p.parsePostUrls(pageUrl)

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

func (p *Parser) parseArticle(url string) Article {
	article := Article{}
	url = p.Rule.HrefTemplate + url
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
	article.Title = doc.Find(p.Rule.ArticleTitleRule).Text()
	article.Description = doc.Find(p.Rule.ArticleDescriptionRule).Text()
	fmt.Println(article, p.Rule.ArticleTitleRule)
	return article
}

func (p *Parser) parsePostUrls(url string) []string {
	fmt.Println(url)
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
	doc.Find(p.Rule.PostContainerRule).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band, ok := s.Find(p.Rule.PostHrefRule).Attr("href")
		if ok {
			postUrls = append(postUrls, band)
			fmt.Println(band)
		}
	})

	return postUrls
}
