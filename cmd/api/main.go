package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	// Fetch environment variables with default values
	portStr := os.Getenv("PORT")
	if portStr == "" {
		logger.Info("Port not specified, Use deafult port")
		portStr = "4000" // Default port
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Error("Invalid port number:")
		os.Exit(1)
	}
	cfg.port = port

	cfg.env = os.Getenv("ENV")
	if cfg.env == "" {
		logger.Info("No environment variable. Using default environment variable")
		cfg.env = "development" // Default environment
	}

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5  * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}