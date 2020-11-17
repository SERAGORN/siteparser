package middlewares

import (
"github.com/go-chi/chi/middleware"
"net/http"
)

func SetHeaders() (middlewares []func(http.Handler) http.Handler) {
	middlewares = append(middlewares, middleware.SetHeader("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type, Content-Length, Accept-Encoding, X-Forwarded-For, X-Real-IP, Authorization"))
	middlewares = append(middlewares, middleware.SetHeader("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS"))
	middlewares = append(middlewares, middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	middlewares = append(middlewares, middleware.SetHeader("Access-Control-Allow-Credentials", "true"))
	middlewares = append(middlewares, middleware.AllowContentType("application/x-www-form-urlencoded", "application/json"))

	return middlewares
}