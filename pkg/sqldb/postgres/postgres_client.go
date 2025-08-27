package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
)

const driver = "postgres"

func NewReader(
	credentials Credentials,
	opts ...ClientOptionsFunc,
) (*sql.DB, error) {
	options := NewDefaultClientOptions(credentials)
	options.apply(opts...)

	return buildSQLClientFromOptions(options)
}

func NewWriter(
	credentials Credentials,
	opts ...ClientOptionsFunc,
) (*sql.DB, error) {
	options := NewDefaultClientOptions(credentials)
	options.apply(opts...)

	return buildSQLClientFromOptions(options)
}

func buildSQLClientFromOptions(options *ClientOptions) (*sql.DB, error) {
	driverName, buildErr := buildDriverNameFromOptions(options)
	if buildErr != nil {
		return nil, buildErr
	}

	srvAddress, addrErr := buildPostgresConnectionStringFromOptions(options)
	if addrErr != nil {
		return nil, addrErr
	}

	client, clientErr := sql.Open(driverName, srvAddress)
	if clientErr != nil {
		return nil, errutil.NewCriticalError("postgres connection error").Wrap(clientErr)
	}

	if pingErr := client.Ping(); pingErr != nil {
		return nil, errutil.NewCriticalError("postgres ping error").Wrap(pingErr)
	}

	client.SetConnMaxLifetime(time.Duration(options.MaxLifetime) * time.Minute)
	client.SetMaxOpenConns(options.MaxConnections)
	client.SetMaxIdleConns(options.ConnIdle)

	return client, nil
}

func buildDriverNameFromOptions(options *ClientOptions) (string, error) {
	var driverName = driver

	if options.Traceable {
		tracedDriverName, driverRegErr := otelsql.Register(
			driver,
			otelsql.AllowRoot(),
			otelsql.TraceQueryWithoutArgs(),
			otelsql.TraceRowsAffected(),
			otelsql.WithDatabaseName(options.Credentials.Database),
			otelsql.WithSystem(semconv.DBSystemPostgreSQL),
		)

		if driverRegErr != nil {
			return "", errutil.NewCriticalError("postgres OTEL driver register error").Wrap(driverRegErr)
		}
		driverName = tracedDriverName
	}

	return driverName, nil
}

func buildPostgresConnectionStringFromOptions(options *ClientOptions) (string, error) {
	rawAddress := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		driver,
		options.Credentials.User,
		options.Credentials.Password,
		options.Credentials.Host,
		options.Credentials.Port,
		options.Credentials.Database,
		options.SSLMode,
	)

	srvAddress, err := pq.ParseURL(rawAddress)
	if err != nil {
		return "", errutil.NewCriticalError("error parsing postgres URL").Wrap(err)
	}

	return srvAddress, nil
}
