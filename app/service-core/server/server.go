package server

import (
	"app/pkg/auth"
	"service-core/broker"
	"service-core/config"
	"service-core/domain/user"
	"service-core/server/grpc"
	"service-core/server/rest"
	"service-core/storage"
	"service-core/storage/query"
)

type Server struct {
	Config      *config.Config
	Storage     *storage.Storage
	GRPCServer  *grpc.Handler
	RESTServer  *rest.Handler
	AuthService *auth.Service
	UserService *user.Service
}

func New(cfg *config.Config, s *storage.Storage, b *broker.Broker) *Server {
	store := query.New(s.Conn)
	authService := auth.NewService()
	userService := user.NewService(cfg, store)

	return &Server{
		Config:      cfg,
		Storage:     s,
		GRPCServer:  grpc.NewHandler(cfg, authService, userService),
		RESTServer:  rest.NewHandler(cfg, store, authService, userService, b),
		AuthService: authService,
		UserService: userService,
	}
}

func (s *Server) Start() {
	rest.Run(s.RESTServer)
	grpc.Run(s.GRPCServer)
}
