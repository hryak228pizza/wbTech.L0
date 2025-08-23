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
	"go.uber.org/zap"
)

func main() {

    // looger init
    logger := zap.NewExample()
    defer logger.Sync()

	dsn := "user=postgres password=123 dbname=postgres sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil { 
        logger.Info("failed to open database")
    }     
    defer db.Close()
        
    handlers := &h.Handler{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("templates/*")),
	}

    r := mux.NewRouter()
    r.HandleFunc("/", handlers.Page).Methods("GET")
	r.HandleFunc("/order/{id}", handlers.List).Methods("GET")

    logger.Info("starting server",
		zap.String("logger", "ZAP"),
		zap.String("host", "localhost"),
		zap.Int("port", 8080),
	)
	//fmt.Println("starting server at :8080")
    go c.Consumer(handlers.DB)
    go p.Producer()
	http.ListenAndServe(":8080", r)

    
}