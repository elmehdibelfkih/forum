package middleware

import (
	"forum/internal/ratelimiter"
	"net/http"
	"time"
)

func RateLimiterMiddleware(next http.Handler, limit ratelimiter.Limit, brust int) http.Handler {
	ipLimiterMap := make(map[string]*ratelimiter.Limiter)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ratelimiter.GetIP(r)

		limiter, exists := ipLimiterMap[ip]
		if !exists {
			limiter = ratelimiter.NewLimiter(limit, brust)
			ipLimiterMap[ip] = limiter
		}

		reservation := limiter.Reserve(time.Now(), 1)
		if !reservation.Ok {
			// w.WriteHeader(http.StatusTooManyRequests)
			http.ServeFile(w, r, "./templates/rate_limiting.html")
			return
		}

		next.ServeHTTP(w, r)
	})
}
