package handlers

import (
	"github.com/SERAGORN/siteparser/domain"
	"github.com/SERAGORN/siteparser/middlewares"
	"github.com/SERAGORN/siteparser/respond"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strings"
)

type ErrReason string

const (
	requestParsingProblem    = ErrReason("request_parsing")
	requestValidationProblem = ErrReason("request_validation")
	internalServiceProblem   = ErrReason("service_internal")
	articlesNotFound         = ErrReason("articles_not_found")
	ParsingProblem           = ErrReason("parsing_problem")
	SaveParsedProblem        = ErrReason("save_parsed_problem")
)

const (
	searchValuePattern = "searchValue"
)

type RouterDependencies struct {
	ArticleService domain.ArticleService
	MySql          *sqlx.DB
	Validate       *validator.Validate
}

func MakeRoutes(routerDependencies *RouterDependencies) (http.Handler, error) {

	responder, err := respond.NewResponder()
	if err != nil {
		return nil, err
	}

	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	router := chi.NewRouter()

	router.Use(sentryHandler.Handle)
	router.Use(middlewares.SetHeaders()...)
	router.Use(CaselessMatcher)

	articleHandler, err := newArticleHandler(routerDependencies.ArticleService, routerDependencies.Validate, responder)

	if err != nil {
		return nil, err
	}

	parserHandler, err := newParserHandler(routerDependencies.ArticleService, routerDependencies.Validate, responder)

	if err != nil {
		return nil, err
	}

	router.Route("/api", func(r chi.Router) {
		r.Route("/search", func(r chi.Router) {
			r.Get("/{"+searchValuePattern+"}", articleHandler.handleSearchArticles())
		})
		r.Post("/parse", parserHandler.handleInitParser())
	})

	return router, err
}

func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
