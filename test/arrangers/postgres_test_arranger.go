package testarrangers

import (
	"context"
	"fmt"
	"sync"

	sq "github.com/Masterminds/squirrel"

	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
)

type PostgresSQLArranger struct {
	schema string
	pool   *postgres.ConnectionPool
	ignore map[string]struct{}
}

func NewPostgresSQLArranger(
	schema string,
	pool *postgres.ConnectionPool,
) *PostgresSQLArranger {
	return &PostgresSQLArranger{
		schema: schema,
		pool:   pool,
		ignore: map[string]struct{}{"migrations": {}},
	}
}

func (pa *PostgresSQLArranger) MustArrange(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	pa.arrangeDB(ctx, wg)
}

func (pa *PostgresSQLArranger) arrangeDB(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	rows, err := sq.Select("table_name").
		From(pa.schema + "." + "information_schema.tables").
		Where(sq.Eq{"table_schema": pa.schema}).
		PlaceholderFormat(sq.Dollar).
		RunWith(pa.pool.Reader()).
		QueryContext(ctx)

	if err != nil {
		panic(err)
	}

	defer func() {
		sqldb.CloseRows(rows)
		_ = rows.Close()
	}()

	if rowsErr := rows.Err(); rowsErr != nil {
		panic(rowsErr)
	}

	var tableName string

	for rows.Next() {
		if scanErr := rows.Scan(&tableName); scanErr != nil {
			panic(scanErr)
		}

		if _, ok := pa.ignore[tableName]; ok {
			continue
		}

		truncateSQL := fmt.Sprintf("TRUNCATE TABLE %s", pa.schema+"."+tableName)
		if _, execErr := pa.pool.Writer().ExecContext(ctx, truncateSQL); execErr != nil {
			panic(execErr)
		}
	}
}
