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
	Rule      Rule
	Url       string
	PagesNums int
	PagesUrl  []string
	PostsUrl  []string
}

// Rule of parse news
type Rule struct {
	//Url Page of news example"https://itproger.com/news/"
	Url string
	//PostContainerRule html container location from father to child,
	//using tags html classes and id,
	//the more nesting the more accurate example: ".allArticles .article"
	PostContainerRule string
	//PostHrefRule html href of post container, usually: tag <a>
	//example: "a"
	PostHrefRule string
	//ArticleTitleRule  html post title location from article page,
	//using tags html classes and id,
	//the more nesting the more accurate example: ".title"
	ArticleTitleRule string
	//ArticleTitleRule  html post description location from article page,
	//using tags html classes and id,
	//the more nesting the more accurate example: ".title"
	ArticleDescriptionRule string
	//PagesNum nums of main news site pages count "https://itproger.com/news/1" "https://itproger.com/news/2" etc...
	//parsed only two pages of articles example: 2
	PagesNum int
	//HrefTemplate template of pages home
	//example "https://itproger.com/news/"
	HrefTemplate string
	//PagesStruct appends of pages index and build url with HrefTemplate
	//Example: "%d/"
	PagesStruct string
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
	for i := 1; i < p.PagesNums; i++ {
		p.PagesUrl = append(p.PagesUrl, fmt.Sprintf("%s%s%d", p.Rule.HrefTemplate, p.Rule.PagesStruct, i))
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
	url = p.Rule.HrefTemplate + "/" + url
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("status code error: %d %s, %s", res.StatusCode, res.Status, url)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("status code error: %d %s, %s", res.StatusCode, res.Status, url)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
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
