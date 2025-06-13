package ratelimiter

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	limit  Limit
	burst  int
	tokens float64
	// last is the last time the limiter's tokens field was updated
	last time.Time
	// lastEvent is the latest time of a rate-limited event (past or future)
	lastEvent time.Time
}

func NewLimiter(r Limit, b int) *Limiter {
	return &Limiter{
		limit:  r,
		burst:  b,
		tokens: float64(b),
	}
}

func GetIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Error parssing IP :%v", err)
		return ""
	}
	return host
}
