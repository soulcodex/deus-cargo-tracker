package utils

import (
	"sync"
	"time"
)

// DateTimeProvider represents anything that can supply the current time.
type DateTimeProvider interface {
	// Now returns the current time.
	Now() time.Time
}

// SystemTimeProvider uses the real system clock.
type SystemTimeProvider struct{}

// NewSystemTimeProvider constructs a SystemTimeProvider.
func NewSystemTimeProvider() *SystemTimeProvider {
	return &SystemTimeProvider{}
}

// Now returns the current local system time.
func (stp *SystemTimeProvider) Now() time.Time {
	return time.Now()
}

// FixedTimeProvider returns a single, deterministic timestamp.
// On the very first Now() call it either uses the supplied time
// (via NewFixedTimeProviderAt) or captures time.Now().Round(time.Millisecond).
type FixedTimeProvider struct {
	now  time.Time
	once sync.Once
}

// NewFixedTimeProvider constructs a provider that freezes at the
// first Now() invocation (rounded to the nearest millisecond).
func NewFixedTimeProvider() *FixedTimeProvider {
	return new(FixedTimeProvider)
}

// NewFixedTimeProviderAt constructs a provider frozen at the
// given time.Time from the get-go.
func NewFixedTimeProviderAt(t time.Time) *FixedTimeProvider {
	return &FixedTimeProvider{now: t, once: sync.Once{}}
}

// Now returns the frozen time. On the first call (and only once),
// if the provider was created with NewFixedTimeProvider(), it captures
// time.Now().Round(time.Millisecond).
func (ftp *FixedTimeProvider) Now() time.Time {
	ftp.once.Do(func() {
		if ftp.now.IsZero() {
			ftp.now = time.Now().Round(time.Millisecond)
		}
	})
	return ftp.now
}
