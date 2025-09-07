package cargopersistence

import (
	"database/sql"
	"time"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	cargotrackingdomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain/tracking"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
)

var (
	ErrScanningCargoTrackingRow = errutil.NewError("error scanning cargo tracking row")
)

func newPostgresCargoTrackingDecoder() postgres.DecodeFunc[cargotrackingdomain.TrackingItem] {
	return func(rows *sql.Rows) (cargotrackingdomain.TrackingItem, error) {
		var (
			id              string
			cargoID         string
			entryType       string
			rawStatusBefore *sql.NullString
			rawStatusAfter  *sql.NullString
			createdAt       time.Time
		)

		err := rows.Scan(
			&id, &cargoID, &entryType, &rawStatusBefore, &rawStatusAfter, &createdAt,
		)
		if err != nil {
			return cargotrackingdomain.TrackingItem{}, ErrScanningCargoTrackingRow.Wrap(err)
		}

		var statusBefore *string
		if rawStatusBefore != nil && rawStatusBefore.Valid {
			statusBefore = &rawStatusBefore.String
		}

		var statusAfter *string
		if rawStatusAfter != nil && rawStatusAfter.Valid {
			statusAfter = &rawStatusAfter.String
		}

		primitives := cargotrackingdomain.TrackingItemPrimitives{
			ID:           id,
			CargoID:      cargoID,
			EntryType:    entryType,
			StatusBefore: statusBefore,
			StatusAfter:  statusAfter,
			CreatedAt:    createdAt,
		}

		return cargotrackingdomain.NewTrackingItemFromPrimitives(primitives), nil
	}
}

func newPostgresCargoTrackingEncoder() func(cargodomain.CargoID) postgres.EncodeFunc[cargotrackingdomain.TrackingItem] {
	return func(cargoID cargodomain.CargoID) postgres.EncodeFunc[cargotrackingdomain.TrackingItem] {
		return func(tracking cargotrackingdomain.TrackingItem) ([]any, error) {
			primitives := cargotrackingdomain.NewTrackingItemPrimitives(cargoID.String(), tracking)

			return []any{
				primitives.ID,
				primitives.CargoID,
				primitives.EntryType,
				primitives.StatusBefore,
				primitives.StatusAfter,
				primitives.CreatedAt,
			}, nil
		}
	}
}
