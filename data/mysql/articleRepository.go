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
	searchArticle     = "SELECT * FROM tbl_article WHERE MATCH (title,description) AGAINST (?);"
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

func (r *articleRepository) SearchArticles(ctx context.Context, searchText string) (*[]domain.Article, error) {
	var articles []domain.Article
	fmt.Println(searchText)
	err := r.db.SelectContext(ctx, &articles, searchArticle, searchText)
	fmt.Println(articles)
	if err != nil {
		return nil, err
	}
	return &articles, err
}

// SaveArticles save multiple articles returns only err
func (r *articleRepository) SaveArticles(ctx context.Context, articles []domain.Article) error {
	insertArticleString := insertArticle
	values := []interface{}{}

	for i := range articles {
		insertArticleString += "(?,?,?,?),"
		values = append(values, articles[i].Title, articles[i].Description, articles[i].SourceUrl, articles[i].SiteUrl)
	}

	insertArticleString = insertArticleString[0 : len(insertArticleString)-1]

	stmt, err := r.db.Prepare(insertArticleString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}
	return nil
}
