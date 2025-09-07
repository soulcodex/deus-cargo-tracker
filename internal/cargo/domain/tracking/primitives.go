package cargotrackingdomain

import (
	"time"
)

type TrackingPrimitives []TrackingItemPrimitives
type TrackingItemPrimitives struct {
	ID           string
	CargoID      string
	EntryType    string
	StatusBefore *string
	StatusAfter  *string
	CreatedAt    time.Time
}

func NewTrackingPrimitives(cargoID string, items Tracking) TrackingPrimitives {
	if len(items) == 0 {
		return make(TrackingPrimitives, 0)
	}

	primitives := make(TrackingPrimitives, len(items))
	for i, item := range items {
		primitives[i] = NewTrackingItemPrimitives(cargoID, item)
	}

	return primitives
}

func NewTrackingItemPrimitives(cargoID string, t TrackingItem) TrackingItemPrimitives {
	var statusBefore, statusAfter *string

	if t.statusBefore != nil {
		sb := *t.statusBefore
		statusBefore = &sb
	}

	if t.statusAfter != nil {
		sa := *t.statusAfter
		statusAfter = &sa
	}

	return TrackingItemPrimitives{
		ID:           t.id.String(),
		CargoID:      cargoID,
		EntryType:    t.entryType.String(),
		StatusBefore: statusBefore,
		StatusAfter:  statusAfter,
		CreatedAt:    t.createdAt,
	}
}
