package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/util"
)

type SessionCall func(sessionId uuid.UUID) (*models.Session, error)

func SessionHandler(session SessionCall, handler func(w http.ResponseWriter, r *http.Request, s *models.Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := extractSession(session, r)
		if err != nil {
			log.Printf("Session error: %s\n", err.Error())
			responses.WriteJSON(w, responses.RespondUnauthorized(nil, "you're not authorized to access this endpoint"))
			return
		}
		handler(w, r, s)
	}
}

func CtxSessionHandler(handler func(w http.ResponseWriter, r *http.Request, s *models.Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value("session").(*models.Session)
		if !ok {
			log.Println("Session type is incorrect")
			responses.WriteJSON(w, responses.RespondUnauthorized(nil, "you're not authorized to access this endpoint"))
			return
		}
		handler(w, r, session)
	}
}

func SessionMiddleware(session SessionCall) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, err := extractSession(session, r)
			if err != nil {
				log.Printf("Session error: %s\n", err.Error())
				responses.WriteJSON(w, responses.RespondUnauthorized(nil, "you're not authorized to access this endpoint"))
				return
			}
			ctx := context.WithValue(r.Context(), "session", s)
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

	session, err := sessionCall(claims.SessionID)
	if err != nil {
		return nil, err
	}

	return session, err
}

func extractJWTFromHeader(r *http.Request) (*util.JWTClaims, error) {
	header := strings.Join(r.Header["Authorization"], "")

	if !strings.HasPrefix(header, "Bearer") {
		return nil, errors.New("permission denied")
	}

	tokenString := strings.Split(header, " ")[1]
	claims, err := util.ValidateJWT(tokenString)
	if err != nil {
		return nil, errors.New("permission denied")
	}

	return claims, nil
}
