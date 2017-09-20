package tally

import (
	"sync"
	"time"
)

// Span represents the count for the given duration
type Span struct {
	Count    int64
	Duration time.Duration
}

type track struct {
	count    int64
	duration time.Duration
	last     time.Time
}

// Tally holds the counts for the specific durations
type Tally struct {
	sync.RWMutex
	start time.Time
	spans []track
}

// NewTally returns a tally to track counters over durations
func NewTally(durations ...time.Duration) *Tally {
	now := time.Now()
	t := &Tally{
		start: now,
		spans: make([]track, len(durations)),
	}
	for i, d := range durations {
		t.spans[i].duration = d
		t.spans[i].last = now
	}
	return t
}

// Reset resets the counters for each duration
func (t *Tally) Reset(amt int64) {
	now := time.Now()
	t.Lock()
	for i := range t.spans {
		t.spans[i].count = 0
		t.spans[i].last = now
	}
	t.start = now
	t.Unlock()
}

// Increment increments the count for each
func (t *Tally) Increment(amt int64) {
	now := time.Now()
	t.Lock()
	for i, s := range t.spans {
		if s.duration < now.Sub(s.last) {
			t.spans[i].count += amt
		} else {
			t.spans[i].count = amt
		}
	}
	t.Unlock()
}

func (t *Tally) Tallies() []Span {
	spans := make([]Span, len(t.spans))
	t.RLock()
	for i, s := range t.spans {
		spans[i] = Span{Duration: s.duration, Count: s.count}
	}
	t.RUnlock()
	return spans
}
