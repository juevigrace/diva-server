package handler

import (
	"errors"
	"net/http"

	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

func handleReqError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrSessionNotFound),
		errors.Is(err, models.ErrNotAuthorized),
		errors.Is(err, models.ErrHeaderNotValid),
		errors.Is(err, models.ErrInvalidCredentials),
		errors.Is(err, models.ErrTokenNotValid),
		errors.Is(err, models.ErrBadAudience),
		errors.Is(err, models.ErrBadIssuer),
		errors.Is(err, models.ErrSessionInvalid),
		errors.Is(err, models.ErrTokenExpired):
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, err.Error()))
	case errors.Is(err, models.ErrForbidden),
		errors.Is(err, models.ErrPermissionDenied),
		errors.Is(err, models.ErrAdminAccessRequired):
		responses.WriteJSON(w, responses.RespondForbbiden(nil, err.Error()))
	case errors.Is(err, models.ErrUserNotFound),
		errors.Is(err, models.ErrActionNotFound):
		responses.WriteJSON(w, responses.RespondNotFound(nil, err.Error()))
	case errors.Is(err, models.ErrUsernameTaken),
		errors.Is(err, models.ErrEmailTaken),
		errors.Is(err, models.ErrSamePassword):
		responses.WriteJSON(w, responses.RespondConflict(nil, err.Error()))
	default:
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
	}
}
