package cargodomain

import (
	"context"

	"github.com/soulcodex/deus-cargo-tracker/pkg/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

var (
	ErrCargoUpdateFailed = errutil.NewError("cargo update failed")
)

type CargoUpdater struct {
	repository CargoRepository
	publisher  domain.EventPublisher
}

func NewCargoUpdater(repository CargoRepository, publisher domain.EventPublisher) *CargoUpdater {
	return &CargoUpdater{
		repository: repository,
		publisher:  publisher,
	}
}

func (cu *CargoUpdater) Update(ctx context.Context, id string, opts ...CargoUpdateOpt) error {
	cargoID, err := NewCargoID(id)
	if err != nil {
		return ErrCargoUpdateFailed.Wrap(err)
	}

	cargo, err := cu.repository.Find(ctx, cargoID)
	if err != nil {
		return ErrCargoUpdateFailed.Wrap(err)
	}

	if updateErr := cargo.Update(ctx, opts...); updateErr != nil {
		return ErrCargoUpdateFailed.Wrap(updateErr)
	}

	// The following approach would be better if we had an outbox pattern implemented.
	// However, for simplicity, we are directly publishing events after saving the cargo.
	events := cargo.PullEvents()
	if saveErr := cu.repository.Save(ctx, cargo); saveErr != nil {
		return ErrCargoUpdateFailed.Wrap(saveErr)
	}

	if publishErr := cu.publisher.Publish(ctx, events...); publishErr != nil {
		return ErrCargoUpdateFailed.Wrap(publishErr)
	}

	return nil
}
