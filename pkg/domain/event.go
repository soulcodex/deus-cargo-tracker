package domain

import (
	"github.com/soulcodex/deus-cargo-tracker/pkg/messaging"
)

// Event alias to message just for simplicity in a real world scenario
// we would have to decouple this even more.
type Event = messaging.Message

// EventPublisher alias to message publisher just for simplicity in a real world scenario
// we would have to decouple this even more.
type EventPublisher = messaging.Publisher
