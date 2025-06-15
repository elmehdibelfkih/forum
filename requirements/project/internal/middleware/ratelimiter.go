package middleware

import (
	forumerror "forum/internal/error"
	"forum/internal/ratelimiter"
	"net/http"
	"sync"
)

func RateLimiterMiddleware(next http.Handler, fillRate float64, burst uint64) http.Handler {
	ipLimiterMap := make(map[string]*ratelimiter.TokenBucketLimiter)
	var mu sync.Mutex
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the ip of the request
		ip := ratelimiter.GetIP(r)

		// create a limeter if the corresponding ip doest have one
		mu.Lock()
		limiter, exists := ipLimiterMap[ip]
		if !exists {
			limiter = ratelimiter.NewTokenBucketLimiter(fillRate, burst)
			ipLimiterMap[ip] = limiter
		}
		mu.Unlock()

		// error if we reached the limit

		if !limiter.Allow() {
			forumerror.TooManyRequests(w, r, "Request")
			return
		}

		next.ServeHTTP(w, r)
	})
}
