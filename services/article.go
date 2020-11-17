package services

import (
	"context"
	"github.com/SERAGORN/siteparser/domain"
)

type articleService struct {
	articleRepository domain.ArticleRepository
}

func NewArticleService(articleRepository domain.ArticleRepository) (domain.ArticleService, error) {
	return &articleService{
		articleRepository: articleRepository,
	}, nil
}

func (s *articleService) GetArticle (ctx context.Context, articleId int64) (*domain.Article, error) {

	article, err := s.articleRepository.GetArticle(ctx, articleId)

	return article, err
}
