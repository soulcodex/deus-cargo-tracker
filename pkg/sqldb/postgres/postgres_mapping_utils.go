package postgres

import (
	"database/sql"
)

type DecodeFunc[T any] func(rows *sql.Rows) (T, error)
type EncodeFunc[T any] func(T) ([]any, error)
