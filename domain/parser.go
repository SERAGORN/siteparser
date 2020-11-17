package domain

import "context"

type ParserParams struct {
	Resource string
}

type ParserRepository interface {
	GetArticles() (*[]Article, error)
}

type ParserService interface {
	GetArticles(ctx context.Context, params ParserParams) (*[]Article, error)
}