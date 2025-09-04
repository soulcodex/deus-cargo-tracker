package cargocommands

import (
	"context"
	"fmt"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

type CreateCargoCommand struct {
	ID       string
	VesselID string
	Items    []struct {
		Name   string `json:"name"`
		Weight uint64 `json:"weight"`
	}
}

func (c *CreateCargoCommand) Type() string {
	return "create_cargo_command"
}

type CreateCargoCommandHandler struct {
	creator      *cargodomain.CargoCreator
	timeProvider utils.DateTimeProvider
}

func NewCreateCargoCommandHandler(
	creator *cargodomain.CargoCreator,
	timeProvider utils.DateTimeProvider,
) *CreateCargoCommandHandler {
	return &CreateCargoCommandHandler{
		creator:      creator,
		timeProvider: timeProvider,
	}
}

func (h *CreateCargoCommandHandler) Handle(ctx context.Context, cmd *CreateCargoCommand) (interface{}, error) {
	input := cargodomain.CargoCreateInput{
		ID:       cmd.ID,
		VesselID: cmd.VesselID,
		Items:    h.buildCargoItemInputs(cmd),
		At:       h.timeProvider.Now(),
	}

	if _, err := h.creator.Create(ctx, input); err != nil {
		return nil, fmt.Errorf("error creating cargo: %w", err)
	}

	return struct{}{}, nil
}

func (h *CreateCargoCommandHandler) buildCargoItemInputs(cmd *CreateCargoCommand) []cargodomain.CargoItemInput {
	items := make([]cargodomain.CargoItemInput, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		items = append(items, cargodomain.CargoItemInput{
			Name:   item.Name,
			Weight: item.Weight,
		})
	}

	return items
}
