package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/storage"
)

type ServerRouter struct {
	Chi *chi.Mux
	db  storage.Storage
}

func NewServerRouter(db storage.Storage) *ServerRouter {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5000", "http://127.0.0.1:5000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		res := responses.RespondNotFound(nil, "Route not found")
		responses.WriteJSON(w, res)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		res := responses.RespondNotFound(nil, "Method not allowed for this route")
		responses.WriteJSON(w, res)
	})

	// queries := db.Queries()

	return &ServerRouter{
		Chi: r,
		db:  db,
	}
}
