package event

import (
	"encoding/json"
	"sync"
)

type Event struct {
	Subject string          `json:"subject"`
	Data    json.RawMessage `json:"data"`
}

type Store struct {
	mu     sync.RWMutex
	events []Event
}

func NewStore() *Store {
	return &Store{
		events: make([]Event, 0),
	}
}

func (s *Store) Add(event Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)
}

func (s *Store) GetByTenantID(tenantID string) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var filteredEvents []Event
	for _, event := range s.events {
		// Assuming subject is in the format "tenant.tenant_id"
		if event.Subject == "tenant."+tenantID {
			filteredEvents = append(filteredEvents, event)
		}
	}
	return filteredEvents
}

func (s *Store) GetAll() []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events
}
