package events

type EventType string

const (
	EventCharacterCreationStarted EventType = "EventCharacterCreationStarted"
	EventCharacterInfoSubmitted   EventType = "EventCharacterInfoSubmitted"
	EventCharacterClassSubmitted  EventType = "EventCharacterClassSubmitted"
)

type EventHandler func(data interface{})

type Subscription struct {
	EventType EventType
	Handler   EventHandler
	Channel   chan interface{}
}

type Event struct {
	EventType EventType
	Data      interface{}
}
