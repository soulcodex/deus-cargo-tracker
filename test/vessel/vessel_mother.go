package vesseltest

import (
	"testing"
	"time"

	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
)

type VesselMotherOpt func(*VesselMother)

func WithVesselID(id string) VesselMotherOpt {
	return func(m *VesselMother) {
		m.primitives.ID = id
	}
}

func WithSoftDeletion(at time.Time) VesselMotherOpt {
	return func(m *VesselMother) {
		m.primitives.DeletedAt = &at
	}
}

type VesselMother struct {
	primitives vesseldomain.VesselPrimitives
}

func NewVesselMother(opts ...VesselMotherOpt) *VesselMother {
	mother := &VesselMother{
		primitives: newRocketPrimitives(),
	}

	for _, opt := range opts {
		opt(mother)
	}

	return mother
}

func (m *VesselMother) Build(t *testing.T) *vesseldomain.Vessel {
	t.Helper()

	return vesseldomain.NewVesselFromPrimitives(m.primitives)
}

func newRocketPrimitives() vesseldomain.VesselPrimitives {
	at := time.Now()

	const (
		defaultCapacityInKilograms = 5000
		defaultLatitude            = 37.7749
		defaultLongitude           = -122.4194
	)

	return vesseldomain.VesselPrimitives{
		ID:        "01K43FJ8ZCYAVQ14ZV7EKCPMR8",
		Name:      "Falcon 9",
		Capacity:  defaultCapacityInKilograms,
		Latitude:  defaultLatitude,
		Longitude: defaultLongitude,
		CreatedAt: at,
		UpdatedAt: at,
		DeletedAt: nil,
	}
}
