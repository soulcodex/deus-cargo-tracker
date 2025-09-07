package cargotrackingdomain

const (
	TrackingEntryTypeCreated       TrackingEntryType = "cargo.created"
	TrackingEntryTypeStatusChanged TrackingEntryType = "cargo.status_changed"
)

type TrackingEntryType string

func (s TrackingEntryType) String() string {
	return string(s)
}
