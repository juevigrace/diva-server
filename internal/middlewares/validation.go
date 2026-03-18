package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

var (
	v *models.XValidator = models.NewXValidator()
)

func ValidatedHandler[T any](t T, handler func(w http.ResponseWriter, r *http.Request, body *T)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ValidateRequestBody(t, w, r)
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, "Validation error"))
			return
		}
		handler(w, r, b)
	}
}

func ValidateRequestBody[T any](t T, w http.ResponseWriter, r *http.Request) (*T, error) {
	var b *T = new(T)
	body, err := ValidateBody(b, r)
	if err != nil {
		log.Printf("Validation error: %s\n", err.Error())
		return nil, err
	}
	return body, nil
}

func ValidateBody[T any](body *T, r *http.Request) (*T, error) {
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		return nil, err
	}

	err := v.Validate(body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
