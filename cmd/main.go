package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/hryak228pizza/wbTech.L0/internal/config"
	sqlc "github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/gen"
	"github.com/hryak228pizza/wbTech.L0/internal/logger"
	h "github.com/hryak228pizza/wbTech.L0/internal/transport/handler"
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

	// first initialization of logger
	logger.Logger()
	defer logger.L().Sync()

	// load app config
	cfg := config.LoadCfg()

	// open database
	var db *sql.DB
	var err error
	for i:=0; i<10; i++ {
		db, err = sql.Open("postgres", cfg.Dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
	}	
	if err != nil {
		logger.L().Fatal("failed to open database")
	}
	
	// var for db queries
	dbQueries := sqlc.New(db)

	// create cache
	lru, err := cache.NewCache(cfg.CacheSize, dbQueries)
	if err != nil {
		logger.L().Fatal("failed to create cache",
			zap.String("error", err.Error()),
		)
	}

	// handlers setup
	handlers := &h.Handler{
		DB:        db,
		DbQueries: dbQueries,
		Tmpl:      template.Must(template.ParseGlob("templates/*")),
		Cache:     lru,
	}

	// create router
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Page).Methods("GET")
	r.HandleFunc("/order/{id}", handlers.List).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	logger.L().Info("starting server",
		zap.String("host", "localhost"),
		zap.String("port", cfg.HttpPort),
	)

	// run kafka consumer and producer with server
	go c.Consumer(cfg, lru, handlers.DB)
	go p.Producer(cfg)
	http.ListenAndServe(cfg.HttpPort, r)
}
