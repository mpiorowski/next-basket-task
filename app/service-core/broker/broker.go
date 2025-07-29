package broker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Broker struct {
	nc      *nats.Conn
	msgChan chan *nats.Msg
	subs    map[string]map[string]chan *nats.Msg
	mu      sync.RWMutex
}

func New(natsURL string) (*Broker, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	b := &Broker{
		nc:      nc,
		msgChan: make(chan *nats.Msg, 256),
		subs:    make(map[string]map[string]chan *nats.Msg),
		mu:      sync.RWMutex{},
	}

	_, err = b.nc.ChanSubscribe(">", b.msgChan)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to wildcard subject: %w", err)
	}

	slog.Info("Connected to NATS and subscribed to wildcard subject", "url", natsURL)
	return b, nil
}

func (b *Broker) Run(ctx context.Context) {
	slog.InfoContext(ctx, "Starting broker message handler")
	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.InfoContext(ctx, "Shutting down broker", "reason", ctx.Err())
				b.nc.Close()
				return
			case msg := <-b.msgChan:
				b.mu.RLock()
				if clients, ok := b.subs[msg.Subject]; ok {
					for id, clientChan := range clients {
						go func(_ string, ch chan *nats.Msg, m *nats.Msg) {
							ch <- m
						}(id, clientChan, msg)
					}
				}
				b.mu.RUnlock()
			}
		}
	}()
}

func (b *Broker) Publish(ctx context.Context, subject string, data []byte) error {
	err := b.nc.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish message to subject %s: %w", subject, err)
	}
	slog.InfoContext(ctx, "Published message", "subject", subject)
	return nil
}

func (b *Broker) Subscribe(subject string) (string, chan *nats.Msg) {
	b.mu.Lock()
	defer b.mu.Unlock()

	clientID := uuid.New().String()
	clientChan := make(chan *nats.Msg, 64)

	if _, ok := b.subs[subject]; !ok {
		b.subs[subject] = make(map[string]chan *nats.Msg)
	}
	b.subs[subject][clientID] = clientChan

	slog.Info("Client subscribed", "subject", subject, "clientID", clientID)
	return clientID, clientChan
}

func (b *Broker) Unsubscribe(subject, clientID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if clients, ok := b.subs[subject]; ok {
		if clientChan, ok := clients[clientID]; ok {
			close(clientChan)
			delete(clients, clientID)
			slog.Info("Client unsubscribed", "subject", subject, "clientID", clientID)
		}
		if len(clients) == 0 {
			delete(b.subs, subject)
		}
	}
}

func (b *Broker) Close() {
	b.nc.Close()
}
