package sqldb

import (
	"fmt"

	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

const (
	databaseMigrationErrorMessage = "an error occurred during migration"
)

type DatabaseMigrationError struct {
	*errutil.CriticalError
}

func NewDatabaseMigrationError(platform Platform, schema, databaseName, direction string) *PoolConfigProvidedError {
	return &PoolConfigProvidedError{
		CriticalError: errutil.NewCriticalErrorWithMetadata(
			databaseMigrationErrorMessage,
			errutil.NewErrorMetadata().
				Set("db.namespace", fmt.Sprintf("%s.%s", schema, databaseName)).
				Set("db.operation.parameter.direction", direction).
				Set("db.system.name", platform.String()),
		),
	}
}
