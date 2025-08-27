package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

// ZerologLogger is a type-aliased zerolog.Logger.
type ZerologLogger = zerolog.Logger

// Option is a function that configures the logger.
type Option func(*options)

// WithLogLevel sets the log level.
func WithLogLevel(level string) Option {
	return func(o *options) {
		lvl, err := zerolog.ParseLevel(level)
		if err != nil {
			panic(err)
		}

		o.level = lvl
	}
}

// WithAppVersion sets the application version.
func WithAppVersion(appVersion string) Option {
	return func(o *options) {
		o.appVersion = appVersion
	}
}

type options struct {
	level      zerolog.Level
	appVersion string
}

func defaultOptions() *options {
	return &options{
		level:      zerolog.InfoLevel,
		appVersion: "0.0.0",
	}
}

// NewZerologLogger creates a new ZerologLogger instance.
func NewZerologLogger(ctx context.Context, appName string, opts ...Option) ZerologLogger {
	zeroLoggerOpts := defaultOptions()
	for _, opt := range opts {
		opt(zeroLoggerOpts)
	}

	zeroLogger := zerolog.New(os.Stdout).
		With().
		Ctx(ctx).
		Timestamp().
		Str("app_name", appName).
		Str("app_version", zeroLoggerOpts.appVersion).
		Logger()

	return zeroLogger.Level(zeroLoggerOpts.level)
}

// MustSetGlobalLevel sets the global override for log level. If this
// function is not called, the global level is InfoLevel.
func MustSetGlobalLevel(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(err)
	}

	zerolog.SetGlobalLevel(lvl)
}
