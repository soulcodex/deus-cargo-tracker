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
		setupMocks    func(repo *cargodomainmock.CargoRepositoryMock)
		expectedError string
	}{
		{
			name: "should create cargo successfully",
			input: cargodomain.CargoCreateInput{
				ID: idProvider.New().String(),
				Items: []cargodomain.CargoItemInput{
					{Name: "Fuel", Weight: 100},
					{Name: "Supplies", Weight: 50},
				},
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock) {
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return nil, nil // not found
				}
				repo.SaveFunc = func(ctx context.Context, c *cargodomain.Cargo) error {
					return nil
				}
			},
		},
		{
			name: "should fail when cargo ID is invalid",
			input: cargodomain.CargoCreateInput{
				ID: "!!!invalid-id###",
				Items: []cargodomain.CargoItemInput{
					{Name: "Fuel", Weight: 100},
				},
			},
			setupMocks:    func(_ *cargodomainmock.CargoRepositoryMock) {}, // no repo interaction
			expectedError: "invalid cargo id",
		},
		{
			name: "should fail if cargo already exists",
			input: cargodomain.CargoCreateInput{
				ID: idProvider.New().String(),
				Items: []cargodomain.CargoItemInput{
					{Name: "Fuel", Weight: 100},
				},
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock) {
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return &cargodomain.Cargo{}, nil
				}
			},
			expectedError: "cargo already exists",
		},
		{
			name: "should fail if find returns an error",
			input: cargodomain.CargoCreateInput{
				ID: idProvider.New().String(),
				Items: []cargodomain.CargoItemInput{
					{Name: "Fuel", Weight: 100},
				},
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock) {
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return nil, errors.New("db lookup failed")
				}
			},
			expectedError: "error checking existing cargo: db lookup failed",
		},
		{
			name: "should fail when creating invalid items",
			input: cargodomain.CargoCreateInput{
				ID: idProvider.New().String(),
				Items: []cargodomain.CargoItemInput{
					{Name: "", Weight: 0}, // invalid item
				},
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock) {
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return nil, nil
				}
			},
			expectedError: "error creating cargo items",
		},
		{
			name: "should fail when cargo exceed max items allowed",
			input: cargodomain.CargoCreateInput{
				ID: idProvider.New().String(),
				Items: []cargodomain.CargoItemInput{
					{Name: "Fuel", Weight: 100},
					{Name: "Food", Weight: 100},
					{Name: "Oil", Weight: 100},
					{Name: "Medicine", Weight: 100},
					{Name: "Shampoo", Weight: 100},
					{Name: "Soap", Weight: 100},
					{Name: "Clothes", Weight: 100},
					{Name: "Wine", Weight: 100},
					{Name: "Whisky", Weight: 100},
					{Name: "Fruit", Weight: 100},
					{Name: "Toys", Weight: 100},
				},
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock) {
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return nil, nil
				}
			},
			expectedError: "error creating cargo items",
		},
		{
			name: "should fail if save returns an error",
			input: cargodomain.CargoCreateInput{
				ID: idProvider.New().String(),
				Items: []cargodomain.CargoItemInput{
					{Name: "Fuel", Weight: 100},
				},
			},
			setupMocks: func(repo *cargodomainmock.CargoRepositoryMock) {
				repo.FindFunc = func(ctx context.Context, id cargodomain.CargoID) (*cargodomain.Cargo, error) {
					return nil, nil
				}
				repo.SaveFunc = func(ctx context.Context, c *cargodomain.Cargo) error {
					return errors.New("db error")
				}
			},
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &cargodomainmock.CargoRepositoryMock{}
			if tt.setupMocks != nil {
				tt.setupMocks(repo)
			}

			creator := cargodomain.NewCargoCreator(repo)
			cargo, err := creator.Create(ctx, tt.input, timeProvider.Now())

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
