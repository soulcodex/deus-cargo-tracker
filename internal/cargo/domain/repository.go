package cargodomain

import (
	"context"
)

type CargoFindingOpt func(*CargoFindingOptions)
type CargoFindingOptions struct {
	WithTracking bool
}

func newDefaultCargoFindingOptions() *CargoFindingOptions {
	return &CargoFindingOptions{
		WithTracking: false,
	}
}

func NewCargoFindingOptions(opts ...CargoFindingOpt) *CargoFindingOptions {
	options := newDefaultCargoFindingOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithTracking() CargoFindingOpt {
	return func(o *CargoFindingOptions) {
		o.WithTracking = true
	}
}

type CargoRepositoryReader interface {
	Find(ctx context.Context, id CargoID, opts ...CargoFindingOpt) (*Cargo, error)
}

type CargoRepositoryWriter interface {
	Save(ctx context.Context, c *Cargo) error
}

//go:generate moq -pkg cargodomainmock -out mock/cargo_repository_moq.go . CargoRepository
type CargoRepository interface {
	CargoRepositoryReader
	CargoRepositoryWriter
}
