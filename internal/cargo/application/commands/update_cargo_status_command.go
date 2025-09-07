package cargocommands

import (
	"context"
	"errors"
	"fmt"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

type UpdateCargoStatusCommand struct {
	ID        string
	NewStatus string
}

func (c *UpdateCargoStatusCommand) Type() string {
	return "update_cargo_status_command"
}

func (c *UpdateCargoStatusCommand) BlockingKey() string {
	return "cargo_update:" + c.ID
}

type UpdateCargoStatusCommandHandler struct {
	updater      *cargodomain.CargoUpdater
	timeProvider utils.DateTimeProvider
	idProvider   utils.ULIDProvider
}

func NewUpdateCargoStatusCommandHandler(
	updater *cargodomain.CargoUpdater,
	timeProvider utils.DateTimeProvider,
	idProvider utils.ULIDProvider,
) *UpdateCargoStatusCommandHandler {
	return &UpdateCargoStatusCommandHandler{
		updater:      updater,
		timeProvider: timeProvider,
		idProvider:   idProvider,
	}
}

func (h *UpdateCargoStatusCommandHandler) Handle(ctx context.Context, cmd *UpdateCargoStatusCommand) (interface{}, error) {
	trackingID, at := h.idProvider.New().String(), h.timeProvider.Now()

	updates := []cargodomain.CargoUpdateOpt{
		cargodomain.WithStatus(trackingID, cmd.NewStatus, at),
	}

	if err := h.updater.Update(ctx, cmd.ID, updates...); err != nil {
		if errors.Is(err, cargodomain.ErrStatusUnchanged) {
			return struct{}{}, nil
		}

		return nil, fmt.Errorf("error updating cargo status: %w", err)
	}

	return struct{}{}, nil
}
