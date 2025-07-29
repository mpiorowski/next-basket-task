package rest

import (
	"app/pkg/event"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type PublishRequest struct {
	Subject string          `json:"subject"`
	Data    json.RawMessage `json:"data"`
}

func (h *Handler) handleBrokerPublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeResponse(h.cfg, w, r, nil, fmt.Errorf("method not allowed: %s", r.Method))
		return
	}
	var req PublishRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode request body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	h.eventStore.Add(event.Event{
		Subject: req.Subject,
		Data:    req.Data,
	})

	err = h.broker.Publish(r.Context(), req.Subject, req.Data)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to publish message: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(h.cfg, w, r, nil, nil)
}

func (h *Handler) handleGetAllEvents(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		events := h.eventStore.GetAll()
		writeResponse(h.cfg, w, r, events, nil)
		return
	}
	events := h.eventStore.GetByTenantID(tenantID)
	writeResponse(h.cfg, w, r, events, nil)
}

func (h *Handler) handleBrokerSubscribe(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling broker subscription", "remote_addr", r.RemoteAddr)
	//nolint:exhaustruct
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade connection", "error", err)
		return
	}
	defer conn.Close()

	subject := r.URL.Query().Get("subject")
	if subject == "" {
		slog.Error("Missing subject in query parameter")
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "subject query parameter is required"))
		return
	}

	clientID, clientChan := h.broker.Subscribe(subject)
	defer h.broker.Unsubscribe(subject, clientID)

	var mu sync.Mutex
	// Read from broker channel and write to websocket
	go func() {
		for msg := range clientChan {
			mu.Lock()
			err := conn.WriteMessage(websocket.TextMessage, msg.Data)
			mu.Unlock()
			if err != nil {
				slog.Error("Failed to write message to websocket", "error", err, "clientID", clientID)
				break
			}
		}
	}()

	slog.Info("Client subscribed to subject", "subject", subject, "remote_addr", conn.RemoteAddr())

	// Keep the connection alive by reading from the client
	for {
		_, _, err := conn.NextReader()
		if err != nil {
			slog.Info("Client disconnected", "remote_addr", conn.RemoteAddr(), "reason", err)
			break
		}
	}
}