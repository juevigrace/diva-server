package core

import "github.com/go-chi/chi/v5"

type Service interface{}

type Handler interface {
	Routes(r chi.Router)
}

type Module interface {
	Handler() Handler
	Service() Service
}
