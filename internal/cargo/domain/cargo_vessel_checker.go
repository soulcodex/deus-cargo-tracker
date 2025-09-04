package cargodomain

import (
	"context"
)

//go:generate moq -pkg cargodomainmock -out mock/cargo_vessel_checker_moq.go . CargoVesselChecker
type CargoVesselChecker interface {
	Check(ctx context.Context, vesselID VesselID) error
}
