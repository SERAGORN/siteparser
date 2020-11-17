package parser

import "github.com/SERAGORN/siteparser/domain"

type parserRepository struct {
}

func NewParserRepository() (domain.ParserRepository, error) {
	return &parserRepository{}, nil
}

func (r *parserRepository) GetArticles() (*[]domain.Article, error) {

	rule := Rule{
		Url:                    "https://itproger.com/news/",
		PageStruct:             "page-",
		PostContainerRule:      ".allArticles .article",
		PostHrefRule:           "a",
		ArticleTitleRule:       ".title",
		ArticleDescriptionRule: ".article_block .txt",
		PagesNum:               1,
		HrefTemplate:           "https://itproger.com/news/",
	}

	InitParser(rule)
	return nil, nil
}
