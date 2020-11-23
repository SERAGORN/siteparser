package parser

import (
	"context"
	"github.com/SERAGORN/siteparser/domain"
	"github.com/SERAGORN/siteparser/internal/parser"
)

type parserRepository struct {
}

func NewParserRepository() (domain.ParserRepository, error) {
	return &parserRepository{}, nil
}

func (r *parserRepository) GetArticles(ctx context.Context, parserParams domain.ParserParams) (*[]domain.Article, error) {
	// Пример парсинга https://itproger.com/news/
	//parserParams= domain.ParserParams{
	//	Url:                    "https://itproger.com/news/",
	//	PageStruct:             "/page-",
	//	PostContainerRule:      ".allArticles .article",
	//	PostHrefRule:           "a",
	//	ArticleTitleRule:       ".title",
	//	ArticleDescriptionRule: ".article_block .txt",
	//	PagesNum:               2,
	//	HrefTemplate:           "https://itproger.com/news",
	//}

	pagesNum := parserParams.PagesNum
	url := parserParams.Url

	parser := parser.Parser{Url: url, PagesNums: pagesNum, Rule: parser.Rule{
		Url:                    parserParams.Url,
		PostContainerRule:      parserParams.PostContainerRule,
		PostHrefRule:           parserParams.PostHrefRule,
		ArticleTitleRule:       parserParams.ArticleTitleRule,
		ArticleDescriptionRule: parserParams.ArticleDescriptionRule,
		PagesNum:               parserParams.PagesNum,
		PagesStruct:            parserParams.PageStruct,
		HrefTemplate:           parserParams.HrefTemplate,
	}}

	parser.BuildPages()
	parser.GetPostUrls()
	return parser.GetPosts(), nil
}
