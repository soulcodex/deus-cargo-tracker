package sqldb

import (
	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
)

const (
	poolConfigProvidedErrorMessage = "invalid pool config provided"
)

type PoolConfigProvidedError struct {
	*errutil.CriticalError
}

func NewPoolConfigProvidedError(driverName string) *PoolConfigProvidedError {
	return &PoolConfigProvidedError{
		CriticalError: errutil.NewCriticalErrorWithMetadata(
			poolConfigProvidedErrorMessage,
			errutil.NewErrorMetadata().
				Set("db.operation.parameter.driver_name", driverName),
		),
	}
}
