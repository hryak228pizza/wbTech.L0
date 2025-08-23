package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	c "github.com/hryak228pizza/wbTech.L0/internal/transport/kafka/consumer"
	p "github.com/hryak228pizza/wbTech.L0/internal/transport/kafka/producer"
	h "github.com/hryak228pizza/wbTech.L0/internal/transport/http"
	"github.com/hryak228pizza/wbTech.L0/internal/logger"
	"go.uber.org/zap"
)


func main() {

	// first initialization of logger
	logger.Logger()
    defer logger.L().Sync()

	// open database
    dsn := "user=postgres password=123 dbname=postgres sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil { 
        logger.L().Info("failed to open database")
    }     
    defer db.Close()
        
	// handlers setup
    handlers := &h.Handler{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("templates/*")),
	}

	// create router 
    r := mux.NewRouter()
    r.HandleFunc("/", handlers.Page).Methods("GET")
	r.HandleFunc("/order/{id}", handlers.List).Methods("GET")

    logger.L().Info("starting server",
		zap.String("host", "localhost"),
		zap.Int("port", 8080),
	)
	
	// run kafka consumer and producer with server
    go c.Consumer(handlers.DB)
    go p.Producer()
	http.ListenAndServe(":8080", r)    
}