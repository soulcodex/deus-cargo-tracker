package cargodomain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/soulcodex/deus-cargo-tracker/pkg/messaging"
)

const (
	CargoStatusUpdaterV1DomainEventName = "cargo-status-updated-v1"
)

type CargoStatusUpdatedV1DomainEvent struct {
	*messaging.BaseMessage

	oldStatus  string
	newStatus  string
	occurredOn time.Time
}

func (e *CargoStatusUpdatedV1DomainEvent) OldStatus() string {
	return e.oldStatus
}

func (e *CargoStatusUpdatedV1DomainEvent) NewStatus() string {
	return e.newStatus
}

func (e *CargoStatusUpdatedV1DomainEvent) OccurredOn() time.Time {
	return e.occurredOn
}

func NewCargoStatusUpdatedV1DomainEvent(
	id CargoID,
	oldStatus Status,
	newStatus Status,
	occurredOn time.Time,
) (*CargoStatusUpdatedV1DomainEvent, error) {
	attributes := map[string]any{
		"old_status":  oldStatus.String(),
		"new_status":  newStatus.String(),
		"occurred_on": occurredOn.Format(time.RFC3339),
	}

	eventData, err := json.Marshal(attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cargo status updated v1 domain event: %w", err)
	}

	return &CargoStatusUpdatedV1DomainEvent{
		BaseMessage: messaging.NewBaseMessage(
			CargoStatusUpdaterV1DomainEventName,
			messaging.DefaultMessageSpecVersion,
			"deus.cargo_tracker",
			id.String(),
			"cargo",
			occurredOn,
			eventData,
		),
		oldStatus:  oldStatus.String(),
		newStatus:  newStatus.String(),
		occurredOn: occurredOn,
	}, nil
}
