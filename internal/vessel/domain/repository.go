package vesseldomain

import (
	"context"
)

type VesselRepositoryReader interface {
	Find(ctx context.Context, id VesselID) (*Vessel, error)
}

type VesselRepositoryWriter interface {
	Save(ctx context.Context, v *Vessel) error
}

//go:generate moq -pkg vesseldomainmock -out mock/vessel_repository_moq.go . VesselRepository
type VesselRepository interface {
	VesselRepositoryReader
	VesselRepositoryWriter
}
