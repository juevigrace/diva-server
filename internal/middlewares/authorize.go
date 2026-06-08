package middlewares

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/errs"
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

// TODO: should i cache non userid resources when isOwner calls?
func RequireResourceOwner(
	urlParam string,
	load func(ctx context.Context, reqid, resid uuid.UUID) (any, bool),
	bypassPerms ...models.PermissionAction,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc, err := GetRequestContext(r.Context())
			if err != nil {
				responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
				return
			}

			rid, err := GetUUIDFromURL(r, urlParam)
			if err != nil {
				responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
				return
			}

			if rc.Session.User.Role == models.ROLE_ADMIN {
				next.ServeHTTP(w, r)
				return
			}

			for i := range bypassPerms {
				if perm, exists := rc.Session.User.Permissions[bypassPerms[i]]; exists && perm.Granted {
					if perm.ExpiresAt == nil || time.UnixMilli(*perm.ExpiresAt).After(time.Now().UTC()) {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			cached, ok := load(r.Context(), rc.Session.User.ID, rid)
			if !ok {
				responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrForbidden.Error()))
				return
			}

			if cached != nil {
				rc.Cache[urlParam] = cached
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
