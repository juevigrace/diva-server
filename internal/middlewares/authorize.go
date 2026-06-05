package middlewares

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/errs"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

func Require(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.With(middlewares...).Use(Authorize)
}

func WithRequire(r chi.Router, middlewares ...func(http.Handler) http.Handler) chi.Router {
	return r.With(middlewares...).With(Authorize)
}

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc, err := GetRequestContext(r.Context())
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		for _, req := range rc.Roles {
			if !req.Satisfied {
				responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrForbidden.Error()))
				return
			}
		}
		for _, req := range rc.Ownerships {
			if !req.Satisfied {
				responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrForbidden.Error()))
				return
			}
		}
		for _, req := range rc.Permissions {
			if !req.Satisfied {
				responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrForbidden.Error()))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

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

			satisfied := slices.Contains(roles, rc.Session.User.Role)
			rc.Roles = append(rc.Roles, &RoleRequirement{
				Satisfied: satisfied,
				Roles:     roles,
			})
			next.ServeHTTP(w, r)
		})
	}
}

// TODO: should i cache non userid resources when isOwner calls?
func RequireResourceOwner(urlParam string, load func(ctx context.Context, reqid, resid uuid.UUID) (any, bool)) func(http.Handler) http.Handler {
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

			if perm, exists := rc.Session.User.Permissions[models.PERMISSION_OWNERSHIP_BYPASS]; exists && perm.Granted {
				if perm.ExpiresAt == nil || time.UnixMilli(*perm.ExpiresAt).After(time.Now().UTC()) {
					rc.Ownerships = append(rc.Ownerships, &OwnershipRequirement{
						Satisfied: true,
					})
					next.ServeHTTP(w, r)
					return
				}
			}

			cached, ok := load(r.Context(), rc.Session.User.ID, rid)
			if !ok {
				responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrForbidden.Error()))
				return
			}

			entry := make(map[string]any, 1)
			entry[urlParam] = cached
			rc.Ownerships = append(rc.Ownerships, &OwnershipRequirement{
				Satisfied: true,
				Cache:     entry,
			})
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

			for _, action := range actions {
				perm, exists := rc.Session.User.Permissions[action]
				if !exists || !perm.Granted {
					responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrPermissionDenied.Error()))
					return
				}
				if perm.ExpiresAt != nil && time.UnixMilli(*perm.ExpiresAt).Before(time.Now().UTC()) {
					responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrPermissionDenied.Error()))
					return
				}
			}

			rc.Permissions = append(rc.Permissions, &PermissionRequirement{
				Satisfied:         true,
				PermissionActions: actions,
			})
			next.ServeHTTP(w, r)
		})
	}
}
