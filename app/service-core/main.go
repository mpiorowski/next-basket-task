package main

import (
	"app/pkg"
	"context"
	"log/slog"
	"service-core/broker"
	"service-core/config"
	"service-core/server"
	"service-core/storage"
)

func main() {
	// Load the configuration
	cfg := config.LoadConfig()

	// Set up the logger
	pkg.InitLogger(cfg.LogLevel)

	// Connect to the database
	s, clean, err := storage.NewStorage(cfg)
	if err != nil {
		slog.Error("Error opening database", "error", err)
		panic(err)
	}
	defer clean()

	err = s.Conn.PingContext(context.Background())
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		panic(err)
	}
	slog.Info("Database connected")

	// Initialize the broker
	b, err := broker.New(cfg.NatsURL)
	if err != nil {
		slog.Error("Error initializing broker", "error", err)
		panic(err)
	}
	// Run the broker
	b.Run(context.Background())

	// Set up the servers
	srv := server.New(cfg, s, b)

	// Run the servers
	srv.Start()

	select {}
}
