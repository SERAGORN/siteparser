package articlesrvc

import (
	"context"
	"github.com/SERAGORN/siteparser/domain"
)

type articleService struct {
	articleRepository domain.ArticleRepository
	parserRepository  domain.ParserRepository
}

func NewArticleService(articleRepository domain.ArticleRepository, parserRepository domain.ParserRepository) (domain.ArticleService, error) {
	return &articleService{
		articleRepository: articleRepository,
		parserRepository:  parserRepository,
	}, nil
}

func (s *articleService) ParseArticles(ctx context.Context, params domain.ParserParams) error {
	articles, err := s.parserRepository.GetArticles(ctx, params)
	if err != nil {
		return err
	}

	return s.articleRepository.SaveArticles(ctx, *articles)
}

func (s *articleService) GetArticle(ctx context.Context, articleId int64) (*domain.Article, error) {

	article, err := s.articleRepository.GetArticle(ctx, articleId)

	return article, err
}

func (s *articleService) SearchArticles(ctx context.Context, searchText string) (*[]domain.Article, error) {
	return s.articleRepository.SearchArticles(ctx, searchText)
}

func (s *articleService) SaveArticles(ctx context.Context, articles []domain.Article) error {

	err := s.articleRepository.SaveArticles(ctx, articles)

	return err
}
