package cargotrackingdomain

import (
	"time"
)

type Tracking []TrackingItem

func NewEmptyTracking() Tracking {
	return make(Tracking, 0)
}

type TrackingItem struct {
	id           TrackingID
	entryType    TrackingEntryType
	createdAt    time.Time
	statusBefore *string
	statusAfter  *string
}

func NewTrackingItemFromPrimitives(p TrackingItemPrimitives) TrackingItem {
	return TrackingItem{
		id:           TrackingID(p.ID),
		entryType:    TrackingEntryType(p.EntryType),
		createdAt:    p.CreatedAt,
		statusBefore: p.StatusBefore,
		statusAfter:  p.StatusAfter,
	}
}

func NewTrackingOnCargoCreated(trackingID TrackingID, cargoStatus string, createdAt time.Time) TrackingItem {
	return TrackingItem{
		id:           trackingID,
		entryType:    TrackingEntryTypeCreated,
		createdAt:    createdAt,
		statusBefore: &cargoStatus,
		statusAfter:  &cargoStatus,
	}
}

func NewTrackingOnCargoStatusChanged(
	id TrackingID,
	createdAt time.Time,
	statusBefore, statusAfter string,
) TrackingItem {
	return TrackingItem{
		id:           id,
		entryType:    TrackingEntryTypeStatusChanged,
		createdAt:    createdAt,
		statusBefore: &statusBefore,
		statusAfter:  &statusAfter,
	}
}
