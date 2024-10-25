package events

import "sync"

type EventManager struct {
	subscribers map[string][]chan interface{}
	mutex       sync.RWMutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscribers: make(map[string][]chan interface{}),
	}
}

func (em *EventManager) Subscribe(eventType string, ch chan interface{}) {
	em.mutex.Lock()
	defer em.mutex.Unlock()
	em.subscribers[eventType] = append(em.subscribers[eventType], ch)
}

func (em *EventManager) Publish(eventType string, data interface{}) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	if channels, found := em.subscribers[eventType]; found {
		for _, ch := range channels {
			go func(c chan interface{}) {
				c <- data
			}(ch)
		}
	}
}

func (em *EventManager) Unsubscribe(eventType string, ch chan interface{}) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if subscribers, found := em.subscribers[eventType]; found {
		for i, subscriber := range subscribers {
			if subscriber == ch {
				em.subscribers[eventType] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
	}
}
