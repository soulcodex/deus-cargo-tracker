package cargotest

import (
	"testing"
	"time"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
)

type CargoMotherOpt func(*CargoMother)

func WithID(id string) CargoMotherOpt {
	return func(m *CargoMother) {
		m.primitives.ID = id
	}
}

func WithVesselID(id string) CargoMotherOpt {
	return func(m *CargoMother) {
		m.primitives.VesselID = id
	}
}

func WithSoftDeletion(at time.Time) CargoMotherOpt {
	return func(m *CargoMother) {
		m.primitives.DeletedAt = &at
	}
}

type CargoMother struct {
	primitives cargodomain.CargoPrimitives
}

func NewCargoMother(opts ...CargoMotherOpt) *CargoMother {
	mother := &CargoMother{
		primitives: newCargoPrimitives(),
	}

	for _, opt := range opts {
		opt(mother)
	}

	return mother
}

func (m *CargoMother) Build(t *testing.T) *cargodomain.Cargo {
	t.Helper()

	return cargodomain.NewCargoFromPrimitives(m.primitives)
}

func newCargoPrimitives() cargodomain.CargoPrimitives {
	at := time.Now()

	const (
		itemOneWeight = 1500
		itemTwoWeight = 2000
	)

	return cargodomain.CargoPrimitives{
		ID:       "01K43FJ8ZCYAVQ14ZV7EKCPMR8",
		VesselID: "01K4B43REGN4HBFQETVZZ484A3",
		Items: []cargodomain.ItemsPrimitives{
			{
				Name:   "Electronics",
				Weight: uint64(itemOneWeight),
			},
			{
				Name:   "Clothing",
				Weight: uint64(itemTwoWeight),
			},
		},
		Status:    cargodomain.StatusPending.String(),
		CreatedAt: at,
		UpdatedAt: at,
		DeletedAt: nil,
	}
}
