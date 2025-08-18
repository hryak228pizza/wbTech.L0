package main

import (
    "database/sql"
    "net/http"
    "github.com/gorilla/mux"
    "html/template"
    "fmt"
    "encoding/json"
    _ "time"
    _ "log"
    _ "github.com/lib/pq"
    //kafka "github.com/segmentio/kafka-go"

    model "github.com/hryak228pizza/wbTech.L0/cmd"
)


type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func Add(db *sql.DB) (sql.Result, error){
    //result, err := db.Exec("insert into delivery (order_uid, name, phone, zip, city, address, region, email) values ($1, $2, $3, $4, $5, $6, $7, $8)", 
    //    "aaaabbbbbcccccctest", "testname", "+9420000042", "152625", "testcity", "testaddress", "testregion", "testmail@com")
	return db.Exec("insert into orders (order_uid, track_number) values ($1, $2)",
		"12345", "TESTTRACKERSECOND")    
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
	id := vars["id"]

    orderInfo := &model.Order{}

    // order info parse
    ordersRow := h.DB.QueryRow("SELECT * FROM orders WHERE order_uid = $1", id)
    if err := ordersRow.Scan(&orderInfo.OrderUID, &orderInfo.TrackNumber, &orderInfo.Entry, &orderInfo.Locale, &orderInfo.InternalSignature, &orderInfo.CustomerID, &orderInfo.DeliveryService, &orderInfo.ShardKey, &orderInfo.SmID, &orderInfo.DateCreated, &orderInfo.OofShard); err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound) // 404
            // w.Write([]byte(`{"error": "order not found"}`))
            json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
            return
        } else {
            w.WriteHeader(http.StatusInternalServerError) // 500
            // w.Write([]byte(`{"error": "internal server error"}`))
            json.NewEncoder(w).Encode(map[string]string{"error": "internal server error 102"})
            return
        }
    }
    
    // delivery info parse
    delivery := &model.Delivery{}
    deliveryRow := h.DB.QueryRow("SELECT * FROM delivery WHERE order_uid = $1", id)
    if err := deliveryRow.Scan(&delivery.OrderUID, &delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email); err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound) // 404
            // w.Write([]byte(`{"error": "order not found"}`))
            json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
            return
        } else {
            w.WriteHeader(http.StatusInternalServerError) // 500
            // w.Write([]byte(`{"error": "internal server error"}`))
            json.NewEncoder(w).Encode(map[string]string{"error": "internal server error 119"})
            return
        }
    }
    orderInfo.Delivery = *delivery

    // payment info parse
    payment := &model.Payment{}
    paymentRow := h.DB.QueryRow("SELECT * FROM payment WHERE transaction = $1", id)
    if err := paymentRow.Scan(&payment.Transaction, &payment.RequestID, 
        &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDT,
        &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
        ); err != nil {
            if err == sql.ErrNoRows {
                w.WriteHeader(http.StatusNotFound) // 404
                // w.Write([]byte(`{"error": "order not found"}`))
                json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
                return
            } else {
                w.WriteHeader(http.StatusInternalServerError) // 500
                // w.Write([]byte(`{"error": "internal server error"}`))
                json.NewEncoder(w).Encode(map[string]string{"error": "internal server error 137"})
                return
            }
    }
    orderInfo.Payment = *payment

    // items Parse
    items := []*model.Item{}
    itemsRows, err := h.DB.Query("SELECT * FROM items WHERE track_number = $1", orderInfo.TrackNumber)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "internal server error 1"})
        return
    }
    defer itemsRows.Close()
    for itemsRows.Next() {
        item := &model.Item{}
        if err := itemsRows.Scan(&item.ID, &item.OrderUID, &item.ChrtID,
                &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, 
                &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
            ); err != nil {
                if err == sql.ErrNoRows {
                    //w.WriteHeader(http.StatusNotFound) // 404
                    //json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
                    //return
                } else {
                    w.WriteHeader(http.StatusInternalServerError) // 500
                    json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
                    return
                }
        }
        items = append(items, item)
    }
    orderInfo.Items = items    

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(orderInfo)
}

func (h *Handler) Page(w http.ResponseWriter, r *http.Request) {
    err := h.Tmpl.ExecuteTemplate(w, "index.html", "")
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	dsn := "user=postgres password=123 dbname=postgres sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil { panic(err) }     
    defer db.Close()
    
    // result, err := Add(db)
    // if err != nil{ panic(err) }	
    // fmt.Println(result.RowsAffected())
    
    handlers := &Handler{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("templates/*")),
	}

    r := mux.NewRouter()
    r.HandleFunc("/", handlers.Page).Methods("GET")
	r.HandleFunc("/order/{id}", handlers.List).Methods("GET")

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", r)

    
}