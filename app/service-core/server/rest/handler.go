package rest

import (
	"app/pkg/auth"
	"app/pkg/event"
	"service-core/broker"
	"service-core/config"
	"service-core/domain/user"
	"service-core/storage/query"
)

type Handler struct {
	cfg         *config.Config
	store       *query.Queries
	authService *auth.Service
	userService *user.Service
	broker      *broker.Broker
	eventStore  *event.Store
}

func NewHandler(
	config *config.Config,
	store *query.Queries,
	authService *auth.Service,
	userService *user.Service,
	broker *broker.Broker,
	eventStore *event.Store,
) *Handler {
	return &Handler{
		cfg:         config,
		store:       store,
		authService: authService,
		userService: userService,
		broker:      broker,
		eventStore:  eventStore,
	}
}
