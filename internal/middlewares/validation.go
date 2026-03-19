package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/juevigrace/diva-server/internal/models"
)

var (
	v *models.XValidator = models.NewXValidator()
)

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
