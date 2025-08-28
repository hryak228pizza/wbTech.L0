package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/hryak228pizza/wbTech.L0/internal/logger"
	h "github.com/hryak228pizza/wbTech.L0/internal/transport/handler"
	c "github.com/hryak228pizza/wbTech.L0/internal/transport/kafka/consumer"
	p "github.com/hryak228pizza/wbTech.L0/internal/transport/kafka/producer"
    _ "github.com/hryak228pizza/wbTech.L0/internal/transport/handler/docs"
	"github.com/swaggo/http-swagger"
	"github.com/hryak228pizza/wbTech.L0/pkg/cache"
	"go.uber.org/zap"
	"github.com/joho/godotenv"
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

	// load .env variables
	err := godotenv.Load()
	if err != nil {
		logger.L().Fatal("Error loading .env file")
	}

	// open database
    dsn := os.Getenv("DB_CONN_STR")
    db, err := sql.Open("postgres", dsn)
    if err != nil { 
        logger.L().Fatal("failed to open database")
    }     
    defer db.Close()

	// create cache
	size, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		logger.L().Fatal("failed to parse cachesize from .env", 
			zap.String("error", err.Error()),
		)
	}
	lru, err := cache.NewCache(size, db)
	if err != nil {
		logger.L().Fatal("failed to create cache", 
			zap.String("error", err.Error()),
		)
	}
        
	// handlers setup
    handlers := &h.Handler{
		DB:    db,
		Tmpl:  template.Must(template.ParseGlob("templates/*")),
		Cache: lru,
	}

	// create router 
    r := mux.NewRouter()
    r.HandleFunc("/", handlers.Page).Methods("GET")
	r.HandleFunc("/order/{id}", handlers.List).Methods("GET")	
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    logger.L().Info("starting server",
		zap.String("host", "localhost"),
		zap.Int("port", 8080),
	)
	
	// run kafka consumer and producer with server
    go c.Consumer(lru, handlers.DB)
    go p.Producer()
	http.ListenAndServe(":8080", r)    
}