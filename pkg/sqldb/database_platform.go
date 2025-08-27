package sqldb

const (
	PostgresSQLPlatform Platform = "postgres"
)

type Platform string

func (p Platform) String() string {
	return string(p)
}
