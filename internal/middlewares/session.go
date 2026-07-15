package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/pkg/errs"
	"github.com/juevigrace/diva-server/pkg/jwt"
)

type SessionCall func(ctx context.Context, sid uuid.UUID) (*models.Session, error)
type UserCall func(ctx context.Context, uid uuid.UUID) (*models.User, error)

func RequiresSession(session SessionCall, user UserCall) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, err := extractSession(r, session, user)
			if err != nil {
				responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrNotAuthorized.Error()))
				return
			}
			rc := NewRequestContext(s)
			ctx := SetRequestContext(r.Context(), rc)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func extractSession(r *http.Request, sessionCall SessionCall, userCall UserCall) (*models.Session, error) {
	claims, err := extractJWTFromHeader(r)
	if err != nil {
		return nil, err
	}

	session, err := sessionCall(r.Context(), claims.SessionID)
	if err != nil {
		return nil, err
	}

	if session.Status != models.SESSION_ACTIVE {
		return nil, errs.ErrNotAuthorized
	}

	user, err := userCall(r.Context(), session.User.ID)
	if err != nil {
		return nil, err
	}

	session.User = *user

	return session, nil
}

func extractJWTFromHeader(r *http.Request) (*jwt.JWTClaims, error) {
	authHeader := r.Header.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer") {
		return nil, errs.ErrHeaderNotValid
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[1] == "" {
		return nil, errs.ErrHeaderNotValid
	}
	tokenString := parts[1]
	claims, err := jwt.ValidateJWT(tokenString)
	if err != nil {
		slog.Warn("jwt validation failed", "error", err)
		return nil, errs.ErrNotAuthorized
	}

	return claims, nil
}
