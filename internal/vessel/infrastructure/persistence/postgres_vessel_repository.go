package vesselpersistence

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
)

var (
	_ vesseldomain.VesselRepository = (*PostgresVesselRepository)(nil)

	ErrFetchingVesselRows = errutil.NewError("error fetching vessel rows")
	ErrRunningQuery       = errutil.NewError("error running vessel query")
	ErrSavingVessel       = errutil.NewError("error saving vessel to the database")
)

type PostgresVesselRepository struct {
	tableName    string
	pool         sqldb.ConnectionPool
	decoder      postgres.DecodeFunc[*vesseldomain.Vessel]
	errorHandler *postgres.ErrorHandler
	fields       []string
}

func NewPostgresVesselRepository(schema string, pool sqldb.ConnectionPool) *PostgresVesselRepository {
	errorHandlers := postgres.ErrorHandlers{
		postgres.UniqueViolationErrorCode: uniqueViolationPostgresVesselRepoErrorHandler(),
	}

	return &PostgresVesselRepository{
		tableName: schema + "." + "vessels",
		pool:      pool,
		decoder:   NewPostgresVesselDecoder(),
		fields: []string{
			"id",
			"name",
			"capacity",
			"latitude",
			"longitude",
			"created_at",
			"updated_at",
			"deleted_at",
		},
		errorHandler: postgres.NewErrorHandler(errorHandlers),
	}
}

func (r *PostgresVesselRepository) Find(ctx context.Context, id vesseldomain.VesselID) (*vesseldomain.Vessel, error) {
	rows, err := r.vesselSelectBuilder(1, sq.Eq{"id": id}).RunWith(r.pool.Reader()).QueryContext(ctx)
	if err != nil {
		return nil, ErrRunningQuery.Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	if !rows.Next() {
		return nil, vesseldomain.NewVesselNotExistsError(id)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, ErrFetchingVesselRows.Wrap(rowsErr)
	}

	v, err := r.decoder(rows)
	if err != nil {
		return nil, ErrFetchingVesselRows.Wrap(err)
	}

	return v, nil
}

func (r *PostgresVesselRepository) Save(ctx context.Context, v *vesseldomain.Vessel) error {
	primitives := v.Primitives()

	query := sq.Insert(r.tableName).
		Columns(r.fields...).
		Values(
			primitives.ID,
			primitives.Name,
			primitives.Capacity,
			primitives.Latitude,
			primitives.Longitude,
			primitives.CreatedAt,
			primitives.UpdatedAt,
			primitives.DeletedAt,
		).Suffix("ON CONFLICT (id) DO UPDATE SET " +
		"name = EXCLUDED.name, " +
		"capacity = EXCLUDED.capacity, " +
		"latitude = EXCLUDED.latitude, " +
		"longitude = EXCLUDED.longitude, " +
		"updated_at = EXCLUDED.updated_at, " +
		"deleted_at = EXCLUDED.deleted_at",
	).PlaceholderFormat(sq.Dollar)

	_, err := query.RunWith(r.pool.Writer()).ExecContext(ctx)
	if err != nil {
		if pgError, match := postgres.IsPostgresError(err); match {
			return ErrSavingVessel.Wrap(r.errorHandler.Handle(v, pgError))
		}

		return ErrSavingVessel.Wrap(err)
	}

	return nil
}

func (r *PostgresVesselRepository) vesselSelectBuilder(limit uint64, wheres ...sq.Eq) sq.SelectBuilder {
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

func uniqueViolationPostgresVesselRepoErrorHandler() postgres.ErrorHandlerFunc {
	return func(resource interface{}, err *pq.Error) error {
		switch res := resource.(type) {
		case *vesseldomain.Vessel:
			return vesseldomain.NewVesselAlreadyExistsError(res.ID()).Wrap(err)
		case vesseldomain.VesselID:
			return vesseldomain.NewVesselAlreadyExistsError(res).Wrap(err)
		default:
			return err
		}
	}
}
