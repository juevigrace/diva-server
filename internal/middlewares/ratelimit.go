package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/juevigrace/diva-server/internal/models/responses"
)

type visitor struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()
		v, ok := rl.visitors[ip]
		now := time.Now()
		if !ok || now.After(v.resetAt) {
			v = &visitor{count: 0, resetAt: now.Add(rl.window)}
			rl.visitors[ip] = v
		}
		v.count++
		rl.mu.Unlock()

		if v.count > rl.limit {
			responses.WriteJSON(w, responses.RespondTooManyRequests(nil, "rate limit exceeded"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.After(v.resetAt) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
