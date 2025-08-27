package httpserver

const (
	defaultHost                = "0.0.0.0"
	defaultPort                = 8080
	defaultWriteTimeoutSeconds = 30
	defaultReadTimeoutSeconds  = 30
)

type RouterConfigFunc func(*RouterConfig)

type RouterConfig struct {
	Host         string
	Port         int
	WriteTimeout int
	ReadTimeout  int
	Middlewares  []Middleware
}

func NewRouterConfig(opts ...RouterConfigFunc) RouterConfig {
	rc := RouterConfig{
		Host:         defaultHost,
		Port:         defaultPort,
		WriteTimeout: defaultWriteTimeoutSeconds,
		ReadTimeout:  defaultReadTimeoutSeconds,
		Middlewares:  make([]Middleware, 0),
	}

	for _, opt := range opts {
		opt(&rc)
	}

	return rc
}

func NewDefaultRouterConfig() []RouterConfigFunc {
	return []RouterConfigFunc{
		WithHost(defaultHost),
		WithPort(defaultPort),
		WithWriteTimeoutSeconds(defaultWriteTimeoutSeconds),
		WithReadTimeoutSeconds(defaultReadTimeoutSeconds),
		WithCORSMiddleware(),
	}
}

func WithHost(host string) RouterConfigFunc {
	return func(rc *RouterConfig) {
		rc.Host = host
	}
}

func WithPort(port int) RouterConfigFunc {
	return func(rc *RouterConfig) {
		rc.Port = port
	}
}

func WithWriteTimeoutSeconds(seconds int) RouterConfigFunc {
	return func(rc *RouterConfig) {
		rc.WriteTimeout = seconds
	}
}

func WithReadTimeoutSeconds(seconds int) RouterConfigFunc {
	return func(rc *RouterConfig) {
		rc.ReadTimeout = seconds
	}
}

func WithMiddleware(middleware Middleware) RouterConfigFunc {
	return func(rc *RouterConfig) {
		rc.Middlewares = append(rc.Middlewares, middleware)
	}
}

func WithCORSMiddleware() RouterConfigFunc {
	return func(rc *RouterConfig) {
		rc.Middlewares = append(rc.Middlewares, enableCorsMiddleware)
	}
}
