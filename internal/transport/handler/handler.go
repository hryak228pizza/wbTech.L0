package handler

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/hryak228pizza/wbTech.L0/internal/model"
	"github.com/hryak228pizza/wbTech.L0/pkg/cache"
	"go.uber.org/zap"
)

type Handler struct {
	DB    *sql.DB
	Tmpl  *template.Template
	Cache *cache.Cache
}

// List godoc
// @Summary Get order data by ID
// @Description Returns an order by its unique identifier (UID).
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order UID"
// @Success 200 {object} model.Order
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
//
// formatting all order data to json
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

	// parse query param
	vars := mux.Vars(r)
	id := vars["id"]

	// check cashed data
	if order, ok := h.Cache.GetOrder(id); ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
		return
	} else {
		orderInfo := &model.Order{}

		// order info parse
		ordersRow := h.DB.QueryRow("SELECT * FROM orders WHERE order_uid = $1", id)
		if err := ordersRow.Scan(&orderInfo.OrderUID, &orderInfo.TrackNumber, &orderInfo.Entry,
			&orderInfo.Locale, &orderInfo.InternalSignature, &orderInfo.CustomerID, &orderInfo.DeliveryService,
			&orderInfo.ShardKey, &orderInfo.SmID, &orderInfo.DateCreated, &orderInfo.OofShard); err != nil {
			if err == sql.ErrNoRows {
				// return code 404
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
				return
			} else {
				// return code 500
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "internal server error 102"})
				return
			}
		}

		// delivery info parse
		delivery := &model.Delivery{}
		deliveryRow := h.DB.QueryRow("SELECT * FROM delivery WHERE order_uid = $1", id)
		if err := deliveryRow.Scan(&delivery.OrderUID, &delivery.Name, &delivery.Phone,
			&delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email); err != nil {
			if err == sql.ErrNoRows {
				// return code 404
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
				return
			} else {
				// return code 500
				w.WriteHeader(http.StatusInternalServerError)
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
				// return code 404
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
				return
			} else {
				// return code 500
				w.WriteHeader(http.StatusInternalServerError)
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
					// no items?
				} else {
					// return code 500
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
					return
				}
			}
			items = append(items, item)
		}
		orderInfo.Items = items

		// save current order into cache
		h.Cache.SetOrder(orderInfo)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orderInfo)
	}
}

// Page godoc
// @Summary Display web page
// @Description Returns an HTML page with an order search form.
// @Tags page
// @Produce  html
// @Success 200 {string} string "HTML content"
// @Router / [get]
//
// show a web page to user
func (h *Handler) Page(w http.ResponseWriter, r *http.Request) {

	// looger init
	logger := zap.NewExample()
	defer logger.Sync()

	err := h.Tmpl.ExecuteTemplate(w, "index.html", "")
	if err != nil {
		logger.Info("failed to execute html template",
			zap.String("url", r.URL.Path),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
