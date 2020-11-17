package handlers

import (
	"github.com/SERAGORN/siteparser/middlewares"
	"github.com/SERAGORN/siteparser/domain"
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
	emptyRequestProblem      = ErrReason("empty_request")
	requestParsingProblem    = ErrReason("request_parsing")
	requestDecodingProblem   = ErrReason("request_decoding")
	requestValidationProblem = ErrReason("request_validation")
	userUnauthorized         = ErrReason("user_unauthorized")
	internalServiceProblem   = ErrReason("service_internal")
	emptySearchParams        = ErrReason("empty_search_params")
	testNotUnique            = ErrReason("test_not_unique")
	calculatorNotFound       = ErrReason("articles_not_found")
	calculatorNotValid       = ErrReason("calculator_not_valid")

	messageTestNotUnique        = "Test with this order type is already enabled"
	messageTestCreated          = "Test successfully created"
	messageTestStopped          = "Test successfully stopped"
	messageCalculatorNotValid   = "Validation error"
	messageServiceError         = "Service error"
	messageStatsAddedToQueue    = "Stats added to queue"
	messageStatsOnMakingProcess = "Stats on making process"
	messageEmptyCityID          = "Empty cityID"
)

type RouterDependencies struct {
	ArticleService domain.ArticleService
	ParserService  domain.ParserService
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

	parserHandler, err := newParserHandler(routerDependencies.ParserService, routerDependencies.Validate, responder)

	if err != nil {
		return nil, err
	}

	router.Route("/api",  func(r chi.Router) {
		r.Get("/article", articleHandler.handleGetArticles())
		r.Get("/parse", parserHandler.handleInitParser())
	})

	return router, err
}

func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}