package vesselpersistence

import (
	"database/sql"
	"time"

	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
)

var (
	ErrScanningVesselRow = errutil.NewError("error scanning vessel row")
)

func NewPostgresVesselDecoder() postgres.DecodeFunc[*vesseldomain.Vessel] {
	return func(rows *sql.Rows) (*vesseldomain.Vessel, error) {
		var (
			id           string
			name         string
			capacity     uint64
			latitude     float64
			longitude    float64
			createdAt    time.Time
			updatedAt    time.Time
			rawDeletedAt sql.NullTime
		)

		err := rows.Scan(
			&id, &name, &capacity, &latitude, &longitude,
			&createdAt, &updatedAt, &rawDeletedAt,
		)
		if err != nil {
			return nil, ErrScanningVesselRow.Wrap(err)
		}

		var deletedAt *time.Time
		if rawDeletedAt.Valid && !rawDeletedAt.Time.IsZero() {
			deletedAt = &rawDeletedAt.Time
		}

		primitives := vesseldomain.VesselPrimitives{
			ID:        id,
			Name:      name,
			Capacity:  capacity,
			Latitude:  latitude,
			Longitude: longitude,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			DeletedAt: deletedAt,
		}

		return vesseldomain.NewVesselFromPrimitives(primitives), nil
	}
}
