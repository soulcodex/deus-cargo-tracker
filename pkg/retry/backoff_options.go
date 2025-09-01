package retry

import (
	"time"

	"github.com/rs/zerolog"

	"github.com/soulcodex/deus-cargo-tracker/pkg/logger"
)

const (
	defaultMaxRetries          = 5
	defaultInitialIntervalMs   = 100
	defaultMaxInterval         = 3 * time.Second
	defaultMultiplier          = 2.0
	defaultMaxElapsedTime      = 10 * time.Second
	defaultRandomizationFactor = 0.0
)

// Options configures Backoff behavior.
type Options struct {
	MaxRetries          int
	InitialInterval     time.Duration
	MaxInterval         time.Duration
	Multiplier          float64
	MaxElapsedTime      time.Duration
	RandomizationFactor float64
	OnRetryHook         OnRetryHook
	Logger              *logger.ZerologLogger
}

// OptionFunc applies a configuration to Options.
type OptionFunc func(*Options)

// NewOptions returns Options populated with defaults,
// then applies any provided OptionFunc.
func NewOptions(opts ...OptionFunc) Options {
	defaultLogger := zerolog.New(zerolog.Nop())
	defaultRetryHook := func(_ int, _ time.Duration, _ error) bool {
		return false
	}
	defaultOptions := Options{
		MaxRetries:          defaultMaxRetries,
		InitialInterval:     defaultInitialIntervalMs * time.Millisecond,
		MaxInterval:         defaultMaxInterval,
		Multiplier:          defaultMultiplier,
		MaxElapsedTime:      defaultMaxElapsedTime,
		RandomizationFactor: defaultRandomizationFactor,
		Logger:              &defaultLogger,
		OnRetryHook:         defaultRetryHook,
	}
	for _, fn := range opts {
		fn(&defaultOptions)
	}
	return defaultOptions
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(n int) OptionFunc {
	return func(o *Options) { o.MaxRetries = n }
}

// WithInitialInterval sets the initial backoff interval.
func WithInitialInterval(d time.Duration) OptionFunc {
	return func(o *Options) { o.InitialInterval = d }
}

// WithMaxInterval sets the maximum backoff interval.
func WithMaxInterval(d time.Duration) OptionFunc {
	return func(o *Options) { o.MaxInterval = d }
}

// WithMultiplier sets the backoff multiplier.
func WithMultiplier(m float64) OptionFunc {
	return func(o *Options) { o.Multiplier = m }
}

// WithMaxElapsedTime sets the total timeout across all retries.
func WithMaxElapsedTime(d time.Duration) OptionFunc {
	return func(o *Options) { o.MaxElapsedTime = d }
}

// WithRandomizationFactor sets jitter factor for backoff.
func WithRandomizationFactor(f float64) OptionFunc {
	return func(o *Options) { o.RandomizationFactor = f }
}

// WithOnRetryHook sets a hook called after each failed attempt.
func WithOnRetryHook(hook OnRetryHook) OptionFunc {
	return func(o *Options) { o.OnRetryHook = hook }
}

// WithLogger sets a ZerologLogger for retry events.
func WithLogger(log logger.ZerologLogger) OptionFunc {
	return func(o *Options) { o.Logger = &log }
}
