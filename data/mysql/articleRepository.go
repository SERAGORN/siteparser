package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/SERAGORN/siteparser/domain"
	"github.com/jmoiron/sqlx"
)

type articleRepository struct {
	db *sqlx.DB
}

const (
	selectArticleById = "select id, title, description from tbl_article where id = ? ;"
	insertArticle     = "insert into tbl_article (title, description ,source_url, site_url) values"
)

var ErrNilDBHandle = errors.New("provided db handle is nil")

func NewArticleRepository(db *sqlx.DB) (domain.ArticleRepository, error) {
	if db == nil {
		return nil, ErrNilDBHandle
	}
	return &articleRepository{db: db}, nil
}

func (r *articleRepository) GetArticle(ctx context.Context, articleId int64) (*domain.Article, error) {

	article := domain.Article{}

	err := r.db.GetContext(ctx, &article, selectArticleById, articleId)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) SaveArticle(ctx context.Context, article domain.Article) error {

	_, err := r.db.ExecContext(ctx, insertArticle, article.Title, article.Description)

	return err
}

func (r *articleRepository) SaveArticles(ctx context.Context, articles []domain.Article) error {

	var res []interface{}
	insertArticleString := insertArticle
	fmt.Println(len(articles))
	for i := range articles {
		if i == 0 {
			insertArticleString = insertArticleString + " (?,?,?,?)"
		} else {
			insertArticleString = insertArticleString + ", (?,?,?,?)"
		}

		res = append(res, articles[i].Title, "", articles[i].SourceUrl, articles[i].SiteUrl)
	}

	insertArticleString = insertArticleString + ";"

	_, err := r.db.ExecContext(ctx, insertArticle, res...)
	return err
}
