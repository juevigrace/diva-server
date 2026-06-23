package middlewares

import (
	"net"
	"net/http"

	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/pkg/errs"
)

func RequireLocalhost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			responses.WriteJSON(w, responses.RespondForbbiden(nil, errs.ErrForbidden.Error()))
			return
		}
		if host != "127.0.0.1" && host != "::1" {
			responses.WriteJSON(w, responses.RespondForbbiden(nil, errs.ErrForbidden.Error()))
			return
		}
		next.ServeHTTP(w, r)
	})
}
