package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/hryak228pizza/wbTech.L0/internal/infrastructure/db/repository"
	"github.com/hryak228pizza/wbTech.L0/pkg/cache"
	"go.uber.org/zap"
)

type Handler struct {
	Repo  repository.OrderRepository
	Tmpl  *template.Template
	Cache *cache.Cache
}

func NewHandler(repo repository.OrderRepository, c *cache.Cache, tmpl *template.Template) *Handler {
	return &Handler{
		Repo:  repo,
		Tmpl:  tmpl,
		Cache: c,
	}
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

		// order info parse
		order, err := h.Repo.GetByUID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJson(w, http.StatusNotFound, map[string]string{"error": "order not found"})
				return
			}
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			return
		}

		// save current order into cache
		h.Cache.SetOrder(order)

		writeJson(w, http.StatusOK, order)
	}
}

// Page godoc
// @Summary Display web page
// @Description Returns an HTML page with an order search form.
// @Tags page
// @Produce  html
// @Success 200 {string} string "HTML content"
// @Router / [get]
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

// helper for writing json
func writeJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
