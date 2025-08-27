package distributedsync

import (
	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
)

const (
	errLockingMutexMessage   = "an error occurred while acquiring processes"
	errUnlockingMutexMessage = "an error occurred while releasing processes"
)

type MutexLockingError struct {
	*errutil.CriticalError
}

func NewMutexLockingError(key string) *MutexLockingError {
	return &MutexLockingError{
		CriticalError: errutil.NewCriticalErrorWithMetadata(
			errLockingMutexMessage,
			errutil.NewErrorMetadata().Set("db.operation.parameter.mutex_key", key),
		),
	}
}

type MutexUnlockingError struct {
	*errutil.CriticalError
}

func NewMutexUnlockingError(key string) *MutexUnlockingError {
	return &MutexUnlockingError{
		CriticalError: errutil.NewCriticalErrorWithMetadata(
			errUnlockingMutexMessage,
			errutil.NewErrorMetadata().Set("db.operation.parameter.mutex_key", key),
		),
	}
}
