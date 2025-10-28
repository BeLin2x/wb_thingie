package nats

import (
	"encoding/json"
	"log"
	"order_service/internal/cache"
	"order_service/internal/database"
	"order_service/internal/models"

	"github.com/nats-io/stan.go"
)

type Subscriber struct {
	clusterID string
	clientID  string
	subject   string
	cache     *cache.Cache
	db        *database.DB
}

func New(clusterID, clientID, subject string, cache *cache.Cache, db *database.DB) *Subscriber {
	return &Subscriber{
		clusterID: clusterID,
		clientID:  clientID,
		subject:   subject,
		cache:     cache,
		db:        db,
	}
}

func (s *Subscriber) Start() error {
	sc, err := stan.Connect(s.clusterID, s.clientID)
	if err != nil {
		return err
	}

	_, err = sc.Subscribe(s.subject, func(msg *stan.Msg) {
		if err := s.processMessage(msg.Data); err != nil {
			log.Printf("Error processing message: %v", err)
		}
	}, stan.DurableName("order-service"))

	if err != nil {
		return err
	}

	log.Printf("Subscribed to NATS subject: %s", s.subject)
	return nil
}

func (s *Subscriber) processMessage(data []byte) error {
	if !json.Valid(data) {
		log.Printf("Invalid JSON received")
		return nil
	}

	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		log.Printf("Error unmarshaling order: %v", err)
		return nil
	}

	if order.OrderUID == "" || order.TrackNumber == "" {
		log.Printf("Invalid order data: missing required fields")
		return nil
	}

	if err := s.db.SaveOrder(&order); err != nil {
		log.Printf("Error saving order to DB: %v", err)
		return err
	}

	s.cache.Set(&order)

	log.Printf("Order processed successfully: %s", order.OrderUID)
	return nil
}