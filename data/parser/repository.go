package parser

import "github.com/SERAGORN/siteparser/domain"

type parserRepository struct {

}

func NewParserRepository() (domain.ParserRepository, error) {
	return &parserRepository{}, nil
}

func (r *parserRepository) GetArticles() (*[]domain.Article, error) {

	InitParser()
	return nil, nil
}