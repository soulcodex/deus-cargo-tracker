package postgres

const (
	defaultUser     = "postgres"
	defaultPassword = "postgres"
	defaultHost     = "localhost"
	defaultSchema   = "public"
	defaultPort     = 5432

	defaultMaxConnections  = 20
	defaultConnectionsIdle = 50
	defaultMaxLifetime     = 3
)

type Credentials struct {
	User     string
	Password string
	Host     string
	Port     uint16
	Database string
	Schema   string
}

func NewCredentials(
	user string,
	password string,
	host string,
	port uint16,
	schema,
	database string,
) Credentials {
	return Credentials{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
		Schema:   schema,
	}
}

func NewDefaultCredentials(database string) Credentials {
	return Credentials{
		User:     defaultUser,
		Password: defaultPassword,
		Host:     defaultHost,
		Port:     defaultPort,
		Database: database,
		Schema:   defaultSchema,
	}
}

type ClientOptionsFunc func(co *ClientOptions)

type ClientOptions struct {
	Credentials    Credentials
	Traceable      bool
	MaxConnections int
	ConnIdle       int
	MaxLifetime    int
	SSLMode        string
}

func NewDefaultClientOptions(credentials Credentials) *ClientOptions {
	return &ClientOptions{
		Credentials:    credentials,
		Traceable:      false,
		MaxConnections: defaultMaxConnections,
		ConnIdle:       defaultConnectionsIdle,
		MaxLifetime:    defaultMaxLifetime,
		SSLMode:        VerifyFullMode.String(),
	}
}

func (co *ClientOptions) apply(options ...ClientOptionsFunc) {
	for _, opt := range options {
		opt(co)
	}
}

func WithTraces() ClientOptionsFunc {
	return func(co *ClientOptions) {
		co.Traceable = true
	}
}

func WithMaxConnections(maxConnections int) ClientOptionsFunc {
	return func(co *ClientOptions) {
		co.MaxConnections = maxConnections
	}
}

func WithConnIdle(connIdle int) ClientOptionsFunc {
	return func(co *ClientOptions) {
		co.ConnIdle = connIdle
	}
}

func WithMaxLifetime(maxLifetime int) ClientOptionsFunc {
	return func(co *ClientOptions) {
		co.MaxLifetime = maxLifetime
	}
}

func WithSSLMode(mode SSLMode) ClientOptionsFunc {
	return func(co *ClientOptions) {
		co.SSLMode = mode.String()
	}
}
