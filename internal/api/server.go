package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router   *chi.Mux
}

func (s *Server) ListenAndServe(addr string) error {
	log.Default().Printf("api server started on %s\n", addr)
	if err := http.ListenAndServe(addr, s.router); err != nil {
		return err
	}
	return nil
}

func NewApiServer() *Server {
	srv := &Server{
		router: chi.NewRouter(),
	}

	srv.router.Use(middleware.Logger)
	srv.router.Post("/mailing_list", srv.createMailingList)

	log.Default().Println("api server initialized")

	return srv
}
