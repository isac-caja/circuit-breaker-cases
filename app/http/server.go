package http

import (
	gohttp "net/http"

	"github.com/go-chi/chi"
)

func NewRouter(handler Handler) chi.Router {
	router := chi.NewRouter()

	router.Get("/ping", handler.HandlePing)
	router.Get("/process", handler.HandleProcess)
	return router
}

func NewServer(router chi.Router) *gohttp.Server {
	port := ":8080"
	return &gohttp.Server{
		Addr:    port,
		Handler: router,
	}
}
