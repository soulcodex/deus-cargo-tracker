package cargotrackingdomain

import (
	"time"
)

type Tracking []TrackingItem
type TrackingItem struct {
	id           TrackingID
	entryType    TrackingEntryType
	createdAt    time.Time
	statusBefore *string
	statusAfter  *string
}

func NewTrackingOnCargoCreated(trackingID TrackingID, cargoStatus string, createdAt time.Time) TrackingItem {
	return TrackingItem{
		id:          trackingID,
		entryType:   TrackingEntryTypeCreated,
		createdAt:   createdAt,
		statusAfter: &cargoStatus,
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
