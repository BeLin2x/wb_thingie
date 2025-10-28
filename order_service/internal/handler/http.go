package handler

import (
	"encoding/json"
	"net/http"
	"order_service/internal/cache"

	"github.com/gorilla/mux"
)

type Handler struct {
	cache *cache.Cache
}

func New(cache *cache.Cache) *Handler {
	return &Handler{
		cache: cache,
	}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["id"]

	order, exists := h.cache.Get(orderUID)
	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders := h.cache.GetAll()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}