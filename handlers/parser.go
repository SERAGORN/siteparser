package handlers

import (
	"encoding/json"
	"errors"
	"github.com/SERAGORN/siteparser/domain"
	"github.com/SERAGORN/siteparser/respond"
	"github.com/getsentry/sentry-go"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type parserHandler struct {
	articleService domain.ArticleService
	validate       *validator.Validate
	respond        *respond.Responder
}

func newParserHandler(articleService domain.ArticleService, validate *validator.Validate, responder *respond.Responder) (*parserHandler, error) {

	if validate == nil {
		return nil, errors.New("city handler: provided validator is nil")
	}

	if responder == nil {
		return nil, errors.New("city handler: provided responder is nil")
	}

	return &parserHandler{
		articleService: articleService,
		validate:       validate,
		respond:        responder,
	}, nil
}

func (h *parserHandler) handleInitParser() http.HandlerFunc {
	type (
		response struct {
			ErrorReason ErrReason `json:"errorReason,omitempty"`
			Result      string    `json:"result,omitempty"`
		}
	)

	return func(w http.ResponseWriter, r *http.Request) {
		request := domain.ParserParams{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			sentry.CaptureException(err)
			h.respond.BadRequest(w, response{ErrorReason: requestParsingProblem})
			return
		}

		if err := h.validate.Struct(request); err != nil {
			sentry.CaptureException(err)
			h.respond.BadRequest(w, response{ErrorReason: requestValidationProblem})
			return
		}

		err := h.articleService.ParseArticles(r.Context(), request)
		if err != nil {
			sentry.CaptureException(err)
			h.respond.BadRequest(w, response{ErrorReason: ParsingProblem})
			return
		}

		h.respond.Ok(w, response{Result: "success"})
	}
}
