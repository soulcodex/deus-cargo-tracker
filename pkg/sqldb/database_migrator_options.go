package sqldb

const (
	defaultMigratorMigrationPath       = "migrations"
	defaultMigratorMigrationsTableName = "migrations"
	defaultMigratorPlatform            = PostgresSQLPlatform
	defaultMigratorSchema              = "public"
	defaultDatabaseName                = "unknown"
)

type MigratorOptFunc func(*MigratorOptions)

type MigratorOptions struct {
	MigrationsPath      string
	MigrationsTableName string
	Platform            Platform
	DatabaseName        string
	Schema              string
}

func NewDatabaseMigratorOptions(opts ...MigratorOptFunc) *MigratorOptions {
	options := &MigratorOptions{
		MigrationsPath:      defaultMigratorMigrationPath,
		MigrationsTableName: defaultMigratorMigrationsTableName,
		Platform:            defaultMigratorPlatform,
		Schema:              defaultMigratorSchema,
		DatabaseName:        defaultDatabaseName,
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithMigrationsPath(path string) func(*MigratorOptions) {
	return func(opts *MigratorOptions) {
		opts.MigrationsPath = path
	}
}

func WithMigrationsTableName(name string) func(*MigratorOptions) {
	return func(opts *MigratorOptions) {
		opts.MigrationsTableName = name
	}
}

func WithPlatform(platform Platform) func(*MigratorOptions) {
	return func(opts *MigratorOptions) {
		opts.Platform = platform
	}
}

func WithDatabaseName(name string) func(*MigratorOptions) {
	return func(opts *MigratorOptions) {
		opts.DatabaseName = name
	}
}

func WithSchema(schema string) func(*MigratorOptions) {
	return func(opts *MigratorOptions) {
		opts.Schema = schema
	}
}
