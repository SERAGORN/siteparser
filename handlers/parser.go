package handlers

import (
	"errors"
	"github.com/SERAGORN/siteparser/domain"
	"github.com/SERAGORN/siteparser/respond"
	"github.com/getsentry/sentry-go"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type parserHandler struct {
	parserService  domain.ParserService
	articleService domain.ArticleService
	validate       *validator.Validate
	respond        *respond.Responder
}

func newParserHandler(parserService domain.ParserService, articleService domain.ArticleService, validate *validator.Validate, responder *respond.Responder) (*parserHandler, error) {
	if parserService == nil {
		return nil, errors.New("city handler: provided parser service is nil")
	}

	if validate == nil {
		return nil, errors.New("city handler: provided validator is nil")
	}

	if responder == nil {
		return nil, errors.New("city handler: provided responder is nil")
	}

	return &parserHandler{
		articleService: articleService,
		parserService:  parserService,
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
		err := r.ParseForm()
		if err != nil {
			sentry.CaptureException(err)
			h.respond.BadRequest(w, response{ErrorReason: requestParsingProblem})
			return
		}

		params := domain.ParserParams{
			Resource: r.Form.Get("resource"),
		}

		parsedArticles, err := h.parserService.GetArticles(r.Context(), params)
		if err != nil {
			sentry.CaptureException(err)
			h.respond.BadRequest(w, response{ErrorReason: ParsingProblem})
			return
		}

		err = h.articleService.SaveArticles(r.Context(), *parsedArticles)

		if err != nil {
			sentry.CaptureException(err)
			h.respond.BadRequest(w, response{ErrorReason: SaveParsedProblem})
			return
		}

		h.respond.Ok(w, response{Result: "success"})
	}
}
