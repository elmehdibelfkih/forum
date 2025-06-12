package middleware

import (
	"forum/internal/ratelimiter"
	"net/http"
	"time"
)

func RateLimiterMiddleware(next http.Handler, limit ratelimiter.Limit, brust int) http.Handler {
	ipLimiterMap := make(map[string]*ratelimiter.Limiter)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the ip of the request
		ip := ratelimiter.GetIP(r)

		// create a limeter if the corresponding ip doest have one
		limiter, exists := ipLimiterMap[ip]
		if !exists {
			limiter = ratelimiter.NewLimiter(limit, brust)
			ipLimiterMap[ip] = limiter
		}

		// error if we reached the limit
		reservation := limiter.Reserve(time.Now(), 2)
		if !reservation.Ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			// json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})

			return
		}

		next.ServeHTTP(w, r)
	})
}
