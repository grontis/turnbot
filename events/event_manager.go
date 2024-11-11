package events

import "sync"

type EventManager struct {
	subscribers map[EventType][]chan interface{}
	mutex       sync.RWMutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscribers: make(map[EventType][]chan interface{}),
	}
}

func (em *EventManager) Subscribe(sub Subscription) chan interface{} {
	if sub.Channel == nil {
		sub.Channel = make(chan interface{})
	}

	em.mutex.Lock()
	defer em.mutex.Unlock()
	em.subscribers[sub.EventType] = append(em.subscribers[sub.EventType], sub.Channel)

	go func() {
		for event := range sub.Channel {
			sub.Handler(event)
		}
	}()

	return sub.Channel
}

func (em *EventManager) Publish(event Event) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	if channels, found := em.subscribers[event.EventType]; found {
		for _, ch := range channels {
			go func(c chan interface{}) {
				c <- event.Data
			}(ch)
		}
	}
}

func (em *EventManager) Unsubscribe(sub Subscription) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if subscribers, found := em.subscribers[sub.EventType]; found {
		for i, subscriber := range subscribers {
			if subscriber == sub.Channel {
				em.subscribers[sub.EventType] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
	}
}
