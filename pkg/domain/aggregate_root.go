package domain

type AggregateRoot struct {
	events []Event
}

func NewAggregateRoot() *AggregateRoot {
	return &AggregateRoot{
		events: make([]Event, 0),
	}
}

func (ar *AggregateRoot) RecordEvent(event Event) {
	ar.events = append(ar.events, event)
}

func (ar *AggregateRoot) PullEvents() []Event {
	events := ar.events
	ar.events = make([]Event, 0)
	return events
}
