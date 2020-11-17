package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/SERAGORN/siteparser/domain"
	"github.com/SERAGORN/siteparser/respond"
	"github.com/getsentry/sentry-go"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
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

func (h *articleHandler) handleGetArticles() http.HandlerFunc {
	type (
		response struct {
			ErrorReason ErrReason       `json:"errorReason,omitempty"`
			Result      *domain.Article `json:"result,omitempty"`
		}
	)

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			sentry.CaptureException(err)
			h.respond.BadRequest(w, response{ErrorReason: requestParsingProblem})
			return
		}

		params := domain.GetArticleParams{
			Id: r.Form.Get("id"),
		}

		articleId, err := strconv.ParseInt(params.Id, 10, 64)
		if err == nil {
			fmt.Printf("%d of type %T", articleId, articleId)
		}

		if err := h.validate.Struct(params); err != nil {
			sentry.CaptureException(err)
			h.respond.BadRequest(w, response{ErrorReason: requestValidationProblem})
			return
		}

		articles, err := h.articleService.GetArticle(r.Context(), articleId)
		if err != nil {

			if errors.Is(err, context.Canceled) {
				return
			}

			if errors.Is(err, sql.ErrNoRows) {
				h.respond.NotFound(w, response{ErrorReason: articlesNotFound})
				return
			}
			sentry.CaptureException(err)
			h.respond.InternalServerError(w, response{ErrorReason: internalServiceProblem})
			return
		}

		h.respond.Ok(w, response{Result: articles})
	}
}
