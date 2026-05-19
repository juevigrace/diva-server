package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/util"
)

type contextKey string

const sessionContextKey contextKey = "session"

type SessionCall func(ctx context.Context, sessionId uuid.UUID) (*models.Session, error)

func SessionMiddleware(session SessionCall) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, err := extractSession(session, r)
			if err != nil {
				responses.WriteJSON(w, responses.RespondUnauthorized(nil, models.ErrNotAuthorized.Error()))
				return
			}
			ctx := context.WithValue(r.Context(), sessionContextKey, s)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func GetSessionFromContext(ctx context.Context) (*models.Session, bool) {
	session, ok := ctx.Value(sessionContextKey).(*models.Session)
	return session, ok
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
		return nil, models.ErrHeaderNotValid
	}

	tokenString := strings.Split(authHeader, " ")[1]
	claims, err := util.ValidateJWT(tokenString)
	if err != nil {
		log.Printf("jwt invalid: %s\n", err.Error())
		return nil, models.ErrNotAuthorized
	}

	return claims, nil
}
