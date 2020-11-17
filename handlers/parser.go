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
	parserService domain.ParserService
	validate      *validator.Validate
	respond       *respond.Responder
}

func newParserHandler(parserService domain.ParserService, validate *validator.Validate, responder *respond.Responder) (*parserHandler, error) {
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
		parserService: parserService,
		validate:      validate,
		respond:       responder,
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

		h.parserService.GetArticles(r.Context(), params)

		h.respond.Ok(w, response{Result: "success"})
	}
}
