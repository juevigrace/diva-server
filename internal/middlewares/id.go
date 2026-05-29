package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
)

func GetUUIDFromURL(r *http.Request, param string) (uuid.UUID, error) {
	idParam := chi.URLParam(r, param)
	if idParam == "" {
		return uuid.Nil, models.ErrIDRequired
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
