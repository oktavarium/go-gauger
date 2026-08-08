package gaugeserver

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/cipher"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/grpcserver"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/gzip"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/handlers"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/hash"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/ipsec"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// GaugeServer - управляющий сервис для сбора метрик
type GaugerServer struct {
	ctx  context.Context
	srv  *http.Server
	grpc *grpcserver.GrpcServer
}

// NewGaugeServer - конструктор сервиса ядл сбора метрик
func NewGaugerServer(
	ctx context.Context,
	addr string,
	grpcAddr string,
	filename string,
	restore bool,
	timeout time.Duration,
	dsn string,
	key string,
	pkFile string,
	subnet string,
) (*GaugerServer, error) {

	router := chi.NewRouter()
	var s storage.Storage
	var err error

	if len(dsn) == 0 {
		s, err = storage.NewInMemoryStorage(ctx, filename, restore, timeout)
	} else {
		s, err = storage.NewPostgresqlStorage(dsn)
	}
	if err != nil {
		return nil, fmt.Errorf("error on creating storage: %w", err)
	}

	handler := handlers.NewHandler(s)

	if len(pkFile) != 0 {
		c, err := cipher.NewCipher(pkFile)
		if err != nil {
			return nil, fmt.Errorf("error on creating cipher: %w", err)
		}
		router.Use(c.CipherMiddleware)
	}

	sec, err := ipsec.NewIPSec(subnet)
	if err != nil {
		return nil, fmt.Errorf("error on creating ipsec: %w", err)
	}

	router.Use(logger.LoggerMiddleware)
	router.Use(sec.IPSecMiddleware)
	if len(key) != 0 {
		router.Use(hash.HashMiddleware([]byte(key)))
	}
	router.Use(gzip.GzipMiddleware)
	router.Get("/", handler.GetHandle)
	router.Get("/ping", handler.PingHandle)
	router.Post("/update/", handler.UpdateJSONHandle)
	router.Post("/updates/", handler.UpdatesHandle)
	router.Post("/value/", handler.ValueJSONHandle)
	router.Post("/update/{type}/{name}/{value}", handler.UpdateHandle)
	router.Get("/value/{type}/{name}", handler.ValueHandle)

	server := &GaugerServer{
		ctx: ctx,
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		grpc: grpcserver.NewGrpcServer(grpcAddr, s),
	}

	return server, nil
}

// ListenAndServer - запуск сервиса
func (g *GaugerServer) ListenAndServe() error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-g.ctx.Done()
		if err := g.srv.Shutdown(context.Background()); err != nil {
			logger.Logger().Error("error",
				zap.String("func", "ListenAndServer"),
				zap.Error(err))
		}

		g.grpc.Shutdown()
	}()

	eg := errgroup.Group{}
	eg.Go(g.grpc.ListenAndServe)
	eg.Go(g.srv.ListenAndServe)

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error on listen and serve: %w", err)
	}

	wg.Wait()

	return nil
}
