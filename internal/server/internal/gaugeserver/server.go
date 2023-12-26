package gaugeserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/gzip"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/handlers"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/hash"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
)

// GaugeServer - управляющий сервис для сбора метрик
type GaugerServer struct {
	router *chi.Mux
	addr   string
}

// NewGaugeServer - конструктор сервиса ядл сбора метрик
func NewGaugerServer(addr string,
	filename string,
	restore bool,
	timeout time.Duration,
	dsn string,
	key string) (*GaugerServer, error) {
	server := &GaugerServer{
		router: chi.NewRouter(),
		addr:   addr,
	}
	var s storage.Storage
	var err error
	if len(dsn) == 0 {
		s, err = storage.NewInMemoryStorage(filename, restore, timeout)
	} else {
		s, err = storage.NewPostgresqlStorage(dsn)
	}
	if err != nil {
		return nil, fmt.Errorf("error on creating storage: %w", err)
	}

	handler := handlers.NewHandler(s)

	server.router.Use(logger.LoggerMiddleware)
	if len(key) != 0 {
		server.router.Use(hash.HashMiddleware([]byte(key)))
	}
	server.router.Use(gzip.GzipMiddleware)
	server.router.Get("/", handler.GetHandle)
	server.router.Get("/ping", handler.PingHandle)
	server.router.Post("/update/", handler.UpdateJSONHandle)
	server.router.Post("/updates/", handler.UpdatesHandle)
	server.router.Post("/value/", handler.ValueJSONHandle)
	server.router.Post("/update/{type}/{name}/{value}", handler.UpdateHandle)
	server.router.Get("/value/{type}/{name}", handler.ValueHandle)

	return server, nil
}

// ListenAndServer - запуск сервиса
func (g *GaugerServer) ListenAndServe() error {
	return http.ListenAndServe(g.addr, g.router)
}
