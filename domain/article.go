package domain

import "context"

type ArticleService interface {
	GetArticle(ctx context.Context, articleId int64) (*Article, error)
	SaveArticles(ctx context.Context, article []Article) error
}

type ArticleRepository interface {
	GetArticle(ctx context.Context, articleId int64) (*Article, error)
	SaveArticle(ctx context.Context, article Article) error
	SaveArticles(ctx context.Context, articles []Article) error
}

type Article struct {
	ID          int64  `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	SiteUrl     string `json:"site_url" db:"site_url"`
	SourceUrl   string `json:"source_url" db:"source_url"`
}

type GetArticleParams struct {
	Id string `json:"id"`
}
