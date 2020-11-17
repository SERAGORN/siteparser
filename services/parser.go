package services

import (
	"context"
	"github.com/SERAGORN/siteparser/domain"
)

type parserService struct {
	parserRepository domain.ParserRepository
}

func NewParserService(parserRepository domain.ParserRepository) (domain.ParserService, error) {
	return &parserService{
		parserRepository: parserRepository,
	}, nil
}

func (s *parserService) GetArticles(ctx context.Context, params domain.ParserParams) (*[]domain.Article, error){

	return s.parserRepository.GetArticles()
}

