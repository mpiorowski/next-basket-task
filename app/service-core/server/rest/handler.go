package rest

import (
	"app/pkg/auth"
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
}

func NewHandler(
	config *config.Config,
	store *query.Queries,
	authService *auth.Service,
	userService *user.Service,
	broker *broker.Broker,
) *Handler {
	return &Handler{
		cfg:         config,
		store:       store,
		authService: authService,
		userService: userService,
		broker:      broker,
	}
}
