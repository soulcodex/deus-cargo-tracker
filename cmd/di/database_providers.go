package di

import (
	"context"

	"github.com/soulcodex/deus-cargo-tracker/configs"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
)

func initPostgresDBPool(_ context.Context, cfg *configs.Config) *postgres.ConnectionPool {
	dbCredentials := postgres.NewCredentials(
		cfg.PostgresUser,
		cfg.PostgresPass,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresSchema,
		cfg.PostgresDB,
	)

	dbSSLMode, err := postgres.SSLModeFromString(cfg.PostgresSSL)
	if err != nil {
		panic(err)
	}

	dbWriter, err := postgres.NewWriter(dbCredentials, postgres.WithSSLMode(dbSSLMode))
	if err != nil {
		panic(err)
	}

	dbPool, err := postgres.WithWriterOnly(dbWriter)
	if err != nil {
		panic(err)
	}

	return dbPool
}

func initSQLMigrator(_ context.Context, cfg *configs.Config, dbPool sqldb.ConnectionPool) sqldb.Migrator {
	return sqldb.NewSQLDatabaseMigrator(
		dbPool.Writer(),
		sqldb.WithPlatform(sqldb.PostgresSQLPlatform),
		sqldb.WithSchema(cfg.PostgresSchema),
		sqldb.WithDatabaseName(cfg.PostgresDB),
		sqldb.WithMigrationsPath(cfg.MigrationsPath),
		sqldb.WithMigrationsTableName(cfg.MigrationsTable),
	)
}
