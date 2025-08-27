package sqldb

const (
	noMigrationsRunCount = 0
)

type Migrator interface {
	Up() (int, error)
	Down() (int, error)
}
