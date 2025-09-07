package cargopersistence

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	cargotrackingdomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain/tracking"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
)

var (
	ErrFetchingCargoTrackingRows = errutil.NewError("error fetching cargo tracking rows")
	ErrRunningCargoTrackingQuery = errutil.NewError("error running cargo tracking query")
	ErrSavingCargoTracking       = errutil.NewError("error saving cargo to the database")
)

type postgresCargoTrackingRepository struct {
	tableName string
	pool      sqldb.ConnectionPool
	decoder   postgres.DecodeFunc[cargotrackingdomain.TrackingItem]
	encoder   func(cargodomain.CargoID) postgres.EncodeFunc[cargotrackingdomain.TrackingItem]
	fields    []string
}

func newPostgresCargoTrackingRepository(schema string, pool sqldb.ConnectionPool) *postgresCargoTrackingRepository {
	return &postgresCargoTrackingRepository{
		tableName: schema + "." + "cargoes_tracking",
		pool:      pool,
		decoder:   newPostgresCargoTrackingDecoder(),
		encoder:   newPostgresCargoTrackingEncoder(),
		fields: []string{
			"id",
			"cargo_id",
			"entry_type",
			"status_before",
			"status_after",
			"created_at",
		},
	}
}

func (r *postgresCargoTrackingRepository) findByCargoID(
	ctx context.Context,
	cargoID cargodomain.CargoID,
) (cargotrackingdomain.Tracking, error) {
	rows, err := r.trackingSelectBuilder(sq.Eq{"cargo_id": cargoID}).RunWith(r.pool.Reader()).QueryContext(ctx)
	if err != nil {
		return nil, ErrRunningCargoTrackingQuery.Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, ErrFetchingCargoTrackingRows.Wrap(rowsErr)
	}

	trackingItems := make(cargotrackingdomain.Tracking, 0)
	for rows.Next() {
		trackingItem, decodeErr := r.decoder(rows)
		if decodeErr != nil {
			return cargotrackingdomain.Tracking{}, decodeErr
		}
		trackingItems = append(trackingItems, trackingItem)
	}

	return trackingItems, nil
}

func (r *postgresCargoTrackingRepository) save(
	ctx context.Context,
	tx *sql.Tx,
	cargoID cargodomain.CargoID,
	tracking cargotrackingdomain.Tracking,
) error {
	if len(tracking) == 0 {
		return nil
	}

	queryBuilder := sq.Insert(r.tableName).Columns(r.fields...).PlaceholderFormat(sq.Dollar)
	encodeFunc := r.encoder(cargoID)

	for _, item := range tracking {
		values, encodeErr := encodeFunc(item)
		if encodeErr != nil {
			return encodeErr
		}
		queryBuilder = queryBuilder.Values(values...)
	}

	_, err := queryBuilder.RunWith(tx).ExecContext(ctx)
	if err != nil {
		return ErrSavingCargoTracking.Wrap(err)
	}

	return nil
}

func (r *postgresCargoTrackingRepository) trackingSelectBuilder(wheres ...sq.Eq) sq.SelectBuilder {
	qb := sq.Select(r.fields...).From(r.tableName).PlaceholderFormat(sq.Dollar)

	if wheres == nil {
		wheres = make([]sq.Eq, 0)
	}

	for _, where := range wheres {
		qb = qb.Where(where)
	}

	return qb
}
