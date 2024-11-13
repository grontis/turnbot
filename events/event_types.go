package events

import "time"

type EventType string

const (
	EventCharacterCreationStarted EventType = "EventCharacterCreationStarted"
	EventCharacterInfoSubmitted   EventType = "EventCharacterInfoSubmitted"
	EventCharacterClassSubmitted  EventType = "EventCharacterClassSubmitted"
	EventCharacterUpdated         EventType = "EventCharacterUpdated"
)

type CharacterCreationStartedData struct {
	UserID    string
	Timestamp time.Time
}

type CharacterInfoSubmittedData struct {
	UserID string
	Name   string
	Age    string
}

type CharacterClassSubmittedData struct {
	UserID    string
	ClassName string
	Level     int
}

type CharacterUpdatedData struct {
	UserID    string
	Timestamp time.Time
}

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
