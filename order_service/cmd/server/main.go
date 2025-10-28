package main

import (
	"log"
	"net/http"
	"order_service/internal/cache"
	"order_service/internal/database"
	"order_service/internal/handler"
	"order_service/internal/nats"
	"path/filepath"

	"github.com/gorilla/mux"
)

func main() {
	const (
		clusterID  = "test-cluster"
		clientID   = "order-service"
		subject    = "orders"
		connStr    = "host=localhost port=5432 user=order_service_user password=hail dbname=order_service_db sslmode=disable"
		httpAddr   = ":8080"
	)

	// Инициализация базы данных
	db, err := database.New(connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database")

	// Инициализация кэша
	cache := cache.New()

	// Восстановление кэша из базы данных при запуске
	orders, err := db.RestoreCache()
	if err != nil {
		log.Fatal("Failed to restore cache from database:", err)
	}
	cache.Restore(orders)
	log.Printf("Restored %d orders from database", len(orders))

	// Инициализация NATS подписчика
	subscriber := nats.New(clusterID, clientID, subject, cache, db)
	if err := subscriber.Start(); err != nil {
		log.Fatal("Failed to start NATS subscriber:", err)
	}
	log.Println("NATS subscriber started")

	// Инициализация HTTP сервера
	handler := handler.New(cache)
	router := mux.NewRouter()

	// API маршруты
	router.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")
	router.HandleFunc("/orders", handler.GetOrders).Methods("GET")

	// Веб-маршруты
	router.HandleFunc("/", serveIndex)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	log.Printf("HTTP server starting on %s", httpAddr)
	log.Printf("Web interface available at: http://localhost%s", httpAddr)
	if err := http.ListenAndServe(httpAddr, router); err != nil {
		log.Fatal("HTTP server error:", err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "templates", "index.html"))
}