package cargopersistence

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	cargotrackingdomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain/tracking"
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
	encoder      postgres.EncodeFunc[*cargodomain.Cargo]
	errorHandler *postgres.ErrorHandler
	fields       []string

	trackingRepo *postgresCargoTrackingRepository
}

func NewPostgresCargoRepository(schema string, pool sqldb.ConnectionPool) *PostgresCargoRepository {
	errorHandlers := postgres.ErrorHandlers{
		postgres.UniqueViolationErrorCode: uniqueViolationPostgresCargoRepoErrorHandler(),
	}

	return &PostgresCargoRepository{
		tableName: schema + "." + "cargoes",
		pool:      pool,
		encoder:   newPostgresCargoEncoder(),
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
		trackingRepo: newPostgresCargoTrackingRepository(schema, pool),
	}
}

func (r *PostgresCargoRepository) Find(
	ctx context.Context,
	id cargodomain.CargoID,
	opts ...cargodomain.CargoFindingOpt,
) (*cargodomain.Cargo, error) {
	options := cargodomain.NewCargoFindingOptions(opts...)

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

	tracking := cargotrackingdomain.NewEmptyTracking()
	if options.WithTracking {
		tracking, err = r.trackingRepo.findByCargoID(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	v, err := newPostgresCargoDecoder(tracking)(rows)
	if err != nil {
		return nil, ErrFetchingCargoRows.Wrap(err)
	}

	return v, nil
}

func (r *PostgresCargoRepository) Save(ctx context.Context, c *cargodomain.Cargo) error {
	tx, txErr := r.pool.Writer().BeginTx(ctx, nil)
	if txErr != nil {
		return ErrSavingCargo.Wrap(txErr)
	}

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

	_, err := query.RunWith(tx).ExecContext(ctx)
	if err != nil {
		if pgError, match := postgres.IsPostgresError(err); match {
			return ErrSavingCargo.Wrap(r.errorHandler.Handle(c, pgError))
		}

		return ErrSavingCargo.Wrap(err)
	}

	if saveTrackingErr := r.trackingRepo.save(ctx, tx, c.ID(), c.Tracking()); saveTrackingErr != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return ErrSavingCargo.Wrap(rbErr)
		}

		return ErrSavingCargo.Wrap(saveTrackingErr)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return ErrSavingCargo.Wrap(commitErr)
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
