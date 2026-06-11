package responses

import (
	"errors"
	"net/http"

	"github.com/juevigrace/diva-server/pkg/errs"
)

func HandleReqError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrSessionNotFound),
		errors.Is(err, errs.ErrNotAuthorized),
		errors.Is(err, errs.ErrHeaderNotValid),
		errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrTokenNotValid),
		errors.Is(err, errs.ErrBadAudience),
		errors.Is(err, errs.ErrBadIssuer),
		errors.Is(err, errs.ErrSessionInvalid),
		errors.Is(err, errs.ErrTokenExpired):
		WriteJSON(w, RespondUnauthorized(nil, err.Error()))
	case errors.Is(err, errs.ErrForbidden),
		errors.Is(err, errs.ErrPermissionDenied),
		errors.Is(err, errs.ErrAdminAccessRequired):
		WriteJSON(w, RespondForbbiden(nil, err.Error()))
	case errors.Is(err, errs.ErrUserNotFound),
		errors.Is(err, errs.ErrActionNotFound):
		WriteJSON(w, RespondNotFound(nil, err.Error()))
	case errors.Is(err, errs.ErrUsernameTaken),
		errors.Is(err, errs.ErrEmailTaken),
		errors.Is(err, errs.ErrSamePassword):
		WriteJSON(w, RespondConflict(nil, err.Error()))
	default:
		WriteJSON(w, RespondBadRequest(nil, err.Error()))
	}
}
