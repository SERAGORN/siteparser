package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/SERAGORN/siteparser/domain"
	"github.com/SERAGORN/siteparser/respond"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type articleHandler struct {
	articleService domain.ArticleService
	validate       *validator.Validate
	respond        *respond.Responder
}

func newArticleHandler(articleService domain.ArticleService, validate *validator.Validate, responder *respond.Responder) (*articleHandler, error) {
	if articleService == nil {
		return nil, errors.New("city handler: provided city service is nil")
	}

	if validate == nil {
		return nil, errors.New("city handler: provided validator is nil")
	}

	if responder == nil {
		return nil, errors.New("city handler: provided responder is nil")
	}

	return &articleHandler{
		articleService: articleService,
		validate:       validate,
		respond:        responder,
	}, nil
}

func (h *articleHandler) handleSearchArticles() http.HandlerFunc {
	type (
		response struct {
			ErrorReason ErrReason         `json:"errorReason,omitempty"`
			Result      *[]domain.Article `json:"result,omitempty"`
		}
	)

	return func(w http.ResponseWriter, r *http.Request) {
		searchValue := chi.URLParam(r, searchValuePattern)

		articles, err := h.articleService.SearchArticles(r.Context(), searchValue)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			if errors.Is(err, sql.ErrNoRows) {
				h.respond.NotFound(w, response{ErrorReason: articlesNotFound})
				return
			}
			fmt.Println(err)
			sentry.CaptureException(err)
			h.respond.InternalServerError(w, response{ErrorReason: internalServiceProblem})
			return
		}

		h.respond.Ok(w, response{Result: articles})
	}
}
