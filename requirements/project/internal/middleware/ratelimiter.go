package middleware

import (
	"forum/internal/ratelimiter"
	"net/http"
	"sync"
)

func RateLimiterMiddleware(next http.Handler, limit float64, burst uint64) http.Handler {
	ipLimiterMap := make(map[string]*ratelimiter.TokenBucketLimiter)
	var mu sync.Mutex

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ratelimiter.GetIP(r)

		limiter, exists := ipLimiterMap[ip]
		mu.Lock()
		if !exists {
			limiter = ratelimiter.NewTokenBucketLimiter(limit, burst)
			ipLimiterMap[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			// w.WriteHeader(http.StatusTooManyRequests)
			http.ServeFile(w, r, "./templates/rate_limiting.html")
			return
		}

		next.ServeHTTP(w, r)
	})
}
