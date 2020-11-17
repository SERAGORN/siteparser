package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/SERAGORN/siteparser/domain"
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

func GetParser(rule Rule) *[]domain.Article {
	pagesNum := rule.PagesNum
	url := rule.Url
	pageTemplate := url + rule.PageStruct
	parser := Parser{Url: url, PagesTemplate: pageTemplate, PagesNums: pagesNum, Rule: rule}
	parser.BuildPages()
	parser.GetPostUrls()

	return parser.GetPosts()
}

func (p *Parser) GetPosts() *[]domain.Article {
	articles := make(chan domain.Article)
	var resArticles []domain.Article
	var wg sync.WaitGroup
	wg.Add(len(p.PostsUrl))
	for i := range p.PostsUrl {
		go func(postUrl string) {
			defer wg.Done()
			result := p.parseArticle(postUrl)
			result.SiteUrl = p.Url
			articles <- result
		}(p.PostsUrl[i])
	}

	go func() {
		wg.Wait()
		close(articles)
	}()

	for article := range articles {
		resArticles = append(resArticles, article)
	}

	return &resArticles
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

func (p *Parser) parseArticle(url string) domain.Article {
	article := domain.Article{}
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
	article.SourceUrl = url
	article.Title = doc.Find(p.Rule.ArticleTitleRule).Text()
	article.Description = doc.Find(p.Rule.ArticleDescriptionRule).Text()
	return article
}

func (p *Parser) parsePostUrls(url string) []string {
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
		}
	})

	return postUrls
}
