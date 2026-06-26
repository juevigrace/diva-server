package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/juevigrace/diva-server/pkg/validator"
)

var (
	v = validator.GetInstance()
)

func ValidateBody[T any](body *T, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		return err
	}

	err := v.Validate(body)
	if err != nil {
		return err
	}

	return nil
}
