package cargodomain

import (
	"context"
	"time"

	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

var (
	ErrCargoUpdateFailed = errutil.NewError("cargo update failed")
)

type CargoUpdater struct {
	repository CargoRepository
}

func NewCargoUpdater(repository CargoRepository) *CargoUpdater {
	return &CargoUpdater{
		repository: repository,
	}
}

func (cu *CargoUpdater) Update(ctx context.Context, id string, at time.Time, opts ...CargoUpdateOpt) error {
	cargoID, err := NewCargoID(id)
	if err != nil {
		return ErrCargoUpdateFailed.Wrap(err)
	}

	cargo, err := cu.repository.Find(ctx, cargoID)
	if err != nil {
		return ErrCargoUpdateFailed.Wrap(err)
	}

	if updateErr := cargo.Update(ctx, at, opts...); updateErr != nil {
		return ErrCargoUpdateFailed.Wrap(updateErr)
	}

	if saveErr := cu.repository.Save(ctx, cargo); saveErr != nil {
		return ErrCargoUpdateFailed.Wrap(saveErr)
	}

	return nil
}
