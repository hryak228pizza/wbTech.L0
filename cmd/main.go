package main

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/hryak228pizza/wbTech.L0/internal/config"
	sqlc "github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/gen"
	"github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/repository"
	"github.com/hryak228pizza/wbTech.L0/internal/logger"
	"github.com/hryak228pizza/wbTech.L0/internal/transport/handler"
	_ "github.com/hryak228pizza/wbTech.L0/internal/transport/handler/docs"
	c "github.com/hryak228pizza/wbTech.L0/internal/transport/kafka/consumer"
	p "github.com/hryak228pizza/wbTech.L0/internal/transport/kafka/producer"
	"github.com/hryak228pizza/wbTech.L0/pkg/cache"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

//	@title			Order service API
//	@version		1.0
//	@description	Service for processing, storing and displaying order data

//	@host		localhost:8080
//	@BasePath	/

func main() {

	// root context for graceful shutdown
	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// first initialization of logger
	logger.Logger()
	defer logger.L().Sync()

	// load app config
	cfg := config.LoadCfg()

	// open database
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", cfg.Dsn)
		if err == nil {
			err = db.PingContext(rootCtx)
		}
		if err == nil {
			break
		}
		logger.L().Warn("failed to connect to database, retrying...",
			zap.Int("attempt", i+1),
			zap.String("error", err.Error()),
		)
	}
	if err != nil {
		logger.L().Fatal("failed to open database")
	}
	defer db.Close()

	// var for db queries
	queries := sqlc.New(db)

	// var for repo
	repo := repository.NewOrderRepository(db, queries)

	// var for template
	tmpl := template.Must(template.ParseGlob("templates/*"))

	// create cache
	lru, err := cache.NewCache(cfg.CacheSize, repo)
	if err != nil {
		logger.L().Fatal("failed to create cache",
			zap.String("error", err.Error()),
		)
	}

	// handlers setup
	h := handler.NewHandler(repo, lru, tmpl)

	// create router
	r := mux.NewRouter()
	r.HandleFunc("/", h.Page).Methods("GET")
	r.HandleFunc("/order/{id}", h.List).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// HTTP Server setup
	srv := &http.Server{
		Addr:    cfg.HttpPort,
		Handler: r,
	}

	logger.L().Info("starting server",
		zap.String("host", "localhost"),
		zap.String("port", cfg.HttpPort),
	)

	// new wait group
	var wg sync.WaitGroup

	// run http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.L().Info("http server starting", zap.String("addr", "localhost"+cfg.HttpPort))
		if err := http.ListenAndServe(cfg.HttpPort, r); err != nil && err != http.ErrServerClosed {
			logger.L().Fatal("http.ListenAndServe failed", zap.Error(err))
		}
	}()

	// run kafka consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Consumer(rootCtx, cfg, h.Cache, repo)
		logger.L().Info("kafka consumer stopped")

	}()

	// run kafka producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.Producer(cfg)
		logger.L().Info("kafka producer stopped")
	}()

	// signal channel
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// wait for signal
	sig := <-sigCh
	logger.L().Info("shutdown signal received", zap.String("signal", sig.String()))

	// stop kafka cosumer and producer
	cancel()

	// graceful HTTP shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.L().Error("http server shutdown error", zap.Error(err))
	} else {
		logger.L().Info("http server stopped gracefully")
	}

	// wait for goroutines to finish
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		logger.L().Info("all goroutines finished")
	case <-time.After(20 * time.Second):
		logger.L().Warn("timeout waiting for goroutines, forcing exit")
	}

	logger.L().Info("service stopped")
}
