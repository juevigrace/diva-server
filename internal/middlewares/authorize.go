package middlewares

import (
	"context"
	"maps"
	"net/http"
	"slices"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/pkg/errs"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

func RequireRole(roles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc, err := GetRequestContext(r.Context())
			if err != nil {
				responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
				return
			}

			if rc.Session.User.Role == models.ROLE_ADMIN {
				next.ServeHTTP(w, r)
				return
			}

			if !slices.Contains(roles, rc.Session.User.Role) {
				responses.WriteJSON(w, responses.RespondForbbiden(nil, errs.ErrForbidden.Error()))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type RequireOwnerParams struct {
	// TODO: consider replacing this
	UrlParams []string
	Perms     []models.PermissionAction
}

func RequireResourceOwner(
	params *RequireOwnerParams,
	load func(ctx context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc, err := GetRequestContext(r.Context())
			if err != nil {
				responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
				return
			}

			if rc.Session.User.Role == models.ROLE_ADMIN {
				next.ServeHTTP(w, r)
				return
			}

			resParams := make([]string, len(params.UrlParams))
			for i := range params.UrlParams {
				param := chi.URLParam(r, params.UrlParams[i])
				if param == "" {
					responses.WriteJSON(w, responses.RespondBadRequest(nil, errs.ErrParamRequired.Error()))
					return
				}
				resParams[i] = param
			}

			for i := range params.Perms {
				if perm, exists := rc.Session.User.Permissions[params.Perms[i]]; exists && perm.Granted {
					if perm.ExpiresAt == nil || time.UnixMilli(*perm.ExpiresAt).After(time.Now().UTC()) {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			cached, ok := load(r.Context(), rc.Session.User.ID, resParams)
			if !ok {
				responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrForbidden.Error()))
				return
			}

			if cached != nil {
				maps.Copy(rc.Cache, cached)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequirePermission(actions ...models.PermissionAction) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc, err := GetRequestContext(r.Context())
			if err != nil {
				responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
				return
			}

			if rc.Session.User.Role == models.ROLE_ADMIN {
				next.ServeHTTP(w, r)
				return
			}

			for i := range actions {
				perm, exists := rc.Session.User.Permissions[actions[i]]
				if !exists || !perm.Granted {
					responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrPermissionDenied.Error()))
					return
				}
				if perm.ExpiresAt != nil && time.UnixMilli(*perm.ExpiresAt).Before(time.Now().UTC()) {
					responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrPermissionDenied.Error()))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
