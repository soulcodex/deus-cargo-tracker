package cargotrackingdomain

const (
	trackingEntryTypeCreated       TrackingEntryType = "cargo.created"
	trackingEntryTypeStatusChanged TrackingEntryType = "cargo.status_changed"
)

type TrackingEntryType string

func (s TrackingEntryType) String() string {
	return string(s)
}
