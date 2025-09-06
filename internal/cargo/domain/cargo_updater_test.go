package cargodomain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	cargodomainmock "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain/mock"
	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
	cargotest "github.com/soulcodex/deus-cargo-tracker/test/cargo"
)

func TestCargoUpdater_Update(t *testing.T) {
	ctx := context.Background()
	idProvider := utils.NewFixedULIDProvider()

	now := time.Now()
	later := now.Add(1 * time.Hour)
	earlier := now.Add(-1 * time.Hour)

	tests := []struct {
		name          string
		setupCargo    func() *cargodomain.Cargo
		id            string
		at            time.Time
		opts          []cargodomain.CargoUpdateOpt
		setupMocks    func(repo *cargodomainmock.CargoRepositoryMock, cargo *cargodomain.Cargo)
		expectedError string
	}{
		{
			name: "should update cargo successfully",
			id:   idProvider.New().String(),
			at:   later,
			setupCargo: func() *cargodomain.Cargo {
				return cargotest.NewCargoMother(
					cargotest.WithID(idProvider.New().String()),
					cargotest.WithVesselID(idProvider.New().String()),
					cargotest.WithTimestamps(now, now),
				).Build(t)
			},
			opts: []cargodomain.CargoUpdateOpt{cargodomain.WithStatus("in_transit")},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, cargo *cargodomain.Cargo) {
				repo.FindFunc = func(_ context.Context, _ cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return cargo, nil
				}
				repo.SaveFunc = func(_ context.Context, _ *cargodomain.Cargo) error {
					return nil
				}
			},
		},
		{
			name:          "should fail when cargo ID is invalid",
			id:            "!!!invalid-id###",
			at:            now,
			setupCargo:    func() *cargodomain.Cargo { return nil },
			setupMocks:    func(_ *cargodomainmock.CargoRepositoryMock, _ *cargodomain.Cargo) {},
			expectedError: "cargo update failed",
		},
		{
			name:       "should fail when cargo not found",
			id:         idProvider.New().String(),
			at:         now,
			setupCargo: func() *cargodomain.Cargo { return nil },
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, _ *cargodomain.Cargo) {
				repo.FindFunc = func(_ context.Context, _ cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return nil, errors.New("not found")
				}
			},
			expectedError: "cargo update failed: not found",
		},
		{
			name: "should skip update when timestamp is older",
			id:   idProvider.New().String(),
			at:   earlier,
			setupCargo: func() *cargodomain.Cargo {
				return cargotest.NewCargoMother(
					cargotest.WithID(idProvider.New().String()),
					cargotest.WithVesselID(idProvider.New().String()),
					cargotest.WithTimestamps(now, now),
				).Build(t)
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, cargo *cargodomain.Cargo) {
				repo.FindFunc = func(_ context.Context, _ cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return cargo, nil
				}

				repo.SaveFunc = func(_ context.Context, _ *cargodomain.Cargo) error {
					return nil
				}
			},
		},
		{
			name: "should fail when cargo is deleted",
			id:   idProvider.New().String(),
			at:   later,
			setupCargo: func() *cargodomain.Cargo {
				return cargotest.NewCargoMother(
					cargotest.WithID(idProvider.New().String()),
					cargotest.WithVesselID(idProvider.New().String()),
					cargotest.WithSoftDeletion(time.Now()),
				).Build(t)
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, cargo *cargodomain.Cargo) {
				repo.FindFunc = func(_ context.Context, _ cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return cargo, nil
				}
			},
			expectedError: "cargo update failed: cargo is not modifiable",
		},
		{
			name: "should fail when update option fails",
			id:   idProvider.New().String(),
			at:   later,
			setupCargo: func() *cargodomain.Cargo {
				return cargotest.NewCargoMother(
					cargotest.WithID(idProvider.New().String()),
					cargotest.WithVesselID(idProvider.New().String()),
					cargotest.WithTimestamps(now, now),
				).Build(t)
			},
			opts: []cargodomain.CargoUpdateOpt{
				func(_ *cargodomain.Cargo) error { return errors.New("bad update") },
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, cargo *cargodomain.Cargo) {
				repo.FindFunc = func(_ context.Context, _ cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return cargo, nil
				}
			},
			expectedError: "cargo update failed: bad update",
		},
		{
			name: "should fail when save fails",
			id:   idProvider.New().String(),
			at:   later,
			setupCargo: func() *cargodomain.Cargo {
				return cargotest.NewCargoMother(
					cargotest.WithID(idProvider.New().String()),
					cargotest.WithVesselID(idProvider.New().String()),
					cargotest.WithTimestamps(now, now),
				).Build(t)
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, cargo *cargodomain.Cargo) {
				repo.FindFunc = func(_ context.Context, _ cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return cargo, nil
				}
				repo.SaveFunc = func(_ context.Context, _ *cargodomain.Cargo) error {
					return errors.New("db error")
				}
			},
			expectedError: "cargo update failed: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &cargodomainmock.CargoRepositoryMock{}
			cargo := tt.setupCargo()
			if tt.setupMocks != nil {
				tt.setupMocks(repo, cargo)
			}

			updater := cargodomain.NewCargoUpdater(repo)
			err := updater.Update(ctx, tt.id, tt.at, tt.opts...)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
