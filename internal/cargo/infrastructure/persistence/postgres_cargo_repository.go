package cargopersistence

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
)

var (
	_ cargodomain.CargoRepository = (*PostgresCargoRepository)(nil)

	ErrFetchingCargoRows = errutil.NewError("error fetching cargo rows")
	ErrRunningQuery      = errutil.NewError("error running cargo query")
	ErrSavingCargo       = errutil.NewError("error saving cargo to the database")
)

type PostgresCargoRepository struct {
	tableName    string
	pool         sqldb.ConnectionPool
	decoder      postgres.DecodeFunc[*cargodomain.Cargo]
	encoder      postgres.EncodeFunc[*cargodomain.Cargo]
	errorHandler *postgres.ErrorHandler
	fields       []string
}

func NewPostgresCargoRepository(schema string, pool sqldb.ConnectionPool) *PostgresCargoRepository {
	errorHandlers := postgres.ErrorHandlers{
		postgres.UniqueViolationErrorCode: uniqueViolationPostgresCargoRepoErrorHandler(),
	}

	return &PostgresCargoRepository{
		tableName: schema + "." + "cargoes",
		pool:      pool,
		decoder:   NewPostgresCargoDecoder(),
		encoder:   NewPostgresCargoEncoder(),
		fields: []string{
			"id",
			"vessel_id",
			"items",
			"status",
			"created_at",
			"updated_at",
			"deleted_at",
		},
		errorHandler: postgres.NewErrorHandler(errorHandlers),
	}
}

func (r *PostgresCargoRepository) Find(ctx context.Context, id cargodomain.CargoID) (*cargodomain.Cargo, error) {
	rows, err := r.cargoSelectBuilder(1, sq.Eq{"id": id}).RunWith(r.pool.Reader()).QueryContext(ctx)
	if err != nil {
		return nil, ErrRunningQuery.Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	if !rows.Next() {
		return nil, cargodomain.NewCargoNotExistsError(id)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, ErrFetchingCargoRows.Wrap(rowsErr)
	}

	v, err := r.decoder(rows)
	if err != nil {
		return nil, ErrFetchingCargoRows.Wrap(err)
	}

	return v, nil
}

func (r *PostgresCargoRepository) Save(ctx context.Context, c *cargodomain.Cargo) error {
	bindings, bindingsErr := r.encoder(c)
	if bindingsErr != nil {
		return ErrSavingCargo.Wrap(bindingsErr)
	}

	query := sq.Insert(r.tableName).
		Columns(r.fields...).
		Values(bindings...).
		Suffix("ON CONFLICT (id) DO UPDATE SET " +
			"items = EXCLUDED.items, " +
			"status = EXCLUDED.status, " +
			"updated_at = EXCLUDED.updated_at, " +
			"deleted_at = EXCLUDED.deleted_at",
		).PlaceholderFormat(sq.Dollar)

	_, err := query.RunWith(r.pool.Writer()).ExecContext(ctx)
	if err != nil {
		if pgError, match := postgres.IsPostgresError(err); match {
			return ErrSavingCargo.Wrap(r.errorHandler.Handle(c, pgError))
		}

		return ErrSavingCargo.Wrap(err)
	}

	return nil
}

func (r *PostgresCargoRepository) cargoSelectBuilder(limit uint64, wheres ...sq.Eq) sq.SelectBuilder {
	qb := sq.Select(r.fields...).From(r.tableName).Limit(limit).PlaceholderFormat(sq.Dollar)

	if wheres == nil {
		wheres = make([]sq.Eq, 0)
	}

	wheres = append(wheres, sq.Eq{"deleted_at": nil})

	for _, where := range wheres {
		qb = qb.Where(where)
	}

	return qb
}

func uniqueViolationPostgresCargoRepoErrorHandler() postgres.ErrorHandlerFunc {
	return func(resource interface{}, err *pq.Error) error {
		switch res := resource.(type) {
		case *cargodomain.Cargo:
			return cargodomain.NewCargoAlreadyExistsError(res.ID(), res.VesselID()).Wrap(err)
		case cargodomain.CargoID:
			return cargodomain.NewCargoAlreadyExistsError(res, "").Wrap(err)
		default:
			return err
		}
	}
}
