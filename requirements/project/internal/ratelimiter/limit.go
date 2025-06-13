package ratelimiter

import (
	"math"
	"time"
)

type Limit float64

func (l Limit) TokensFromDuration(d time.Duration) float64 {
	if l <= 0 {
		return 0
	}
	return d.Seconds() * float64(l)
}


func (l Limit) DurationFromTokens(token float64) time.Duration {

	if l <= 0 {
		return time.Duration(math.MaxInt64)
	}

	duration := (token / float64(l)) * float64(time.Second)

	if duration > float64(math.MaxInt64) {
		return time.Duration(math.MaxInt64)
	}

	return time.Duration(duration)
}
