package cargopersistence

import (
	"database/sql"
	"encoding/json"
	"time"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	cargotrackingdomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain/tracking"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
)

var (
	ErrScanningCargoRow     = errutil.NewError("error scanning cargo row")
	ErrInvalidCargoEncoding = errutil.NewError("invalid cargo provided on encoding")
)

func newPostgresCargoDecoder(tracking cargotrackingdomain.Tracking) postgres.DecodeFunc[*cargodomain.Cargo] {
	return func(rows *sql.Rows) (*cargodomain.Cargo, error) {
		var (
			id           string
			vesselID     string
			items        sql.RawBytes
			status       string
			createdAt    time.Time
			updatedAt    time.Time
			rawDeletedAt sql.NullTime
		)

		err := rows.Scan(
			&id, &vesselID, &items, &status,
			&createdAt, &updatedAt, &rawDeletedAt,
		)
		if err != nil {
			return nil, ErrScanningCargoRow.Wrap(err)
		}

		var deletedAt *time.Time
		if rawDeletedAt.Valid && !rawDeletedAt.Time.IsZero() {
			deletedAt = &rawDeletedAt.Time
		}

		var cargoItems []cargodomain.ItemsPrimitives
		unmarshalErr := json.Unmarshal(items, &cargoItems)
		if unmarshalErr != nil {
			return nil, ErrScanningCargoRow.Wrap(err)
		}

		primitives := cargodomain.CargoPrimitives{
			ID:        id,
			VesselID:  vesselID,
			Items:     cargoItems,
			Tracking:  cargotrackingdomain.NewTrackingPrimitives(id, tracking),
			Status:    status,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			DeletedAt: deletedAt,
		}

		return cargodomain.NewCargoFromPrimitives(primitives), nil
	}
}

func newPostgresCargoEncoder() postgres.EncodeFunc[*cargodomain.Cargo] {
	return func(cargo *cargodomain.Cargo) ([]any, error) {
		if cargo == nil {
			return nil, ErrInvalidCargoEncoding
		}

		primitives := cargo.Primitives()

		items, marshalErr := json.Marshal(primitives.Items)
		if marshalErr != nil {
			return nil, ErrSavingCargo.Wrap(marshalErr)
		}

		encoded := []any{
			primitives.ID,
			primitives.VesselID,
			sql.RawBytes(items),
			primitives.Status,
			primitives.CreatedAt,
			primitives.UpdatedAt,
			primitives.DeletedAt,
		}

		return encoded, nil
	}
}
