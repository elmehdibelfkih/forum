package middleware

import (
	forumerror "forum/internal/error"
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
		reservation := limiter.Reserve(time.Now(), 1)
		if !reservation.Ok {
			// FIXME: FIX THE COMMENT FLAG
			forumerror.TooManyRequests(w, r, "Request")
			return
		}

		next.ServeHTTP(w, r)
	})
}
