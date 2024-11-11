package events

type EventType string

const (
	EventCharacterInfoSubmitted  EventType = "event_character_info_submit"
	EventCharacterClassSubmitted EventType = "event_character_class_submit"
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
