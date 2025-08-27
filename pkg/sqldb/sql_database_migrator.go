package sqldb

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

type SQLDatabaseMigrator struct {
	client          *sql.DB
	migrationSet    *migrate.MigrationSet
	migrationSource migrate.MigrationSource
	platform        Platform
	databaseName    string
}

func NewSQLDatabaseMigrator(
	client *sql.DB,
	opts ...MigratorOptFunc,
) *SQLDatabaseMigrator {
	options := NewDatabaseMigratorOptions(opts...)
	migrationsSource := &migrate.FileMigrationSource{Dir: options.MigrationsPath}
	migrationsSet := &migrate.MigrationSet{
		SchemaName:         "",
		IgnoreUnknown:      false,
		DisableCreateTable: false,
		TableName:          options.MigrationsTableName,
	}

	if options.Schema != "" {
		migrationsSet.SchemaName = options.Schema
	}

	return &SQLDatabaseMigrator{
		client:          client,
		migrationSet:    migrationsSet,
		migrationSource: migrationsSource,
		platform:        options.Platform,
		databaseName:    options.DatabaseName,
	}
}

func (sdm *SQLDatabaseMigrator) Up() (int, error) {
	appliedMigrations, err := sdm.migrationSet.Exec(sdm.client, sdm.platform.String(), sdm.migrationSource, migrate.Up)
	if err != nil {
		return noMigrationsRunCount, NewDatabaseMigrationError(
			sdm.platform,
			sdm.migrationSet.SchemaName,
			sdm.databaseName,
			"up",
		)
	}

	return appliedMigrations, nil
}

func (sdm *SQLDatabaseMigrator) Down() (int, error) {
	appliedMigrations, err := sdm.migrationSet.Exec(sdm.client, sdm.platform.String(), sdm.migrationSource, migrate.Down)
	if err != nil {
		return noMigrationsRunCount, NewDatabaseMigrationError(
			sdm.platform,
			sdm.migrationSet.SchemaName,
			sdm.databaseName,
			"down",
		)
	}

	return appliedMigrations, nil
}
