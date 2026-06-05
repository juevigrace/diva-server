package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/errs"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/util"
)

type SessionCall func(ctx context.Context, sessionId uuid.UUID) (*models.Session, error)

func RequiresSession(session SessionCall) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, err := extractSession(session, r)
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

func extractSession(sessionCall SessionCall, r *http.Request) (*models.Session, error) {
	claims, err := extractJWTFromHeader(r)
	if err != nil {
		return nil, err
	}

	session, err := sessionCall(r.Context(), claims.SessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func extractJWTFromHeader(r *http.Request) (*util.JWTClaims, error) {
	authHeader := r.Header.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer") {
		return nil, errs.ErrHeaderNotValid
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[1] == "" {
		return nil, errs.ErrHeaderNotValid
	}
	tokenString := parts[1]
	claims, err := util.ValidateJWT(tokenString)
	if err != nil {
		slog.Warn("jwt validation failed", "error", err)
		return nil, errs.ErrNotAuthorized
	}

	return claims, nil
}
