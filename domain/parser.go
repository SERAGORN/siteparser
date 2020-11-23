package domain

import "context"

type ParserParams struct {
	Url                    string `json:"url"`
	PageStruct             string `json:"pageStruct"`
	PostContainerRule      string `json:"postContainerRule"`
	PostHrefRule           string `json:"postHrefRule"`
	ArticleTitleRule       string `json:"articleTitleRule"`
	ArticleDescriptionRule string `json:"articleDescriptionRule"`
	PagesNum               int    `json:"pagesNum"`
	HrefTemplate           string `json:"hrefTemplate"`
}

type ParserRepository interface {
	GetArticles(ctx context.Context, params ParserParams) (*[]Article, error)
}
