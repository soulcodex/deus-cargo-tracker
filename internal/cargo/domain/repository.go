package cargodomain

import (
	"context"
)

type CargoRepositoryReader interface {
	Find(ctx context.Context, id CargoID) (*Cargo, error)
}

type CargoRepositoryWriter interface {
	Save(ctx context.Context, c *Cargo) error
}

//go:generate moq -pkg cargodomainmock -out mock/cargo_repository_moq.go . CargoRepository
type CargoRepository interface {
	CargoRepositoryReader
	CargoRepositoryWriter
}
