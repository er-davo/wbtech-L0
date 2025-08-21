package retry

import (
	"math"
	"math/rand/v2"
	"time"
)

type Backoff interface {
	// Next returns the duration to wait before the next attempt.
	// attempt 0 based.
	Next(attempt int) time.Duration
}

type FixedBackoff struct {
	Interval time.Duration
	Jitter   float64
}

func (f FixedBackoff) Next(attempt int) time.Duration {
	return addJitter(time.Duration(f.Interval), f.Jitter)
}

type LinearBackoff struct {
	Base   time.Duration
	Step   time.Duration
	Max    time.Duration
	Jitter float64
}

func (l LinearBackoff) Next(attempt int) time.Duration {
	d := l.Base + time.Duration(attempt)*l.Step
	if l.Max > 0 && d > l.Max {
		return l.Max
	}
	return addJitter(time.Duration(d), l.Jitter)
}

type ExponentialBackoff struct {
	Base   time.Duration
	Factor float64
	Max    time.Duration
	Jitter float64
}

func (e ExponentialBackoff) Next(attempt int) time.Duration {
	d := float64(e.Base) * math.Pow(e.Factor, float64(attempt))
	if e.Max > 0 && d > float64(e.Max) {
		return e.Max
	}
	return addJitter(time.Duration(d), e.Jitter)
}

func addJitter(d time.Duration, jitter float64) time.Duration {
	if jitter <= 0 || jitter >= 1 {
		return d
	}
	delta := (rand.Float64()*2 - 1) * jitter
	return time.Duration(float64(d) * (1 + delta))
}
