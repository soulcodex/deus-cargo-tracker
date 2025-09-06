package cargocommands

import (
	"context"
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
	return "update_cargo_status_command:" + c.ID
}

type UpdateCargoStatusCommandHandler struct {
	updater      *cargodomain.CargoUpdater
	timeProvider utils.DateTimeProvider
}

func NewUpdateCargoStatusCommandHandler(
	updater *cargodomain.CargoUpdater,
	timeProvider utils.DateTimeProvider,
) *UpdateCargoStatusCommandHandler {
	return &UpdateCargoStatusCommandHandler{
		updater:      updater,
		timeProvider: timeProvider,
	}
}

func (h *UpdateCargoStatusCommandHandler) Handle(ctx context.Context, cmd *UpdateCargoStatusCommand) (interface{}, error) {
	updates := []cargodomain.CargoUpdateOpt{
		cargodomain.WithStatus(cmd.NewStatus),
	}

	if err := h.updater.Update(ctx, cmd.ID, h.timeProvider.Now(), updates...); err != nil {
		return nil, fmt.Errorf("error updating cargo status: %w", err)
	}

	return struct{}{}, nil
}
