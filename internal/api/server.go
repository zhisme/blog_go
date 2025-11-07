package api

import (
	"backend-go/internal/interfaces"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	router               *chi.Mux
	mailingListRepository interfaces.MailingListRepository
}

// ServeHTTP implements http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(addr string) error {
	log.Default().Printf("api server started on %s\n", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func NewApiServer(mailingListRepo interfaces.MailingListRepository) *Server {
	srv := &Server{
		router:               chi.NewRouter(),
		mailingListRepository: mailingListRepo,
	}

	srv.router.Use(middleware.Logger)
	srv.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:1313", "https://zhisme.com/"},
		AllowedMethods:   []string{"POST", "DELETE"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	srv.router.Post("/mailing_list", srv.createMailingList)

	log.Default().Println("api server initialized")

	return srv
}
