package cargodomain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	cargodomainmock "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain/mock"
	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

func TestCargoCreator_Create(t *testing.T) {
	ctx, timeProvider, idProvider := context.Background(), utils.NewFixedTimeProvider(), utils.NewFixedULIDProvider()

	tests := []struct {
		name          string
		input         cargodomain.CargoCreateInput
		setupMocks    func(repo *cargodomainmock.CargoRepositoryMock, checker *cargodomainmock.CargoVesselCheckerMock)
		expectedError string
	}{
		{
			name: "should create cargo successfully",
			input: cargodomain.CargoCreateInput{
				ID:       idProvider.New().String(),
				VesselID: idProvider.New().String(),
				Items: []struct {
					Name   string
					Weight uint64
				}{
					{Name: "Fuel", Weight: 100},
					{Name: "Supplies", Weight: 50},
				},
				At: timeProvider.Now(),
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, checker *cargodomainmock.CargoVesselCheckerMock) {
				checker.CheckFunc = func(ctx context.Context, id cargodomain.VesselID) error {
					return nil
				}
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID, opts ...cargodomain.CargoFindingOpt) (*cargodomain.Cargo, error) {
					return nil, nil
				}
				repo.SaveFunc = func(ctx context.Context, c *cargodomain.Cargo) error {
					return nil
				}
			},
		},
		{
			name: "should fail when vessel check fails",
			input: cargodomain.CargoCreateInput{
				ID:       idProvider.New().String(),
				VesselID: idProvider.New().String(),
				Items: []struct {
					Name   string
					Weight uint64
				}{
					{Name: "Fuel", Weight: 100},
				},
				At: timeProvider.Now(),
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, checker *cargodomainmock.CargoVesselCheckerMock) {
				checker.CheckFunc = func(ctx context.Context, id cargodomain.VesselID) error {
					return errors.New("vessel validation failed")
				}
			},
			expectedError: "error checking vessel",
		},
		{
			name: "should fail when cargo ID is invalid",
			input: cargodomain.CargoCreateInput{
				ID:       "!!!invalid-id###",
				VesselID: idProvider.New().String(),
				Items: []struct {
					Name   string
					Weight uint64
				}{
					{Name: "Fuel", Weight: 100},
				},
				At: timeProvider.Now(),
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, checker *cargodomainmock.CargoVesselCheckerMock) {
				checker.CheckFunc = func(ctx context.Context, id cargodomain.VesselID) error {
					return nil
				}
			},
			expectedError: "invalid cargo id",
		},
		{
			name: "should fail if cargo already exists",
			input: cargodomain.CargoCreateInput{
				ID:       idProvider.New().String(),
				VesselID: idProvider.New().String(),
				Items: []struct {
					Name   string
					Weight uint64
				}{
					{Name: "Fuel", Weight: 100},
				},
				At: timeProvider.Now(),
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, checker *cargodomainmock.CargoVesselCheckerMock) {
				checker.CheckFunc = func(ctx context.Context, id cargodomain.VesselID) error {
					return nil
				}
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID, opts ...cargodomain.CargoFindingOpt) (*cargodomain.Cargo, error) {
					return &cargodomain.Cargo{}, nil
				}
			},
			expectedError: "cargo already exists",
		},
		{
			name: "should fail if find returns an error",
			input: cargodomain.CargoCreateInput{
				ID:       idProvider.New().String(),
				VesselID: idProvider.New().String(),
				Items: []struct {
					Name   string
					Weight uint64
				}{
					{Name: "Fuel", Weight: 100},
				},
				At: timeProvider.Now(),
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, checker *cargodomainmock.CargoVesselCheckerMock) {
				checker.CheckFunc = func(ctx context.Context, id cargodomain.VesselID) error {
					return nil
				}
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID, opts ...cargodomain.CargoFindingOpt) (*cargodomain.Cargo, error) {
					return nil, errors.New("db lookup failed")
				}
			},
			expectedError: "error checking existing cargo: db lookup failed",
		},
		{
			name: "should fail when creating invalid items",
			input: cargodomain.CargoCreateInput{
				ID:       idProvider.New().String(),
				VesselID: idProvider.New().String(),
				Items: []struct {
					Name   string
					Weight uint64
				}{
					{Name: "", Weight: 0},
				},
				At: timeProvider.Now(),
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, checker *cargodomainmock.CargoVesselCheckerMock) {
				checker.CheckFunc = func(ctx context.Context, id cargodomain.VesselID) error {
					return nil
				}
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID, opts ...cargodomain.CargoFindingOpt) (*cargodomain.Cargo, error) {
					return nil, nil
				}
			},
			expectedError: "error creating cargo items",
		},
		{
			name: "should fail when cargo exceed max items allowed",
			input: cargodomain.CargoCreateInput{
				ID:       idProvider.New().String(),
				VesselID: idProvider.New().String(),
				Items: []struct {
					Name   string
					Weight uint64
				}{
					{Name: "A", Weight: 100}, {Name: "B", Weight: 100}, {Name: "C", Weight: 100},
					{Name: "D", Weight: 100}, {Name: "E", Weight: 100}, {Name: "F", Weight: 100},
					{Name: "G", Weight: 100}, {Name: "H", Weight: 100}, {Name: "I", Weight: 100},
					{Name: "J", Weight: 100}, {Name: "K", Weight: 100}, // Exceeds limit
				},
				At: timeProvider.Now(),
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, checker *cargodomainmock.CargoVesselCheckerMock) {
				checker.CheckFunc = func(ctx context.Context, id cargodomain.VesselID) error {
					return nil
				}
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID, opts ...cargodomain.CargoFindingOpt) (*cargodomain.Cargo, error) {
					return nil, nil
				}
			},
			expectedError: "error creating cargo items",
		},
		{
			name: "should fail if save returns an error",
			input: cargodomain.CargoCreateInput{
				ID:       idProvider.New().String(),
				VesselID: idProvider.New().String(),
				Items: []struct {
					Name   string
					Weight uint64
				}{
					{Name: "Fuel", Weight: 100},
				},
				At: timeProvider.Now(),
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock, checker *cargodomainmock.CargoVesselCheckerMock) {
				checker.CheckFunc = func(ctx context.Context, id cargodomain.VesselID) error {
					return nil
				}
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID, opts ...cargodomain.CargoFindingOpt) (*cargodomain.Cargo, error) {
					return nil, nil
				}
				repo.SaveFunc = func(ctx context.Context, c *cargodomain.Cargo) error {
					return errors.New("db error")
				}
			},
			expectedError: "error saving cargo: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &cargodomainmock.CargoRepositoryMock{}
			checker := &cargodomainmock.CargoVesselCheckerMock{}

			if tt.setupMocks != nil {
				tt.setupMocks(repo, checker)
			}

			creator := cargodomain.NewCargoCreator(repo, checker, idProvider)
			cargo, err := creator.Create(ctx, tt.input)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, cargo)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, cargo)
			}
		})
	}
}
