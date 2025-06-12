package ratelimiter

import (
	"fmt"
	"math"
	"time"
)

type Reservation struct {
	Ok        bool
	lim       *Limiter
	tokens    int
	timeToAct time.Time
	// This is the Limit at reservation time, it can change later.
	limit Limit
}

func (limter *Limiter) Reserve(t time.Time, n int) Reservation {
	limter.mu.Lock()
	defer limter.mu.Unlock()

	fmt.Println(limter) // debuging

	if limter.limit == math.MaxFloat64 {
		return Reservation{
			Ok:        true,
			lim:       limter,
			tokens:    n, // number of request in t frame of time
			timeToAct: t, // time.Now()
		}
	}

	tokens := limter.advance(t)

	tokens -= float64(n)

	var waitDuration time.Duration
	if tokens < 0 {
		waitDuration = limter.limit.DurationFromTokens(-tokens)
	}

	ok := n <= limter.burst && waitDuration <= 0

	r := Reservation{
		Ok:    ok,
		lim:   limter,
		limit: limter.limit,
	}

	if ok {
		r.tokens = n
		r.timeToAct = t.Add(waitDuration)

		// Update state
		limter.last = t
		limter.tokens = tokens
		limter.lastEvent = r.timeToAct
	}

	return r

}

func (limiter *Limiter) advance(t time.Time) float64 {
	last := limiter.last
	if t.Before(last) { //
		last = t
	}

	passed := t.Sub(last)
	newtokens := limiter.limit.TokensFromDuration(passed)
	newtokens += limiter.tokens
	burst := float64(limiter.burst)

	if newtokens > burst {
		newtokens = burst
	}
	return newtokens
}
